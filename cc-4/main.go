package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var pl = fmt.Println

func parseColumn(s string) ([]int, error) {
	filteredString := s[2:]

	if strings.Contains(filteredString, ",") {
		array := strings.Split(filteredString, ",")
		intArray := make([]int, len(array))
		for i := range array {
			val, err := strconv.Atoi(array[i])
			if err != nil {
				log.Fatal("Error parsing list of -f flags:", err)
			}
			intArray[i] = val
		}
		return intArray, nil
	} else {
		n, err := strconv.Atoi(s[2:])
		if err != nil {
			return nil, err
		}
		return []int{n}, nil
	}
}

func parseDelimeter(s string) string {
	return s[2:]
}

func main() {
	args := os.Args
	var delimeterFlag string = "\t"
	var inputFile string = ""
	var cols []int

	for i, v := range args[1:] {
		if !strings.HasPrefix(v, "-") {
			inputFile = v
		} else {
			if strings.HasPrefix(v, "-f") {
				column, err := parseColumn(v)
				cols = column
				if err != nil {
					log.Fatal("Error parsing flag:", err)
				}
			} else if strings.HasPrefix(v, "-d") {
				delimeterFlag = parseDelimeter(args[i])
			}
		}
	}

	if inputFile == "" {
		log.Fatal("No input file specified")
		return
	}

	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer f.Close()

	content := make([][]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		array := strings.Split(line, delimeterFlag)
		content = append(content, array)
	}

	for i := 0; i < len(content); i++ {
		line := ""
		for _, v := range cols {
			line += content[i][v] + "\t"
		}
		pl(line)
	}
}
