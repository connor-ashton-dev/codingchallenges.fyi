package main

import (
	"bufio"
	"io"
)

func GetFileSize(r io.Reader) int64 {
	content, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return int64(len(content))
}

func GetLineCount(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}

func GetWordCount(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}

func GetCharacterCount(r io.Reader) int {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	count := 0
	for scanner.Scan() {
		count++
	}
	return count
}
