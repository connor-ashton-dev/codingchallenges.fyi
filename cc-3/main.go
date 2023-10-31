package main

import (
	"bufio"
	"container/heap"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var pl = fmt.Println

func main() {
	var inputFile = flag.String("i", "", "Path to the input file to be compressed.")
	var outputFile = flag.String("o", "", "Path to the output file where the compressed data will be saved.")
	var methodFlag = flag.String("m", "", "Trigger encoding.")
	flag.Parse()

	if *inputFile == "" || *outputFile == "" {
		fmt.Println("Both -i and -o flags are required.")
		flag.Usage() // This prints out how to use the flags
		os.Exit(1)   // Exit with a non-zero status to indicate an error
	}

	if *methodFlag == "" {
		fmt.Println("Decode/Encode option is necessary")
		flag.Usage() // This prints out how to use the flags
		os.Exit(1)   // Exit with a non-zero status to indicate an error
	}

	if *methodFlag == "encode" {
		f, err := os.Open(*inputFile)
		if err != nil {
			log.Fatal("Error opening file")
		}
		defer f.Close()

		m := make(map[rune]int)
		var builder strings.Builder
		var contents string

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			builder.WriteString(line)
			builder.WriteRune('\n')
			for _, c := range line {
				m[c] += 1
			}
			m['\n']++
		}

		contents = builder.String()

		if err := scanner.Err(); err != nil {
			log.Fatal("Error reading file:", err)
		}

		// WANT TO ENCODE

		// Initialize a priority queue with all nodes.
		nq := make(NodeQueue, len(m))
		i := 0
		for char, freq := range m {
			nq[i] = &Node{char: char, freq: freq}
			i++
		}

		heap.Init(&nq)

		// Build the Huffman tree.
		for len(nq) > 1 {
			node1 := heap.Pop(&nq).(*Node)
			node2 := heap.Pop(&nq).(*Node)

			combinedNode := &Node{
				freq:  node1.freq + node2.freq,
				left:  node1,
				right: node2,
			}
			heap.Push(&nq, combinedNode)
		}

		// The root of the Huffman tree is now the only item in the priority queue.
		root := heap.Pop(&nq).(*Node)
		codes := make(map[rune]string)
		generateHuffmanCodes(root, "", codes)
		encoded := encodeText(contents, codes)
		header := serializeHuffmanCodes(codes, m)
		binaryData := encodeToBinary(encoded)
		// Convert header length to a byte slice
		headerLenBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(headerLenBytes, uint32(len(header)))
		// Combine header length, header itself, encodedLen and binary data
		err = os.WriteFile(*outputFile, append(headerLenBytes, append([]byte(header), binaryData...)...), 0644)
		if err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}
	} else if *methodFlag == "decode" {

		fileContent, err := os.ReadFile(*inputFile)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		headerLength := binary.BigEndian.Uint32(fileContent[:4])
		header := fileContent[4 : 4+headerLength]
		encodedContent := fileContent[4+headerLength:]
		encodedText := binaryToEncodedText(encodedContent)
		root := rebuildTree(string(header))
		decodedText := decode(encodedText, root)
		err = os.WriteFile(*outputFile, []byte(decodedText), 0644)
	}
}
