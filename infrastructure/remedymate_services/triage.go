package remedymate_services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
	"remedymate-backend/infrastructure/content"
)

type TriageService struct {
	contentService interfaces.ContentService
	llmClient      interfaces.LLMClient
}

func NewTriageService(contentService interfaces.ContentService, llmClient interfaces.LLMClient) interfaces.TriageService {
	return &TriageService{
		contentService: contentService,
		llmClient:      llmClient,
	}
}

// performs LLM-powered triage classification only (no fallback)
func (ts *TriageService) ClassifySymptoms(ctx context.Context, textInput, lang string) (*entities.TriageResult, error) {
	if err := ts.ValidateInput(textInput, lang); err != nil {
		return nil, err
	}

	triageLevel, detectedFlags, err := ts.classifyWithLLM(ctx, textInput, lang)

	if err != nil {
		return nil, fmt.Errorf("triage classification failed: %w", err)
	}

	// Handle unclear input specially
	if len(detectedFlags) > 0 && detectedFlags[0] == "unclear_input" {
		result := &entities.TriageResult{
			Level:    entities.TriageLevelGreen, // Use green level but with clarification message
			RedFlags: []string{},
			Message:  ts.getClarificationMessage(lang),
		}
		return result, nil
	}

	result := &entities.TriageResult{
		Level:    triageLevel,
		RedFlags: detectedFlags,
		Message:  ts.getTriageMessage(triageLevel, lang),
	}

	return result, nil
}

// returns appropriate message based on triage level
func (ts *TriageService) getTriageMessage(level entities.TriageLevel, language string) string {
	switch level {
	case entities.TriageLevelRed:
		return ts.getRedFlagMessage(language)
	case entities.TriageLevelYellow:
		return ts.getYellowFlagMessage(language)
	case entities.TriageLevelGreen:
		return ts.getGreenFlagMessage(language)
	default:
		return ts.getGreenFlagMessage(language)
	}
}

// validates the symptom input
func (ts *TriageService) ValidateInput(inputText, lang string) error {
	if strings.TrimSpace(inputText) == "" {
		return fmt.Errorf("symptom text cannot be empty")
	}
	if len(inputText) < 3 {
		return fmt.Errorf("symptom text too short (minimum 3 characters)")
	}
	if len(inputText) > 500 {
		return fmt.Errorf("symptom text too long (maximum 500 characters)")
	}
	if lang != "en" && lang != "am" {
		return fmt.Errorf("unsupported language: %s (supported: en, am)", lang)
	}
	return nil
}

// performs LLM-based triage classification using data-driven prompts
func (ts *TriageService) classifyWithLLM(ctx context.Context, inputText, lang string) (entities.TriageLevel, []string, error) {
	redFlagPrompt := ts.formatRedFlagRulesForPrompt(lang)
	yellowFlagPrompt := ts.formatYellowFlagRulesForPrompt(lang)
	approvedTopicsPrompt := ts.formatApprovedTopicsForPrompt(lang)

	prompt := fmt.Sprintf(`
You are a medical triage classifier. Analyze the user input and determine if it describes a medical emergency.
Your ONLY output must be a single JSON object with this exact structure:
{"level": "RED" | "YELLOW" | "GREEN" | "UNCLEAR", "flags": ["flag1", "flag2"]}

CRITICAL RED FLAGS (output RED if you detect any of these):
%s

YELLOW FLAGS (output YELLOW if you detect any of these, but no red flags):
%s

GREEN FLAGS (output GREEN only if the symptom matches one of these approved topics):
%s

UNCLEAR: If the input doesn't clearly match any of the above categories, or if you cannot understand what the user is describing.

Be conservative - when in doubt about red/yellow flags, escalate to YELLOW or RED.
Only return GREEN if the symptom clearly matches one of the approved topics.

User Input (Language: %s): "%s"
	`,
		redFlagPrompt,
		yellowFlagPrompt,
		approvedTopicsPrompt,
		lang,
		inputText)

	response, err := ts.llmClient.ClassifyTriage(ctx, prompt)
	if err != nil {
		return entities.TriageLevelGreen, nil, fmt.Errorf("LLM API call failed: %w", err)
	}

	// Clean the response by removing markdown formatting
	cleanedResponse := strings.TrimSpace(response)
	cleanedResponse = strings.TrimPrefix(cleanedResponse, "```json")
	cleanedResponse = strings.TrimPrefix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSuffix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSpace(cleanedResponse)

	var llmResult struct {
		Level string   `json:"level"`
		Flags []string `json:"flags"`
	}
	if err := json.Unmarshal([]byte(cleanedResponse), &llmResult); err != nil {
		return entities.TriageLevelGreen, nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	var level entities.TriageLevel
	switch llmResult.Level {
	case "RED":
		level = entities.TriageLevelRed
	case "YELLOW":
		level = entities.TriageLevelYellow
	case "GREEN":
		level = entities.TriageLevelGreen
	case "UNCLEAR":
		// Return a special response asking for clarification
		return entities.TriageLevelGreen, []string{"unclear_input"}, nil
	default:
		return entities.TriageLevelGreen, nil, fmt.Errorf("LLM returned invalid triage level: %s", llmResult.Level)
	}

	return level, llmResult.Flags, nil
}

// formats red flag rules for inclusion in LLM prompts
func (ts *TriageService) formatRedFlagRulesForPrompt(language string) string {
	redFlagRules := ts.getRedFlagRulesOnly()
	var ruleDescriptions []string

	for _, rule := range redFlagRules {
		if rule.Language == language {
			desc := fmt.Sprintf("%s: %s", strings.Join(rule.Keywords, ", "), rule.Description)
			ruleDescriptions = append(ruleDescriptions, desc)
		}
	}

	return strings.Join(ruleDescriptions, "\n")
}

// formats yellow flag rules for inclusion in LLM prompts
func (ts *TriageService) formatYellowFlagRulesForPrompt(language string) string {
	yellowFlagRules := ts.getYellowFlagRulesOnly()
	var ruleDescriptions []string

	for _, rule := range yellowFlagRules {
		if rule.Language == language {
			desc := fmt.Sprintf("%s: %s", strings.Join(rule.Keywords, ", "), rule.Description)
			ruleDescriptions = append(ruleDescriptions, desc)
		}
	}

	return strings.Join(ruleDescriptions, "\n")
}

// formatApprovedTopicsForPrompt formats approved topics for inclusion in LLM prompts
func (ts *TriageService) formatApprovedTopicsForPrompt(language string) string {
	approvedBlocks, err := ts.contentService.GetApprovedBlocks()
	if err != nil {
		return ""
	}

	var topicDescriptions []string
	for _, block := range approvedBlocks {
		if translation, exists := block.Translations[language]; exists {
			// Create a description based on the topic key and self-care items
			desc := fmt.Sprintf("%s: %s", block.TopicKey, strings.Join(translation.SelfCare[:min(2, len(translation.SelfCare))], ", "))
			topicDescriptions = append(topicDescriptions, desc)
		}
	}

	return strings.Join(topicDescriptions, "\n")
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// gets only red flag rules from content service
func (ts *TriageService) getRedFlagRulesOnly() []entities.RedFlagRule {
	if cs, ok := ts.contentService.(*content.ContentService); ok {
		return cs.GetRedFlagRulesOnly()
	}
	return nil
}

// gets only yellow flag rules from content service
func (ts *TriageService) getYellowFlagRulesOnly() []entities.RedFlagRule {
	if cs, ok := ts.contentService.(*content.ContentService); ok {
		return cs.GetYellowFlagRules()
	}
	return nil
}

// Message helpers
func (ts *TriageService) getRedFlagMessage(language string) string {
	if language == "am" {
		return "ወዲያውኑ የህክምና እርዳታ ይፈልጉ። ወደ ቅርብ ሆስፒታል ወይም የድንገተኛ ጊዜ አገልግሎት ይሂዱ።"
	}
	return "Seek emergency care immediately. Go to the nearest hospital or emergency service."
}

func (ts *TriageService) getYellowFlagMessage(language string) string {
	if language == "am" {
		return "ምልክቶችዎን በጥንቃቄ ይከታተሉ። ካልተሻሻለ ወይም ከባሰ የህክምና ባለሙያ ያማክሩ።"
	}
	return "Monitor your symptoms closely. Consult a healthcare professional if they don't improve or worsen."
}

func (ts *TriageService) getGreenFlagMessage(language string) string {
	if language == "am" {
		return "ምልክቶችዎ ቀላል ሊሆኑ ይችላሉ። የራስ እንክብካቤ ምክሮችን ይከተሉ።"
	}
	return "Your symptoms appear to be mild. Follow self-care recommendations."
}

func (ts *TriageService) getClarificationMessage(language string) string {
	if language == "am" {
		return "ይቅርታ፣ የገለጹልኝን ምልክቶች በግልጽ ለመረዳት አልቻልኩም። ምልክቶችዎን በሌላ መንገድ ሊያስረዱኝ ወይም ተጨማሪ ዝርዝር ሊሰጡኝ ይችላሉ?"
	}
	return "I'm sorry, I couldn't clearly understand the symptoms you described. Could you please explain your symptoms in a different way or provide more details?"
}
