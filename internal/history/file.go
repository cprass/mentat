package history

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type FileStore struct {
	path string
	file *os.File
}

func NewFileStore(path string) (*FileStore, error) {
	logPath := filepath.Join(path, "reviews.log")
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &FileStore{
		path: logPath,
		file: file,
	}, nil
}

func (f *FileStore) Append(event *ReviewEvent) error {
	file, err := os.OpenFile(f.path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, event.String())
	if err != nil {
		return err
	}
	return nil
}

func (f *FileStore) LoadAll() ([]*ReviewEvent, error) {
	file, err := os.Open(f.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []*ReviewEvent{}, nil // Empty history
		}
		return nil, err
	}
	defer file.Close()

	var events []*ReviewEvent
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // skip empty lines
		}

		event, err := NewReviewEventFromString(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %q: %w", line, err)
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (f *FileStore) Close() error {
	return f.file.Close()
}
