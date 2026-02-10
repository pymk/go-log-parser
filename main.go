package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type LogEntry struct {
	IP           string
	User         string
	Timestamp    time.Time
	Method       string
	Path         string
	Protocol     string
	StatusCode   int
	ResponseSize int
	Referer      string
	UserAgent    string
}

func main() {
	logPath := filepath.Clean("testdata/access.log")

	logs, err := logReader(logPath)
	if err != nil {
		log.Fatal("failed to read log:", err)
	}

	fmt.Printf("Number of logs parsed: %d\n", len(logs))
	if len(logs) > 0 {
		fmt.Println("First line:")
		fmt.Println(logs[0])
	}
}

// logParser parses a single Apache access log line and returns a LogEntry
// struct with the extracted fields. It uses regex to match the expected
// log format and converts timestamp, status code, and response size to
// their appropriate types. Returns an error if the line doesn't match
// the expected format or if type conversion fails.
func logParser(l string) (*LogEntry, error) {
	var entry *LogEntry

	re := regexp.MustCompile(`^(\d+\.\d+\.\d+\.\d+) - (\S+) \[([^\]]+)\] "(\S+) ([^\s]+) (\S+)" (\d+) (\d+) "([^"]*)" "([^"]*)"$`)
	s := re.FindStringSubmatch(l)
	if s == nil {
		return nil, fmt.Errorf("failed to parse line")
	}

	// Datetime conversion
	timeLayout := "02/Jan/2006:15:04:05 -0700"
	timeString := s[3]
	t, err := time.Parse(timeLayout, timeString)
	if err != nil {
		return nil, err
	}

	// Int conversion
	statusInt, err := strconv.Atoi(s[7])
	if err != nil {
		return nil, err
	}
	responseString, err := strconv.Atoi(s[8])
	if err != nil {
		return nil, err
	}

	entry = &LogEntry{
		IP:           s[1],
		User:         s[2],
		Timestamp:    t,
		Method:       s[4],
		Path:         s[5],
		Protocol:     s[6],
		StatusCode:   statusInt,
		ResponseSize: responseString,
		Referer:      s[9],
		UserAgent:    s[10],
	}
	return entry, nil
}

// logReader reads a log file and parses each line into LogEntry structs.
// It takes a file path as input and returns a slice of parsed log entries
// and any error encountered during reading or parsing. Lines that fail to
// parse are skipped with an error message printed to stdout.
func logReader(path string) ([]*LogEntry, error) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the file
	var logs []*LogEntry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		// Parse each line
		parsedln, err := logParser(scanner.Text())
		if err != nil {
			fmt.Println("failed to parse:", err)
			continue
		}

		// Collect the results
		logs = append(logs, parsedln)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}
