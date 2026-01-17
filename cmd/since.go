package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/verygoodsoftwarenotvirus/journal/internal/journal"
)

var sinceCmd = &cobra.Command{
	Use:   "since",
	Short: "Show how long it's been since your last journal entry",
	Long:  `Shows the time elapsed since your most recent journal entry in human-readable format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		entry, _, err := journal.FindMostRecentEntry()
		if err != nil {
			return fmt.Errorf("failed to find last entry: %w", err)
		}

		now := time.Now()
		diff := now.Sub(entry.PublishTime)

		// Format the time difference in human-readable terms
		humanTime := formatDuration(diff)

		fmt.Printf("Last journal entry: %s ago (%s)\n", humanTime, entry.PublishTime.Format("January 2, 2006 at 3:04 PM"))

		return nil
	},
}

// formatDuration formats a duration in human-readable terms
func formatDuration(d time.Duration) string {
	// Convert to seconds for easier calculation
	seconds := int(d.Seconds())

	// Less than a minute
	if seconds < 60 {
		if seconds < 2 {
			return "just now"
		}
		return fmt.Sprintf("%d seconds", seconds)
	}

	// Convert to minutes
	minutes := seconds / 60
	if minutes < 60 {
		if minutes == 1 {
			return "1 minute"
		}
		return fmt.Sprintf("%d minutes", minutes)
	}

	// Convert to hours
	hours := minutes / 60
	if hours < 24 {
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}

	// Convert to days
	days := hours / 24
	if days < 7 {
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}

	// Convert to weeks
	weeks := days / 7
	if weeks < 4 {
		if weeks == 1 {
			return "1 week"
		}
		return fmt.Sprintf("%d weeks", weeks)
	}

	// Convert to months (approximate)
	months := days / 30
	if months < 12 {
		if months == 1 {
			return "1 month"
		}
		return fmt.Sprintf("%d months", months)
	}

	// Convert to years
	years := days / 365
	if years == 1 {
		return "1 year"
	}
	return fmt.Sprintf("%d years", years)
}

func init() {
	rootCmd.AddCommand(sinceCmd)
}
