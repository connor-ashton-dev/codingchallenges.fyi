package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	args := os.Args[1:]
	var command string
	var fileName string
	var content []byte
	var err error

	if len(args) == 2 {
		command = args[0]
		fileName = args[1]
		content, err = os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
	} else if len(args) == 1 {
		command = "all"
		fileName = args[0]
		content, err = os.ReadFile(fileName)
		if err != nil {
			panic(err)
		}
	} else if len(args) == 0 {
		command = "all"
		content, err = io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Invalid arguments")
		return
	}

	if command == "-c" {
		size := GetFileSize(bytes.NewReader(content))
		fmt.Println(size)

	} else if command == "-l" {
		count := GetLineCount(bytes.NewReader(content))
		fmt.Println(count)

	} else if command == "-w" {
		count := GetWordCount(bytes.NewReader(content))
		fmt.Println(count)

	} else if command == "-m" {
		count := GetCharacterCount(bytes.NewReader(content))
		fmt.Println(count)

	} else if command == "all" {
		size := GetFileSize(bytes.NewReader(content))
		count := GetLineCount(bytes.NewReader(content))
		wordCount := GetWordCount(bytes.NewReader(content))
		fmt.Printf("    %d   %d   %d   %s\n", count, wordCount, size, fileName)
	} else {
		fmt.Println("Invalid command")
	}
}

func getFileReader(fileName string) io.Reader {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	return file
}
