package journal

import "time"

// Entry represents a journal entry
type Entry struct {
	Tags           []string  `json:"tags"`
	PublishTime    time.Time `json:"publish_time"`
	WritingStartTime time.Time `json:"writing_start_time"`
	Content        string    `json:"content"`
}

// NewEntry creates a new journal entry with the current time as both publish and writing start time
func NewEntry(content string, tags []string) *Entry {
	now := time.Now()
	return &Entry{
		Tags:            tags,
		PublishTime:     now,
		WritingStartTime: now,
		Content:         content,
	}
}

// NewEntryWithTimes creates a new journal entry with custom times
func NewEntryWithTimes(content string, tags []string, writingStartTime, publishTime time.Time) *Entry {
	return &Entry{
		Tags:            tags,
		PublishTime:     publishTime,
		WritingStartTime: writingStartTime,
		Content:         content,
	}
}

