package llm

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/verygoodsoftwarenotvirus/journal/internal/journal"
)

// Provider defines the interface for LLM backends that can generate
// follow-up questions based on journal entries.
type Provider interface {
	AskFollowUpQuestions(ctx context.Context, entries []*journal.Entry) (string, error)
}

// CredentialsPresent returns true if LLM credentials are configured for the current provider.
// When false, the app should skip LLM features (e.g., follow-up questions).
func CredentialsPresent() bool {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("JOURNAL_LLM_PROVIDER")))
	if provider == "" {
		provider = "claude"
	}

	switch provider {
	case "claude":
		return os.Getenv("JOURNAL_ANTHROPIC_API_KEY") != "" || os.Getenv("ANTHROPIC_API_KEY") != ""
	default:
		return false
	}
}

// NewProvider returns a Provider based on the JOURNAL_LLM_PROVIDER environment variable.
// Empty or "claude" selects the Claude provider. Future values: "openai", etc.
func NewProvider() (Provider, error) {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("JOURNAL_LLM_PROVIDER")))
	if provider == "" {
		provider = "claude"
	}

	switch provider {
	case "claude":
		return NewClaudeProvider()
	default:
		return nil, fmt.Errorf("unknown LLM provider %q (supported: claude)", provider)
	}
}
