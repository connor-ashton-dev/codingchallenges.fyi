package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	INVALID TokenType = iota
	OPEN_BRACKET
	CLOSE_BRACKET
	DQUOTE
	COLON
	SQUOTE
	WORD
	COMMA
	BOOL
	NULL
	NUMBER
)

type Token struct {
	token TokenType
	value string
}

type Lexer struct {
	input  string
	tokens []Token
}

func NewLexer(json string) *Lexer {
	cleanedString := removeSpacesAndNewLines(json)
	myLexer := &Lexer{
		input: cleanedString,
	}
	return myLexer
}

func removeSpacesAndNewLines(json string) string {
	// Remove spaces
	noSpaces := strings.ReplaceAll(json, " ", "")
	// Remove new lines
	noNewLines := strings.ReplaceAll(noSpaces, "\n", "")
	// Remove tabs
	noTabs := strings.ReplaceAll(noNewLines, "\t", "")
	// Remove carriage returns (for Windows compatibility)
	noCarriageReturns := strings.ReplaceAll(noTabs, "\r", "")
	return noCarriageReturns
}

func ClassifyToken(tok rune) TokenType {
	switch tok {
	case '{':
		return OPEN_BRACKET
	case '}':
		return CLOSE_BRACKET
	case '"':
		return DQUOTE
	case ':':
		return COLON
	case ',':
		return COMMA
	default:
		if isLetter(tok) {
			return WORD
		} else if isDigit(tok) {
			return NUMBER
		} else {
			return INVALID
		}
	}
}

func isLetter(tok rune) bool {
	return unicode.IsLetter(tok)
}

func isDigit(tok rune) bool {
	return unicode.IsDigit(tok)
}

func (l *Lexer) Read() ([]Token, error) {
	for pos := 0; pos < len(l.input); {
		c := rune(l.input[pos])
		t := ClassifyToken(c)
		value := ""
		var token Token

		// Handle sequences of letters
		if t == WORD {
			for pos < len(l.input) && isLetter(rune(l.input[pos])) || isDigit(rune(l.input[pos])) {
				value += string(l.input[pos])
				pos++
			}
			token = Token{token: classifyString(value), value: value}
		} else if t == NUMBER {
			for pos < len(l.input) && isDigit(rune(l.input[pos])) {
				value += string(l.input[pos])
				pos++
			}
			token = Token{token: NUMBER, value: value}
		} else {
			token = Token{token: t, value: value}
			pos++
		}

		if token.token == INVALID {
			return nil, errors.New(fmt.Sprintf("Invalid token found: %q", c))
		}
		l.tokens = append(l.tokens, token)
	}
	return l.tokens, nil
}

// Additional function to classify multi-character sequences made of letters
func classifyString(value string) TokenType {
	// For now, this just returns LETTER for all sequences,
	// but you can expand this to classify keywords like "true", "false", "null"
	if value == "true" || value == "false" {
		return BOOL
	}
	if value == "null" {
		return NULL
	}

	return WORD
}
