package scanner

import (
	"bufio"
	"os"
	"strings"
)

// FileReader handles reading and normalizing file content
type FileReader struct{}

// NewFileReader creates a new file reader
func NewFileReader() *FileReader {
	return &FileReader{}
}

// ReadLines reads all lines from a file, handling different line endings
func (r *FileReader) ReadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	// Handle large lines (up to 1MB per line)
	const maxCapacity = 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()
		// Normalize line endings (remove any remaining \r)
		line = strings.TrimRight(line, "\r")
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// FileExists checks if a file exists and is accessible
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// FileReadable checks if a file is readable
func FileReadable(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	file.Close()
	return true
}
