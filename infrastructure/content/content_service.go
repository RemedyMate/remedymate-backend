package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	derrors "remedymate-backend/domain/AppError"
	"remedymate-backend/domain/entities"
	"remedymate-backend/domain/interfaces"
)

type ContentService struct {
	approvedBlocks  []entities.ApprovedBlock
	redFlagRules    []entities.RedFlagRule
	yellowFlagRules []entities.RedFlagRule
	dataPath        string
}

// NewContentService creates a new content service instance
func NewContentService(dataPath string) interfaces.ContentService {
	service := &ContentService{
		dataPath: dataPath,
	}

	// Load content on initialization
	if err := service.LoadContent(); err != nil {
		fmt.Printf("❌ Failed to load content: %v\n", err)
	}

	return service
}

// loads approved blocks from JSON file
func (cs *ContentService) LoadContent() error {
	// Load approved blocks
	blocksPath := filepath.Join(cs.dataPath, "approved_block.json")
	blocksData, err := os.ReadFile(blocksPath)
	if err != nil {
		return fmt.Errorf("failed to read approved blocks file: %w", err)
	}

	var blocks []entities.ApprovedBlock
	if err := json.Unmarshal(blocksData, &blocks); err != nil {
		return fmt.Errorf("failed to parse approved blocks JSON: %w", err)
	}

	cs.approvedBlocks = blocks

	// Load red flag rules
	if err := cs.loadRedFlagRules(); err != nil {
		return fmt.Errorf("failed to load red flag rules: %w", err)
	}

	// Load yellow flag rules
	if err := cs.loadYellowFlagRules(); err != nil {
		return fmt.Errorf("failed to load yellow flag rules: %w", err)
	}

	return nil
}

// returns all approved blocks
func (cs *ContentService) GetApprovedBlocks() ([]entities.ApprovedBlock, error) {
	if len(cs.approvedBlocks) == 0 {
		return nil, fmt.Errorf("no approved blocks loaded")
	}
	return cs.approvedBlocks, nil
}

// returns content for a specific topic and language
func (cs *ContentService) GetContentByTopic(topicKey, language string) (*entities.ContentTranslation, error) {
	for _, block := range cs.approvedBlocks {
		if block.TopicKey == topicKey {
			if content, exists := block.Translations[language]; exists {
				return &content, nil
			}
			return nil, derrors.ErrLanguageNotAvailable
		}
	}
	return nil, derrors.ErrTopicNotFound
}

// loads red flag rules from JSON file only
func (cs *ContentService) loadRedFlagRules() error {
	rulesPath := filepath.Join(cs.dataPath, "red_flag_rules.json")
	rulesData, err := os.ReadFile(rulesPath)
	if err != nil {
		return fmt.Errorf("failed to read red flag rules file: %w", err)
	}

	var rules []entities.RedFlagRule
	if err := json.Unmarshal(rulesData, &rules); err != nil {
		return fmt.Errorf("failed to parse red flag rules JSON: %w", err)
	}

	cs.redFlagRules = rules
	fmt.Printf("✅ Loaded %d red flag rules from %s\n", len(rules), rulesPath)
	return nil
}

// load yellow flag rules from separate file
func (cs *ContentService) loadYellowFlagRules() error {
	rulesPath := filepath.Join(cs.dataPath, "yellow_flag_rules.json")
	rulesData, err := os.ReadFile(rulesPath)
	if err != nil {
		return fmt.Errorf("failed to read yellow flag rules file: %w", err)
	}

	var rules []entities.RedFlagRule
	if err := json.Unmarshal(rulesData, &rules); err != nil {
		return fmt.Errorf("failed to parse yellow flag rules JSON: %w", err)
	}

	cs.yellowFlagRules = rules
	fmt.Printf("✅ Loaded %d yellow flag rules from %s\n", len(rules), rulesPath)
	return nil
}

// returns the red flag rules for triage
func (cs *ContentService) GetRedFlagRules() []entities.RedFlagRule {
	return cs.redFlagRules
}

// returns only red flag rules (level RED)
func (cs *ContentService) GetRedFlagRulesOnly() []entities.RedFlagRule {
	var redRules []entities.RedFlagRule
	for _, rule := range cs.redFlagRules {
		if rule.Level == entities.TriageLevelRed {
			redRules = append(redRules, rule)
		}
	}
	return redRules
}

// returns only yellow flag rules (level YELLOW)
func (cs *ContentService) GetYellowFlagRules() []entities.RedFlagRule {
	var yellowRules []entities.RedFlagRule
	for _, rule := range cs.yellowFlagRules {
		if rule.Level == entities.TriageLevelYellow {
			yellowRules = append(yellowRules, rule)
		}
	}
	return yellowRules
}
