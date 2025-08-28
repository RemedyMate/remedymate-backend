package interfaces

import (
	"context"
)

// LLMClient defines the interface for LLM interactions
type LLMClient interface {
	ClassifyTriage(ctx context.Context, prompt string) (string, error)
}
