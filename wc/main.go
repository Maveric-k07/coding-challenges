package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"unicode"
)

func main() {
	var filePath string
	cFlag := flag.Bool("c", false, "count the number of bytes in a file")
	lFlag := flag.Bool("l", false, "count the number of lines in a file")
	wFlag := flag.Bool("w", false, "count the number of words in a file")
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		filePath = args[0]
	} else {
		fmt.Println("File path is required")
		os.Exit(1)
	}

	if *cFlag {
		byteCount, err := countBytes(filePath)
		if err != nil {
			fmt.Println("Error counting bytes:", err)
			os.Exit(1)
		}
		fmt.Printf("%10d %s\n", byteCount, filePath)
	}

	if *lFlag {
		lines, err := countLines(filePath)
		if err != nil {
			fmt.Println("Error counting lines: ", err)
			os.Exit(1)
		}
		fmt.Printf("%10d %s\n", lines, filePath)
	}

	if *wFlag {
		words, err := countWords(filePath)
		if err != nil {
			fmt.Println("Error counting words:", err)
			os.Exit(1)
		}
		fmt.Printf("%10d %s\n", words, filePath)
	}
}

func countBytes(filePath string) (int, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	return len(file), nil
}

func countLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return lineCount, nil
}

func countWords(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wordCount := 0

	for scanner.Scan() {
		words := splitWords(scanner.Text())
		wordCount += len(words)
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return wordCount, nil
}

func splitWords(text string) []string {
	words := make([]string, 0)
	inWord := false

	for _, char := range text {
		if unicode.IsSpace(char) {
			inWord = false
		} else {
			if !inWord {
				words = append(words, "")
			}
			words[len(words)-1] += string(char)
			inWord = true
		}
	}

	return words
}
