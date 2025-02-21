package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// This service reads two files (a log file and diff containing code changes/commits), looks for error
// indicators in the logs, scans for warning keywords in the code changes, and then prints an overall risk
// level based on some simle heuristics

// future iterations needed - reading a live log tail and reading the code changes directly from git...
// not relying on files...can the diff be read as part of a pre merge check?

// LogFailure represents a failure event parsed from the logs
type LogFailure struct {
	Timestamp string
	Message   string
}

// analyzeLogs scans the provided log file for error & failure messages
// it employs a regex to capture the lines that include ERROR, FATAL, PANIC
func analyzeLogs(logFile string) ([]LogFailure, error) {
	file, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var failures []LogFailure
	scanner := bufio.NewScanner(file)
	// an example log line: "[2025-02-20 12:00:00] ERROR: Something went wrong"
	re := regexp.MustCompile(`\[(.*?)\]\s+(ERROR|FATAL|PANIC):\s+(.*)`)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 4 {
			failures = append(failures, LogFailure{
				Timestamp: matches[1],
				Message:   matches[3],
			})
		}
	}
	return failures, scanner.Err()
}

// analyzeCodeChanges scans a fuile of code changes/comments and assigns a risk score
// it counts occurrences of common warning keywords such as "TODO", "FIXME", "hack" and "unsafe"
func analyzeCodeChanges(changesFile string) (int, error) {
	file, err := os.Open(changesFile)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	riskScore := 0
	scanner := bufio.NewScanner(file)
	keywords := []string{"TODO", "FIXME", "hack", "unsafe"}
	for scanner.Scan() {
		line := scanner.Text()
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				riskScore++
			}
		}
	}
	return riskScore, scanner.Err()
}

// assessRisk computes overall risk level based on the number of log failures and the code change risk score
func assessRisk(failures []LogFailure, codeRisk int) string {
	numFailures := len(failures)
	riskLevel := ""

	// here is a simple heuristic:
	// high risk: > 10 failures or code risk score > 20
	// medium risk: > 5 failures or code risk score > 10
	// low risk: less than the above
	if numFailures > 10 || codeRisk > 20 {
		riskLevel = "High"
	} else if numFailures > 5 || codeRisk > 10 {
		riskLevel = "Medium"
	} else {
		riskLevel = "Low"
	}

	return fmt.Sprintf("Risk Level: %s (Failures: %d, Code Risk Score: %d)", riskLevel, numFailures, codeRisk)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: riskassessor <log_file> <code_changes_file>")
		os.Exit(1)
	}
	logFile := os.Args[1]
	changesFile := os.Args[2]

	failures, err := analyzeLogs(logFile)
	if err != nil {
		fmt.Printf("Error analyzing logs: %v\n", err)
		os.Exit(1)
	}

	codeRisk, err := analyzeCodeChanges(changesFile)
	if err != nil {
		fmt.Printf("Error analyzing code changes: %v\n", err)
		os.Exit(1)
	}

	assessment := assessRisk(failures, codeRisk)
	fmt.Println(assessment)
}
