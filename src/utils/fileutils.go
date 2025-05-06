package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadFile supports three selectors on the given file path:
//  1. <path>                     – entire file
//  2. <path>:<line>             – single 1‑based line
//  3. <path>:<start>-<end>      – inclusive line range
//
// The function returns the selected content as a single string. Whitespace in the
// original file is preserved (including newline characters for ranged selections).
func ReadFile(filePath string) (string, error) {
	// Split path from optional selector.
	parts := strings.SplitN(filePath, ":", 2)
	basePath := parts[0]

	// --- Case 1: no selector – read and return entire file --------------------
	if len(parts) == 1 {
		bytes, err := os.ReadFile(basePath)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}

	selector := parts[1]
	// Decide whether selector is a single line or a range.
	if !strings.Contains(selector, "-") {
		// --- Case 2: single line ------------------------------------------------
		lineNum, err := strconv.Atoi(selector)
		if err != nil || lineNum <= 0 {
			return "", fmt.Errorf("invalid line number in path: %s", filePath)
		}
		return readSingleLine(basePath, lineNum)
	}

	// --- Case 3: range ----------------------------------------------------------
	bounds := strings.SplitN(selector, "-", 2)
	if len(bounds) != 2 {
		return "", fmt.Errorf("invalid line range in path: %s", filePath)
	}
	start, err1 := strconv.Atoi(bounds[0])
	end, err2 := strconv.Atoi(bounds[1])
	if err1 != nil || err2 != nil || start <= 0 || end < start {
		return "", fmt.Errorf("invalid line range in path: %s", filePath)
	}
	return readLineRange(basePath, start, end)
}

// readSingleLine returns the specified 1‑based line from the file.
func readSingleLine(path string, lineNum int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	current := 0
	for scanner.Scan() {
		current++
		if current == lineNum {
			return scanner.Text(), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("line %d out of range in %s", lineNum, path)
}

// readLineRange returns lines [start, end] (inclusive) joined by newline.
func readLineRange(path string, start, end int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	current := 0
	var lines []string
	for scanner.Scan() {
		current++
		if current < start {
			continue
		}
		if current > end {
			break
		}
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "", fmt.Errorf("line range %d-%d out of range in %s", start, end, path)
	}
	return strings.Join(lines, "\n"), nil
}
