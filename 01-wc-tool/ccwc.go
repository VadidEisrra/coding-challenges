package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode/utf8"
)

var p = fmt.Println
var legalOptions = []string{"c", "l", "w", "m"}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, char := range slice {
		if _, present := seen[char]; !present {
			seen[char] = true
			result = append(result, char)
		}
	}
	return result
}

func normalizeOptions(options []string) []string {
	normalized := []string{}

	for _, option := range options {
		stripped := strings.TrimPrefix(option, "-")
		if len(stripped) > 1 {
			for _, char := range stripped {
				normalized = append(normalized, string(char))
			}
		} else {
			normalized = append(normalized, stripped)
		}

	}
	return removeDuplicates(normalized)
}

func validateOptions(options []string) []string {
	for _, flag := range options {
		if !slices.Contains(legalOptions, flag) {
			p("ccwc: illegal option", flag)
			p("usage: ccwc [-clmw] [file ..]")
			os.Exit(1)
		}
	}
	return options
}

func countBytes(data []byte) int {
	byteCount := len(data)
	return byteCount
}

func countLines(data []byte) int {
	lineCount := 0
	for _, char := range string(data) {
		if char == '\n' {
			lineCount++
		}
	}
	return lineCount
}

func countWords(data []byte) int {
	wordString := string(data)
	wordSlice := strings.Fields(wordString)
	wordCount := len(wordSlice)
	return wordCount
}

func countChars(data []byte) int {
	charCount := utf8.RuneCount(data)
	return charCount
}

func counter(options []string, fileContents []byte) map[string]int {
	result := make(map[string]int)

	if len(options) == 0 {
		result["newLines"] = countLines(fileContents)
		result["words"] = countWords(fileContents)
		result["bytes"] = countBytes(fileContents)

		return result
	}
	for _, flag := range options {
		if flag == "c" {
			result["bytes"] = countBytes(fileContents)
		}
		if flag == "l" {
			result["newLines"] = countLines(fileContents)
		}
		if flag == "w" {
			result["words"] = countWords(fileContents)
		}
		if flag == "m" {
			result["bytes"] = countChars(fileContents)
		}
	}
	return result
}

func main() {
	args := os.Args[1:]
	commandFlags := []string{}
	fileNames := []string{}

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			commandFlags = append(commandFlags, arg)
		}
		if !strings.HasPrefix(arg, "-") {
			fileNames = append(fileNames, arg)
		}
	}

	options := normalizeOptions(commandFlags)
	validateOptions(options)

	for _, file := range fileNames {
		dat, err := os.ReadFile(file)
		check(err)

		result := counter(options, dat)
		p(result)
	}
}
