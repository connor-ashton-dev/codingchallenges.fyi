package main

import (
	"container/heap"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func generateHuffmanCodes(node *Node, code string, codes map[rune]string) {
	if node == nil {
		return
	}

	// If it's a leaf node (has a character), store its code
	if node.left == nil && node.right == nil {
		codes[node.char] = code
		return
	}

	// Traverse left and right children
	generateHuffmanCodes(node.left, code+"0", codes)
	generateHuffmanCodes(node.right, code+"1", codes)
}

func encodeText(input string, codes map[rune]string) string {
	var builder strings.Builder
	for _, char := range input {
		code, exists := codes[char]
		if !exists {
			// Handle missing codes. This should ideally not happen if the Huffman tree is built correctly.
			log.Fatalf("Huffman code missing for character: %v", char)
		}
		builder.WriteString(code)
	}
	return builder.String()
}

func decode(encodedText string, root *Node) string {
	var decoded strings.Builder
	currentNode := root

	for _, bit := range encodedText {
		if bit == '0' {
			currentNode = currentNode.left
		} else {
			currentNode = currentNode.right
		}

		if currentNode.left == nil && currentNode.right == nil {
			decoded.WriteRune(currentNode.char)
			currentNode = root
		}
	}

	return decoded.String()
}

func encodeToBinary(encoded string) []byte {
	byteCount := (len(encoded) + 7) / 8 // Calculate the necessary bytes
	result := make([]byte, byteCount)
	for i := 0; i < len(encoded); i++ {
		if encoded[i] == '1' {
			result[i/8] |= 1 << (7 - i%8) // Set the specific bit to 1
		}
	}
	return result
}
func binaryToEncodedText(data []byte) string {
	var result strings.Builder
	for _, b := range data {
		for i := 0; i < 8; i++ {
			if (b>>uint(7-i))&1 == 1 {
				result.WriteByte('1')
			} else {
				result.WriteByte('0')
			}
		}
	}
	return result.String()
}
func serializeHuffmanCodes(codes map[rune]string, freqMap map[rune]int) string {
	var builder strings.Builder
	for char, code := range codes {
		builder.WriteString(fmt.Sprintf("%d|%s|%d ", char, code, freqMap[char]))
	}
	return builder.String()
}

func rebuildTree(header string) *Node {
	nodes := strings.Split(header, " ") // Split by space to get individual character codes
	freqMap := make(map[rune]int)
	for _, node := range nodes {
		if len(node) == 0 {
			continue
		}

		parts := strings.Split(node, "|")
		charValue, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Println("Error converting char value:", parts[0])
			continue
		}
		char := rune(charValue)
		frequency, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Println("Error converting frequency value:", parts[2])
			continue
		}
		freqMap[char] = frequency
	}

	// The logic below builds the Huffman tree just like in the main encoding part.
	nq := make(NodeQueue, len(freqMap))
	i := 0
	for char, freq := range freqMap {
		nq[i] = &Node{char: char, freq: freq}
		i++
	}

	heap.Init(&nq)

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

	root := heap.Pop(&nq).(*Node)
	return root
}
