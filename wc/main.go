package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	var filePath string
	cFlag := flag.Bool("c", false, "count the number of bytes in a file")
	lFlag := flag.Bool("l", false, "count the number of lines in a file")
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
}

func countBytes(filePath string) (int, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	return len(file), nil
}

func countLines(filepath string) (int, error) {
	file, err := os.Open(filepath)
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
