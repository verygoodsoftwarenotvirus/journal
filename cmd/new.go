package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/verygoodsoftwarenotvirus/journal/internal/journal"
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
					var err error
					content, err = openEditor()
					if err != nil {
						return fmt.Errorf("failed to open editor: %w", err)
					}
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

// openEditor opens the user's editor (or vim) and returns the content
func openEditor() (string, error) {
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

	// Close the file so editor can open it
	tmpfile.Close()

	// Open editor
	cmd := exec.Command(editor, tmpPath)
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
