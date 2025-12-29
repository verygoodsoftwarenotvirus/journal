# journal

A simple CLI tool for keeping journal entries in JSON format.

## Features

- Save journal entries as JSON files organized by date (year/month/day)
- Support for tags, publish time, writing start time, and content
- Configurable save location via environment variable
- Interactive and non-interactive modes

## Installation

```bash
go build -o journal .
```

## Usage

### Create a new journal entry

**Interactive mode:**
```bash
journal new -i
```

**With content flag:**
```bash
journal new -c "Today was a great day!" -t "work,personal"
```

**From stdin:**
```bash
echo "My journal entry" | journal new
```

**With tags:**
```bash
journal new -c "My entry" -t tag1,tag2,tag3
```

### Configuration

By default, journal entries are saved to `$HOME/Dropbox/Journal/<year>/<month>/<day>/<timestamp>.json`.

You can override the base path by setting the `JOURNAL_PATH` environment variable:

```bash
export JOURNAL_PATH=/path/to/your/journal
journal new -c "My entry"
```

## Journal Entry Format

Each journal entry is saved as a JSON file with the following structure:

```json
{
  "tags": ["tag1", "tag2"],
  "publish_time": "2024-01-15T10:30:00Z",
  "writing_start_time": "2024-01-15T10:30:00Z",
  "content": "Journal entry content here"
}
```

## Commands

- `journal new` - Create a new journal entry
  - `-c, --content string` - Journal entry content
  - `-t, --tags strings` - Tags for the journal entry (comma-separated)
  - `-i, --interactive` - Interactive mode for entering content