package llm

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	anyllm "github.com/mozilla-ai/any-llm-go"
	"github.com/mozilla-ai/any-llm-go/providers"
	"github.com/mozilla-ai/any-llm-go/providers/anthropic"
	"github.com/verygoodsoftwarenotvirus/journal/internal/journal"
)

const (
	claudeModel      = "claude-3-5-haiku-latest"
	defaultMaxTokens = 1024
)

// ClaudeProvider implements Provider using Anthropic's Claude via any-llm-go.
type ClaudeProvider struct {
	client anyllm.Provider
}

// NewClaudeProvider creates a Claude-backed LLM provider.
// Uses JOURNAL_ANTHROPIC_API_KEY if set, otherwise ANTHROPIC_API_KEY.
func NewClaudeProvider() (*ClaudeProvider, error) {
	opts := []anyllm.Option{}
	if key := os.Getenv("JOURNAL_ANTHROPIC_API_KEY"); key != "" {
		opts = append(opts, anyllm.WithAPIKey(key))
	}
	// Otherwise anthropic.New() will use ANTHROPIC_API_KEY from environment

	client, err := anthropic.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("create anthropic provider: %w", err)
	}

	return &ClaudeProvider{client: client}, nil
}

// AskFollowUpQuestions sends the given journal entries to Claude and returns
// contextual follow-up questions (e.g., about deadlines, plans, events).
func (p *ClaudeProvider) AskFollowUpQuestions(ctx context.Context, entries []*journal.Entry) (string, error) {
	if len(entries) == 0 {
		return "", fmt.Errorf("no entries provided")
	}

	userContent := formatEntriesForPrompt(entries)
	systemPrompt := `You are a thoughtful journaling assistant. You will receive the user's last few journal entries with their dates.

Today's date: ` + time.Now().Format("Monday, January 2, 2006") + `

Your task: Identify commitments, deadlines, plans, events, or topics the user wrote about in the past. Given how much time has passed, generate 2–5 concise, personal follow-up questions.

Examples:
- If they wrote on Monday that a report was due in two days, and today is Friday, ask: "How did the report go?"
- If they mentioned starting a new habit, ask how it's going.
- If they expressed worry about something, ask how it turned out.

Output only the questions, one per line or as a short numbered list. Be warm and conversational.`

	maxTokens := defaultMaxTokens
	resp, err := p.client.Completion(ctx, providers.CompletionParams{
		Model:     claudeModel,
		MaxTokens: &maxTokens,
		Messages: []providers.Message{
			{Role: providers.RoleSystem, Content: systemPrompt},
			{Role: providers.RoleUser, Content: userContent},
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from Claude")
	}

	return resp.Choices[0].Message.ContentString(), nil
}

func formatEntriesForPrompt(entries []*journal.Entry) string {
	var b strings.Builder
	b.WriteString("Journal entries:\n\n")
	for i, e := range entries {
		b.WriteString(fmt.Sprintf("--- Entry %d (%s) ---\n", i+1, e.PublishTime.Format("Monday, January 2, 2006 at 3:04 PM")))
		b.WriteString(e.Content)
		b.WriteString("\n\n")
	}
	return b.String()
}
