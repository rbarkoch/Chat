package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

func main() {
	// Load command line arguments.
	clArgs, err := ParseCommandLineArgs(os.Args)
	if err != nil {
		fmt.Println("Error parsing command line arguments:", err)
		flag.Usage()
		os.Exit(1)
	}

	// Print help text if requested.
	if clArgs.Help {
		flag.Usage()
		os.Exit(0)
	}

	configs, err := LoadConfigurationsFromWorkingDirectory()
	if err != nil {
		fmt.Println("Error while loading configurations:", err)
		os.Exit(1)
	}

	userHomePath, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting user home directory:", err)
		os.Exit(1)
	}

	homeConfigPath := path.Join(userHomePath, ".chatconfig")
	_, err = os.Stat(homeConfigPath)
	if err == nil {
		homeConfig, err := LoadConfiguration(homeConfigPath)
		if err != nil {
			fmt.Println("Error while loading configuration from home directory:", err)
			os.Exit(1)
		}

		configs = append(configs, homeConfig)
	}

	// Merge all configurations.
	mergedConfig := MergeConfigurations(configs)

	chatRequest, err := ConstructChatRequest(clArgs, mergedConfig)
	if err != nil {
		fmt.Println("Error constructing chat request:", err)
		os.Exit(1)
	}

	response, err := Chat(chatRequest)
	if err != nil {
		fmt.Println("Error during chat:", err)
		os.Exit(1)
	}

	fmt.Println(response.Text)
}

func ConstructChatRequest(clArgs CommandLineArgs, config Config) (ChatRequest, error) {
	// Start by constructing the chat request from configuration.
	chatRequest := ChatRequest{}

	chatRequest.ApiKey = config.Get("OPENAI_API_KEY")
	chatRequest.Model = config.Get("MODEL")
	chatRequest.SystemPrompt = config.Get("SYSTEM_PROMPT")

	// Override with command line arguments if provided.
	if clArgs.ApiKey != "" {
		chatRequest.ApiKey = clArgs.ApiKey
	}

	if clArgs.Model != "" {
		chatRequest.Model = clArgs.Model
	}

	if clArgs.SystemPrompt != "" {
		chatRequest.SystemPrompt = clArgs.SystemPrompt
	}

	// Set the rest of the values.
	chatRequest.Prompt = clArgs.Prompt

	if clArgs.File != "" {
		fileContents, err := ReadFile(clArgs.File)
		if err != nil {
			return ChatRequest{}, fmt.Errorf("error reading file: %s", err)
		}

		chatRequest.FileContents = fileContents
	}

	// Validate the chat request.
	if chatRequest.Prompt == "" {
		return ChatRequest{}, fmt.Errorf("prompt is required")
	}

	if chatRequest.ApiKey == "" {
		return ChatRequest{}, fmt.Errorf("API key is required")
	}

	return chatRequest, nil
}

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
