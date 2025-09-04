package conversation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
)

type ConversationServiceImpl struct {
	llmClient interfaces.LLMClient
}

// NewConversationService creates a new conversation service
func NewConversationService(llmClient interfaces.LLMClient) interfaces.ConversationService {
	return &ConversationServiceImpl{
		llmClient: llmClient,
	}
}

// ValidateSymptom validates if the provided symptom is medical and appropriate
func (cs *ConversationServiceImpl) ValidateSymptom(ctx context.Context, symptom, language string) (bool, string, error) {
	// Basic cleanup - just trim whitespace
	symptom = strings.TrimSpace(symptom)

	// Only check for completely empty input
	if len(symptom) == 0 {
		if language == "am" {
			return false, "እባክዎ የሚሰማዎትን የጤና ችግር ይግለጹ።", nil
		}
		return false, "Please describe your health symptom or concern.", nil
	}

	// Let AI handle all validation logic
	prompt := cs.buildSymptomValidationPrompt(symptom, language)

	response, err := cs.llmClient.ClassifyTriage(ctx, prompt)
	if err != nil {
		return false, "", fmt.Errorf("failed to validate symptom: %w", err)
	}

	isValid, feedback, err := cs.parseSymptomValidationResponse(response)
	if err != nil {
		return false, "", fmt.Errorf("failed to parse symptom validation response: %w", err)
	}

	return isValid, feedback, nil
}

// GenerateQuestions generates 5 follow-up questions based on the initial symptom
func (cs *ConversationServiceImpl) GenerateQuestions(ctx context.Context, symptom, language string) ([]entities.Question, error) {
	prompt := cs.buildQuestionGenerationPrompt(symptom, language)

	// Try up to 3 times to get valid questions
	var lastErr error
	for attempt := 1; attempt <= 3; attempt++ {
		response, err := cs.llmClient.ClassifyTriage(ctx, prompt)
		if err != nil {
			lastErr = fmt.Errorf("failed to generate questions (attempt %d): %w", attempt, err)
			continue
		}

		questions, err := cs.parseQuestionsFromResponse(response)
		if err != nil {
			lastErr = fmt.Errorf("failed to parse questions (attempt %d): %w", attempt, err)
			continue
		}

		// Validate we have at least 3 questions (reduced from 5 for better reliability)
		if len(questions) < 3 {
			lastErr = fmt.Errorf("insufficient questions: got %d, need at least 3 (attempt %d)", len(questions), attempt)
			continue
		}

		// Validate question structure
		validQuestions := []entities.Question{}
		for i, q := range questions {
			if q.Text == "" || q.Type == "" {
				lastErr = fmt.Errorf("invalid question structure at index %d (attempt %d)", i, attempt)
				continue
			}
			// Only take the first 5 questions if we got more
			if len(validQuestions) < 5 {
				validQuestions = append(validQuestions, q)
			}
		}

		if len(validQuestions) >= 3 {
			return validQuestions, nil
		}
	}

	// If all attempts failed, generate emergency fallback questions
	return cs.generateEmergencyFallbackQuestions(symptom, language), lastErr
}

// ValidateAnswer validates a user's answer to a question
func (cs *ConversationServiceImpl) ValidateAnswer(ctx context.Context, question entities.Question, answer string) (bool, string, error) {
	prompt := cs.buildValidationPrompt(question, answer)

	response, err := cs.llmClient.ClassifyTriage(ctx, prompt)
	if err != nil {
		return false, "", fmt.Errorf("failed to validate answer: %w", err)
	}

	isValid, feedback := cs.parseValidationResponse(response)
	return isValid, feedback, nil
}

// GenerateHealthReport creates a structured health report from conversation data
func (cs *ConversationServiceImpl) GenerateHealthReport(ctx context.Context, conversation *entities.Conversation) (*entities.HealthReport, error) {
	prompt := cs.buildReportGenerationPrompt(conversation)

	response, err := cs.llmClient.ClassifyTriage(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate health report: %w", err)
	}

	report, err := cs.parseHealthReportFromResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse health report: %w", err)
	}

	return report, nil
}

// buildQuestionGenerationPrompt creates the prompt for generating follow-up questions
func (cs *ConversationServiceImpl) buildQuestionGenerationPrompt(symptom, language string) string {
	langText := "English"
	if language == "am" {
		langText = "Amharic"
	}

	return fmt.Sprintf(`You are a medical AI assistant helping to gather detailed information about a patient's symptoms. 

Generate exactly 5 targeted follow-up questions for a patient reporting: "%s"

IMPORTANT GUIDELINES:
- Questions must be SPECIFIC to the symptom "%s"
- Tailor questions to gather the most relevant clinical information for this particular symptom
- Consider what healthcare providers would need to know for proper assessment
- Questions should progress logically from basic to more detailed information
- Use clear, simple language appropriate for patients
- Generate questions in %s language

QUESTION CATEGORIES (adapt based on symptom):
1. Duration/Timeline: When did this start? How has it changed over time?
2. Location/Distribution: Where exactly is it? Does it spread or move?
3. Severity/Intensity: How severe is it? Scale of 1-10? Impact on daily activities?
4. Triggers/Patterns: What makes it better/worse? Any patterns you notice?
5. Associated symptoms: Any other symptoms occurring with this?

CRITICAL FORMATTING REQUIREMENTS:
- You MUST return ONLY the JSON array
- Do NOT include any explanatory text before or after the JSON
- Do NOT include markdown formatting (no backticks or code blocks)
- Do NOT include any other text or comments
- The response should start with [ and end with ]
- Keep questions concise (under 100 characters each) to prevent truncation
- Ensure the entire response is complete and properly closed

EXACT JSON FORMAT REQUIRED:
[
  {"id": 1, "text": "Concise question here", "type": "duration", "required": true},
  {"id": 2, "text": "Concise question here", "type": "location", "required": true},
  {"id": 3, "text": "Concise question here", "type": "severity", "required": true},
  {"id": 4, "text": "Concise question here", "type": "associated", "required": true},
  {"id": 5, "text": "Concise question here", "type": "triggers", "required": false}
]

EXAMPLES FOR DIFFERENT SYMPTOMS:
- For headache: Ask about location (front/back/sides), triggers (stress/food/sleep), duration, throbbing vs constant
- For chest pain: Ask about location, radiation, breathing relation, exertion, severity
- For fever: Ask about temperature, other symptoms, duration, pattern, associated chills
- For stomach pain: Ask about location, relation to eating, nausea, bowel changes

REMEMBER: Your response must be ONLY the JSON array for symptom: "%s"`, symptom, symptom, langText, symptom)
}

// buildValidationPrompt creates the prompt for validating answers
func (cs *ConversationServiceImpl) buildValidationPrompt(question entities.Question, answer string) string {
	return fmt.Sprintf(`Validate this answer to a medical question.

Question: %s
Question Type: %s
Answer: %s

Requirements:
- Check if the answer is relevant and informative
- For duration: should include time period (days, hours, etc.)
- For location: should specify body part or area
- For severity: should indicate pain level or intensity
- For history: should mention relevant medical background
- For triggers: should describe what causes or worsens the symptom

Respond with JSON format:
{"valid": true/false, "feedback": "explanation if invalid"}

Validation result:`, question.Text, question.Type, answer)
}

// buildReportGenerationPrompt creates the prompt for generating health reports
func (cs *ConversationServiceImpl) buildReportGenerationPrompt(conversation *entities.Conversation) string {
	// Build context from conversation
	context := fmt.Sprintf("Symptom: %s\nLanguage: %s\n", conversation.Symptom, conversation.Language)

	for i, answer := range conversation.Answers {
		if i < len(conversation.Questions) {
			context += fmt.Sprintf("Q%d: %s\nA%d: %s\n",
				answer.QuestionID,
				conversation.Questions[answer.QuestionID-1].Text,
				answer.QuestionID,
				answer.Text)
		}
	}

	return fmt.Sprintf(`Create a structured health report based on this conversation:

%s

Generate a comprehensive health report in JSON format with the following fields:
- symptom: the main symptom
- duration: how long the symptom has been present
- location: where the symptom is located
- severity: how severe the symptom is
- associated_symptoms: any other symptoms mentioned
- medical_history: relevant medical background
- triggers: what causes or worsens the symptom
- possible_conditions: potential diagnoses
- recommendations: suggested next steps
- urgency_level: GREEN/YELLOW/RED based on severity

Format as JSON object.`, context)
}

// parseQuestionsFromResponse parses questions from LLM response
func (cs *ConversationServiceImpl) parseQuestionsFromResponse(response string) ([]entities.Question, error) {
	// Clean the response and try to find JSON
	response = strings.TrimSpace(response)

	// Try multiple strategies to extract JSON array
	var jsonStr string
	var found bool

	// Strategy 1: Find the first complete JSON array with proper bracket matching
	for i := 0; i < len(response); i++ {
		if response[i] == '[' {
			bracketCount := 0
			braceCount := 0
			inString := false
			escaped := false
			start := i

			for j := i; j < len(response); j++ {
				char := response[j]

				if !inString {
					if char == '[' {
						bracketCount++
					} else if char == ']' {
						bracketCount--
						if bracketCount == 0 {
							jsonStr = response[start : j+1]
							found = true
							break
						}
					} else if char == '{' {
						braceCount++
					} else if char == '}' {
						braceCount--
					} else if char == '"' {
						inString = true
					}
				} else {
					if !escaped && char == '"' {
						inString = false
					}
					escaped = !escaped && char == '\\'
				}
			}
			if found {
				break
			}
		}
	}

	// Strategy 2: Try to repair incomplete JSON if we found an opening bracket but no complete structure
	if !found {
		jsonStart := strings.Index(response, "[")
		if jsonStart != -1 {
			// Check if this looks like a truncated JSON array
			remaining := response[jsonStart:]
			if strings.Contains(remaining, `"id":`) && strings.Contains(remaining, `"text":`) {
				// Use regex-based approach to extract complete JSON objects
				var completeObjects []string

				// Split by object boundaries and process each potential object
				text := strings.ReplaceAll(remaining, "\n", " ")
				text = strings.ReplaceAll(text, "  ", " ")

				// Find complete objects using a simple state machine
				objStart := -1
				braceCount := 0
				inString := false
				escaped := false

				for i, char := range text {
					if !inString {
						if char == '{' {
							if braceCount == 0 {
								objStart = i
							}
							braceCount++
						} else if char == '}' {
							braceCount--
							if braceCount == 0 && objStart != -1 {
								obj := strings.TrimSpace(text[objStart : i+1])
								// Clean up any trailing commas inside the object
								obj = strings.ReplaceAll(obj, ",}", "}")
								completeObjects = append(completeObjects, obj)
								objStart = -1
							}
						} else if char == '"' {
							inString = true
						}
					} else {
						if !escaped && char == '"' {
							inString = false
						}
						escaped = !escaped && char == '\\'
					}
				}

				// If we have at least 3 complete objects, use them
				if len(completeObjects) >= 3 {
					jsonStr = "[" + strings.Join(completeObjects, ",") + "]"
					found = true
				}
			}
		}
	}

	// Strategy 3: Fallback to simple index finding
	if !found {
		jsonStart := strings.Index(response, "[")
		jsonEnd := strings.LastIndex(response, "]")

		if jsonStart != -1 && jsonEnd != -1 && jsonStart < jsonEnd {
			jsonStr = response[jsonStart : jsonEnd+1]
			found = true
		}
	}

	if !found {
		return nil, fmt.Errorf("no valid JSON array found in response: %s", response)
	}

	var questions []entities.Question
	err := json.Unmarshal([]byte(jsonStr), &questions)
	if err != nil {
		// Try to provide more helpful error information
		return nil, fmt.Errorf("failed to unmarshal JSON (length: %d chars): %w\nExtracted JSON: %s", len(jsonStr), err, jsonStr)
	}

	// Validate questions
	if len(questions) == 0 {
		return nil, fmt.Errorf("no questions found in response")
	}

	// Validate that we have the expected number of questions (5), but allow fewer if they're valid
	if len(questions) < 3 {
		return nil, fmt.Errorf("insufficient questions generated: got %d, need at least 3", len(questions))
	}

	return questions, nil
}

// parseValidationResponse parses validation result from LLM response
func (cs *ConversationServiceImpl) parseValidationResponse(response string) (bool, string) {
	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return true, "" // Default to valid if can't parse
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var result struct {
		Valid    bool   `json:"valid"`
		Feedback string `json:"feedback"`
	}

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return true, "" // Default to valid if can't parse
	}

	return result.Valid, result.Feedback
}

// parseHealthReportFromResponse parses health report from LLM response
func (cs *ConversationServiceImpl) parseHealthReportFromResponse(response string) (*entities.HealthReport, error) {
	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		// Create a basic report instead of failing
		return &entities.HealthReport{
			Symptom:            "Unknown",
			Duration:           "Unknown",
			Location:           "Unknown",
			Severity:           "Unknown",
			AssociatedSymptoms: []string{},
			MedicalHistory:     "Unknown",
			Triggers:           "Unknown",
			PossibleConditions: []string{},
			Recommendations:    []string{"Please consult a healthcare provider"},
			UrgencyLevel:       "YELLOW",
			GeneratedAt:        time.Now(),
		}, nil
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var report entities.HealthReport
	err := json.Unmarshal([]byte(jsonStr), &report)
	if err != nil {
		return nil, fmt.Errorf("failed to parse health report: %w", err)
	}

	return &report, nil
}

// generateEmergencyFallbackQuestions generates minimal fallback questions when AI generation fails
func (cs *ConversationServiceImpl) generateEmergencyFallbackQuestions(symptom, language string) []entities.Question {
	// Create basic questions that work for any symptom
	questions := []entities.Question{
		{ID: 1, Text: "When did you first notice this symptom?", Type: "duration", Required: true},
		{ID: 2, Text: "Where exactly do you feel this symptom?", Type: "location", Required: true},
		{ID: 3, Text: "How would you rate the severity from 1-10?", Type: "severity", Required: true},
		{ID: 4, Text: "Are you experiencing any other symptoms?", Type: "associated", Required: true},
		{ID: 5, Text: "What makes this symptom better or worse?", Type: "triggers", Required: false},
	}

	// Translate to Amharic if needed
	if language == "am" {
		questions = []entities.Question{
			{ID: 1, Text: "ይህ ምልክት መቼ ጀመረ?", Type: "duration", Required: true},
			{ID: 2, Text: "ይህ ምልክት የት ይሰማዎታል?", Type: "location", Required: true},
			{ID: 3, Text: "ከ1-10 ምን ያህል ከባድ ነው?", Type: "severity", Required: true},
			{ID: 4, Text: "ሌሎች ምልክቶች አሉዎት?", Type: "associated", Required: true},
			{ID: 5, Text: "ምን ያደርገዋል ይህ ምልክት የተሻለ ወይስ የተባሰ?", Type: "triggers", Required: false},
		}
	}

	return questions
}

// buildSymptomValidationPrompt creates the prompt for validating symptoms
func (cs *ConversationServiceImpl) buildSymptomValidationPrompt(symptom, language string) string {
	langText := "English"
	if language == "am" {
		langText = "Amharic"
	}

	return fmt.Sprintf(`You are a medical AI validator determining if user input represents a legitimate health symptom or concern.

ANALYZE THIS INPUT: "%s"

VALIDATION RULES:

1. REJECT these types of inputs:
   - Greetings: "hello", "hi", "hey", "good morning"
   - Test inputs: "test", "testing", "123", "abc", "xyz"
   - System questions: "what can you do?", "how does this work?", "can you help?"
   - General medical questions: "what causes headaches?", "tell me about diabetes"
   - Nonsense text: gibberish, random characters
   - System capability questions: asking about app features or functionality
   - Requests for general medical information (not personal symptoms)

2. ACCEPT these legitimate symptom descriptions:
   - Basic symptom mentions: "I have a headache", "I had headache today", "my stomach hurts"
   - Symptoms with some context: "I have chest pain when walking", "headache for 2 days"
   - Mental health concerns: "I feel anxious", "I'm depressed", "can't sleep"
   - Physical complaints: "my back hurts", "I'm dizzy", "I have fever"
   - Injury descriptions: "I hurt my ankle", "my arm is swollen"

3. KEY PRINCIPLE: 
   - If someone is describing a PERSONAL health experience (even basic), ACCEPT it
   - The follow-up questions will gather more details - don't require all details upfront
   - Focus on rejecting non-medical inputs, not requiring extensive symptom details initially

4. EMERGENCY CLASSIFICATION:
   - EMERGENCY: severe chest pain, difficulty breathing, severe injuries, suicidal thoughts
   - HIGH: significant symptoms needing evaluation
   - MEDIUM: moderate symptoms
   - LOW: basic symptoms or minor concerns

EXAMPLES TO ACCEPT:
- "I have a headache" ✓
- "I had headache today" ✓  
- "my stomach hurts" ✓
- "I feel sick" ✓
- "chest pain" ✓
- "I'm anxious" ✓

EXAMPLES TO REJECT:
- "hello" ✗
- "test" ✗
- "what can you do?" ✗
- "how to treat fever?" ✗
- "abc123" ✗

Be REASONABLE - accept basic symptom descriptions. The purpose is to filter out non-medical inputs, not to require detailed symptom descriptions upfront.

OUTPUT FORMAT (JSON only):
{
  "valid": true/false,
  "feedback": "Brief explanation in %s",
  "urgency_level": "LOW/MEDIUM/HIGH/EMERGENCY",
  "category": "physical/mental/functional/emergency/invalid"
}

VALIDATE: "%s"`, symptom, langText, symptom)
}

// parseSymptomValidationResponse parses the symptom validation result from LLM response
func (cs *ConversationServiceImpl) parseSymptomValidationResponse(response string) (bool, string, error) {
	// Clean and trim the response
	response = strings.TrimSpace(response)

	// If response is empty or too short, reject by default
	if len(response) < 10 {
		return false, "Invalid or unclear input. Please describe your specific health symptom or concern.", nil
	}

	// Try to extract JSON from the response
	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		// If no valid JSON found, be conservative and reject
		return false, "Please describe your specific health symptom or concern clearly.", fmt.Errorf("no valid JSON found in validation response")
	}

	jsonStr := response[jsonStart : jsonEnd+1]

	var result struct {
		Valid        bool   `json:"valid"`
		Feedback     string `json:"feedback"`
		UrgencyLevel string `json:"urgency_level"`
		Category     string `json:"category"`
	}

	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		// If JSON parsing fails, be conservative and reject
		return false, "Please describe your specific health symptom or concern clearly.", fmt.Errorf("failed to unmarshal validation JSON: %w", err)
	}

	// Additional validation on the parsed result
	if !result.Valid {
		// Provide helpful feedback if available, otherwise use default
		feedback := result.Feedback
		if feedback == "" {
			feedback = "Please describe your specific health symptom or concern clearly."
		}
		return false, feedback, nil
	}

	// For valid symptoms, ensure we have proper categorization
	if result.Category == "invalid" {
		return false, result.Feedback, nil
	}

	// Emergency symptoms should be flagged but still considered valid for processing
	if result.UrgencyLevel == "EMERGENCY" {
		emergencyFeedback := "This appears to be an emergency situation. Please seek immediate medical attention or call emergency services."
		if result.Feedback != "" {
			emergencyFeedback = result.Feedback
		}
		return true, emergencyFeedback, nil
	}

	return true, result.Feedback, nil
}
