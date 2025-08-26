package content

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/RemedyMate/remedymate-backend/domain/entities"
	"github.com/RemedyMate/remedymate-backend/domain/interfaces"
)

type ContentService struct {
	approvedBlocks []entities.ApprovedBlock
	redFlagRules   []entities.RedFlagRule
	dataPath       string
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

	return nil
}

// reloads content from files
func (cs *ContentService) ReloadContent() error {
	return cs.LoadContent()
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
			return nil, fmt.Errorf("language '%s' not available for topic '%s'", language, topicKey)
		}
	}
	return nil, fmt.Errorf("topic '%s' not found", topicKey)
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
	for _, rule := range cs.redFlagRules {
		if rule.Level == entities.TriageLevelYellow {
			yellowRules = append(yellowRules, rule)
		}
	}
	return yellowRules
}
