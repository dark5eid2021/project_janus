package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// LogFailure represents a failure event parsed from logs
type LogFailure struct {
	Timestamp string
	Message   string
}

// analyzeLogLine scans the provided the log file for error
// It checks if the line contains "ERROR" and prints the specific message
func analyzeLogLine(line string) {
	if strings.Contains(line, "ERROR") {
		fmt.Printf("Error detected: %s\n", line)
	} else {
		fmt.Printf("Log: %s\n", line)
	}
}

// continuouslyAnalyzedLogStream continually reads log lines from the given reader
// and processes them using the analyzeLogLine function
func continuouslyAnalyzedLogStream(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for {
		// attempt to read the next line
		if scanner.Scan() {
			line := scanner.Text()
			analyzeLogLine(line)
		} else {
			// if no new data is available, wait before retrying
			time.Sleep(100 * time.Millisecond)
		}
		// check for scanner errors - maybe this should be optional?
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading log stream: %v\n", err)
			break
		}
	}
}

func main() {
	// to demonstrate, we open a log file
	// in reality this could be any io.Reader - like a network connection or a pipe
	logFile, err := os.Open("logfile.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open logfile: %v\n", err)
		return
	}
	defer logFile.Close()

	// start the continuous analysis of the log stream
	continuouslyAnalyzedLogStream(logFile)
}
