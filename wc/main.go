package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
)

func main() {
	var filePath string
	cFlag := flag.Bool("c", false, "count the number of bytes in a file")
	lFlag := flag.Bool("l", false, "count the number of lines in a file")
	wFlag := flag.Bool("w", false, "count the number of words in a file")
	mFlag := flag.Bool("m", false, "count the number of characters in a file")
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		filePath = args[0]
	} else {
		countFromStdin(*cFlag, *lFlag, *wFlag, *mFlag)
		return
	}

	if !(*cFlag || *lFlag || *wFlag || *mFlag) {
		byteCount, err := countBytes(filePath)
		if err != nil {
			fmt.Println("Error counting bytes:", err)
			os.Exit(1)
		}
		lines, err := countLines(filePath)
		if err != nil {
			fmt.Println("Error counting lines:", err)
			os.Exit(1)
		}
		words, err := countWords(filePath)
		if err != nil {
			fmt.Println("Error counting words:", err)
			os.Exit(1)
		}

		fmt.Printf("%10d %10d %10d %s\n", lines, words, byteCount, filePath)
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

	if *mFlag {
		charCount, err := countCharacters(filePath)
		if err != nil {
			fmt.Println("Error counting characters:", err)
			os.Exit(1)
		}
		fmt.Printf("%10d %s\n", charCount, filePath)
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

func countCharacters(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	charCount := 0

	buffer := make([]byte, 4096)
	for {
		n, err := file.Read(buffer)
		if err != nil && err != bufio.ErrBufferFull {
			break
		}

		charCount += utf8.RuneCount(buffer[:n])
		if err == bufio.ErrBufferFull {
			continue
		}
		break
	}

	return charCount, nil
}

func countFromStdin(cFlag, lFlag, wFlag, mFlag bool) {
	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 {
		filePath := "stdin"
		byteCount, err := countBytesFromReader(os.Stdin)
		if err != nil {
			fmt.Println("Error counting bytes:", err)
			os.Exit(1)
		}
		lines, err := countLinesFromReader(os.Stdin)
		if err != nil {
			fmt.Println("Error counting lines:", err)
			os.Exit(1)
		}
		words, err := countWordsFromReader(os.Stdin)
		if err != nil {
			fmt.Println("Error counting words:", err)
			os.Exit(1)
		}
		charCount, err := countCharactersFromReader(os.Stdin)
		if err != nil {
			fmt.Println("Error counting characters:", err)
			os.Exit(1)
		}
		if !(cFlag || lFlag || wFlag || mFlag) {
			fmt.Printf("%10d %10d %10d %s\n", lines, words, byteCount, filePath)
		}

		if cFlag {
			fmt.Printf("%10d %s\n", byteCount, filePath)
		}
		if lFlag {
			fmt.Printf("%10d %s\n", lines, filePath)
		}
		if wFlag {
			fmt.Printf("%10d %s\n", words, filePath)
		}
		if mFlag {
			fmt.Printf("%10d %s\n", charCount, filePath)
		}
	} else {
		fmt.Println("No input provided.")
	}
}

// Functions to count from readers (files, stdin)
func countBytesFromReader(r io.Reader) (int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func countLinesFromReader(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	return lineCount, scanner.Err()
}

func countWordsFromReader(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	wordCount := 0
	for scanner.Scan() {
		words := splitWords(scanner.Text())
		wordCount += len(words)
	}
	return wordCount, scanner.Err()
}

func countCharactersFromReader(r io.Reader) (int, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	return len([]rune(string(data))), nil
}
