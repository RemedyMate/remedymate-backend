package triage

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
	"github.com/RemedyMate/remedymate-backend/infrastructure/content"
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
func (ts *TriageService) ClassifySymptoms(ctx context.Context, input entities.SymptomInput) (*entities.TriageResult, error) {
	if err := ts.ValidateInput(input); err != nil {
		return nil, err
	}

	triageLevel, detectedFlags, err := ts.classifyWithLLM(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("triage classification failed: %w", err)
	}

	result := &entities.TriageResult{
		Level:    triageLevel,
		RedFlags: detectedFlags,
		Message:  ts.getTriageMessage(triageLevel, input.Language),
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
func (ts *TriageService) ValidateInput(input entities.SymptomInput) error {
	if strings.TrimSpace(input.Text) == "" {
		return fmt.Errorf("symptom text cannot be empty")
	}
	if len(input.Text) < 3 {
		return fmt.Errorf("symptom text too short (minimum 3 characters)")
	}
	if len(input.Text) > 500 {
		return fmt.Errorf("symptom text too long (maximum 500 characters)")
	}
	if input.Language != "en" && input.Language != "am" {
		return fmt.Errorf("unsupported language: %s (supported: en, am)", input.Language)
	}
	return nil
}

// performs LLM-based triage classification using data-driven prompts
func (ts *TriageService) classifyWithLLM(ctx context.Context, input entities.SymptomInput) (entities.TriageLevel, []string, error) {
	redFlagPrompt := ts.formatRedFlagRulesForPrompt(input.Language)
	yellowFlagPrompt := ts.formatYellowFlagRulesForPrompt(input.Language)

	prompt := fmt.Sprintf(`
You are a medical triage classifier. Analyze the user input and determine if it describes a medical emergency.
Your ONLY output must be a single JSON object with this exact structure:
{"level": "RED" | "YELLOW" | "GREEN", "flags": ["flag1", "flag2"]}

CRITICAL RED FLAGS (output RED if you detect any of these):
%s

YELLOW FLAGS (output YELLOW if you detect any of these, but no red flags):
%s

GREEN: mild symptoms that don't match any red or yellow flags.

Be conservative - when in doubt, escalate to YELLOW or RED.

User Input (Language: %s): "%s"
	`,
		redFlagPrompt,
		yellowFlagPrompt,
		input.Language,
		input.Text)

	response, err := ts.llmClient.ClassifyTriage(ctx, prompt)
	if err != nil {
		return entities.TriageLevelGreen, nil, fmt.Errorf("LLM API call failed: %w", err)
	}

	var llmResult struct {
		Level string   `json:"level"`
		Flags []string `json:"flags"`
	}
	if err := json.Unmarshal([]byte(response), &llmResult); err != nil {
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
