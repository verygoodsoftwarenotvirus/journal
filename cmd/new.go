package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/verygoodsoftwarenotvirus/journal/internal/journal"
	"github.com/verygoodsoftwarenotvirus/journal/internal/llm"
)

var (
	contentFlag     string
	tagsFlag        []string
	interactiveFlag bool
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new journal entry",
	Long:  `Create a new journal entry and save it to the journal directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var content string
		var tags []string

		if interactiveFlag {
			// Interactive mode
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter journal content (press Enter twice to finish):\n")
			var lines []string
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("error reading input: %w", err)
				}
				line = strings.TrimRight(line, "\n\r")
				if line == "" && len(lines) > 0 {
					// Check if previous line was also empty (double enter)
					if lines[len(lines)-1] == "" {
						lines = lines[:len(lines)-1] // Remove the last empty line
						break
					}
				}
				lines = append(lines, line)
			}
			content = strings.Join(lines, "\n")

			fmt.Print("Enter tags (comma-separated, or press Enter for none): ")
			tagsInput, _ := reader.ReadString('\n')
			tagsInput = strings.TrimSpace(tagsInput)
			if tagsInput != "" {
				tags = strings.Split(tagsInput, ",")
				for i := range tags {
					tags[i] = strings.TrimSpace(tags[i])
				}
			}
		} else {
			// Non-interactive mode
			if contentFlag == "" {
				// Check if stdin has data (piped input)
				stdinInfo, _ := os.Stdin.Stat()
				if (stdinInfo.Mode() & os.ModeCharDevice) == 0 {
					// Stdin is piped, read from it
					scanner := bufio.NewScanner(os.Stdin)
					var lines []string
					for scanner.Scan() {
						lines = append(lines, scanner.Text())
					}
					if err := scanner.Err(); err != nil {
						return fmt.Errorf("error reading from stdin: %w", err)
					}
					content = strings.Join(lines, "\n")
				} else {
					// Stdin is a terminal, open editor
					initialContent := ""
					if llm.CredentialsPresent() {
						questions, err := fetchFollowUpQuestions()
						if err != nil {
							// Log but don't fail - user can still write without questions
							fmt.Fprintf(os.Stderr, "Note: could not fetch follow-up questions: %v\n", err)
						} else if questions != "" {
							initialContent = "\n\n---\n\nFollow-up questions from your last entries:\n\n" + strings.TrimSpace(questions) + "\n"
						}
					}
					var err error
					content, err = openEditor(initialContent)
					if err != nil {
						return fmt.Errorf("failed to open editor: %w", err)
					}
					content = extractJournalContent(content)
				}
			} else {
				content = contentFlag
			}

			tags = tagsFlag
		}

		if content == "" {
			return fmt.Errorf("journal entry content cannot be empty")
		}

		entry := journal.NewEntry(content, tags)

		if err := journal.SaveEntry(entry); err != nil {
			return fmt.Errorf("failed to save entry: %w", err)
		}

		fmt.Println("Journal entry saved successfully!")
		return nil
	},
}

// fetchFollowUpQuestions fetches follow-up questions from the LLM based on the last 3 entries.
// Returns empty string on any error.
func fetchFollowUpQuestions() (string, error) {
	entries, err := journal.FindLastNEntries(3)
	if err != nil || len(entries) == 0 {
		return "", err
	}

	provider, err := llm.NewProvider()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return provider.AskFollowUpQuestions(ctx, entries)
}

// extractJournalContent returns the journal entry content, stripping any content
// below the "---" separator (e.g., follow-up questions).
func extractJournalContent(edited string) string {
	if idx := strings.Index(edited, "\n---"); idx >= 0 {
		return strings.TrimRight(edited[:idx], "\n\r ")
	}
	return edited
}

// openEditor opens the user's editor (or vim) and returns the content.
// If initialContent is non-empty, it is written to the temp file before opening
// so the user sees it (e.g., follow-up questions at the bottom).
func openEditor(initialContent string) (string, error) {
	// Get editor from environment, default to vim
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// Create temporary file
	tmpfile, err := os.CreateTemp("", "journal-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpfile.Name()
	defer os.Remove(tmpPath) // Clean up temp file

	if initialContent != "" {
		if _, err := tmpfile.WriteString(initialContent); err != nil {
			tmpfile.Close()
			return "", fmt.Errorf("failed to write initial content: %w", err)
		}
	}

	// Close the file so editor can open it
	tmpfile.Close()

	// Build editor command, enabling word wrap for vim/nvim
	args := []string{}
	if strings.Contains(filepath.Base(editor), "vim") {
		args = append(args, "-c", "set wrap linebreak")
	}
	args = append(args, tmpPath)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %w", err)
	}

	// Read the content from the temp file
	contentBytes, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	// Remove trailing newline if present (common when editor saves)
	content := string(contentBytes)
	content = strings.TrimRight(content, "\n\r")

	return content, nil
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVarP(&contentFlag, "content", "c", "", "Journal entry content")
	newCmd.Flags().StringSliceVarP(&tagsFlag, "tags", "t", []string{}, "Tags for the journal entry (comma-separated)")
	newCmd.Flags().BoolVarP(&interactiveFlag, "interactive", "i", false, "Interactive mode for entering content")
}
