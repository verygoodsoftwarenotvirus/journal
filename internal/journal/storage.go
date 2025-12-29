package journal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetJournalPath returns the journal base path, checking environment variable first
func GetJournalPath() string {
	if path := os.Getenv("JOURNAL_PATH"); path != "" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home can't be determined
		return "./Journal"
	}
	return filepath.Join(home, "Dropbox", "Journal")
}

// SaveEntry saves a journal entry to the filesystem
func SaveEntry(entry *Entry) error {
	basePath := GetJournalPath()
	
	// Create directory structure: year/month/day
	now := entry.PublishTime
	year := fmt.Sprintf("%04d", now.Year())
	month := fmt.Sprintf("%02d", int(now.Month()))
	day := fmt.Sprintf("%02d", now.Day())
	
	dirPath := filepath.Join(basePath, year, month, day)
	
	// Create directories if they don't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Create filename with timestamp
	timestamp := now.Format("20060102-150405")
	filename := fmt.Sprintf("%s.json", timestamp)
	filePath := filepath.Join(dirPath, filename)
	
	// Marshal entry to JSON
	jsonData, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// LoadEntry loads a journal entry from a file path
func LoadEntry(filePath string) (*Entry, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal entry: %w", err)
	}
	
	return &entry, nil
}

