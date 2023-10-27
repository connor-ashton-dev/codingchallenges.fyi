package main

import (
	"errors"
)

type Parser struct {
	bracketStack []string
	tokens       []Token
	input        []string
	state        StateType
}

type StateType int

const (
	NORMAL StateType = iota
	INSIDE_STRING
	AFTER_STRING
	INSIDE_OBJECT
)

func (p *Parser) parseTokens() error {
	for _, t := range p.tokens {
		r, err := classifyToken(t)
		if err != nil {
			return err
		}
		p.input = append(p.input, r)
	}
	return nil
}

func classifyToken(t Token) (string, error) {
	switch t.token {
	case OPEN_BRACKET:
		return "{", nil
	case CLOSE_BRACKET:
		return "}", nil
	case COLON:
		return ":", nil
	case DQUOTE:
		return "\"", nil
	case WORD:
		return t.value, nil
	case NUMBER:
		return t.value, nil
	case COMMA:
		return ",", nil
	case NULL:
		return "null", nil
	case BOOL:
		return "bool", nil
	default:
		return "", errors.New("Invalid token")
	}
}

func newParser(rawTokens []Token) *Parser {
	return &Parser{
		tokens: rawTokens,
		state:  NORMAL,
	}
}

func (p *Parser) scanTokens() error {
	for pos := 0; pos < len(p.input); {
		tokenType := p.tokens[pos].token
		switch {
		case tokenType == OPEN_BRACKET:
			p.state = INSIDE_OBJECT
			p.bracketStack = append(p.bracketStack, "{")

		case tokenType == CLOSE_BRACKET:
			if p.state != INSIDE_OBJECT && p.state != AFTER_STRING {
				return errors.New("Closing bracket found outside of object context")
			}
			if len(p.bracketStack) == 0 || p.bracketStack[len(p.bracketStack)-1] != "{" {
				return errors.New("Mismatched brackets")
			}
			p.bracketStack = p.bracketStack[:len(p.bracketStack)-1]

		case tokenType == DQUOTE:
			if p.state == AFTER_STRING {
				return errors.New("Missing a colon or comma")
			}

			if p.state != INSIDE_STRING {
				p.state = INSIDE_STRING
			} else {
				p.state = AFTER_STRING
			}

		case tokenType == WORD:
			if p.state != INSIDE_STRING {
				return errors.New("Word token found outside of string context")
			}

		case tokenType == COLON:
			if p.state != AFTER_STRING {
				return errors.New("Colon token found outside of object context")
			}

			if p.input[pos-1] != "\"" {
				return errors.New("Before colon not valid")
			}

			if !afterColon(p.tokens[pos+1]) {
				return errors.New("After colon not valid")
			}

			p.state = INSIDE_OBJECT

		case tokenType == COMMA:
			if p.commaValid(pos) == false {
				return errors.New("Comma in bad place")
			}
			p.state = INSIDE_OBJECT

		case tokenType == NUMBER:
			if p.tokens[pos-1].token != COLON && p.tokens[pos-1].token != DQUOTE {
				return errors.New("Number in bad place")
			}

		case tokenType == BOOL:
			if p.tokens[pos-1].token != COLON {
				return errors.New("Bool in bad place")
			}
			p.state = AFTER_STRING

		case tokenType == NULL:
			if p.tokens[pos-1].token != COLON {
				return errors.New("null in bad place")
			}
			p.state = AFTER_STRING

		default:
			// handle other cases or ignore
		}
		pos++
	}

	if p.state == INSIDE_STRING {
		return errors.New("Unclosed string literal")
	}
	if len(p.bracketStack) != 0 {
		return errors.New("Mismatched brackets")
	}
	return nil
}

func afterColon(tok Token) bool {
	t := tok.token
	if t != OPEN_BRACKET && t != CLOSE_BRACKET && t != BOOL && t != NULL && t != DQUOTE && t != WORD && t != NUMBER {
		return false
	}
	return true
}

func (p *Parser) commaValid(pos int) bool {
	// all valid tokens before or after comma position
	s := p.tokens[pos-1].token
	if s != OPEN_BRACKET && s != DQUOTE && s != BOOL && s != NULL && s != CLOSE_BRACKET && s != NUMBER {
		return false
	}
	s = p.tokens[pos+1].token
	if s != DQUOTE && s != WORD && s != BOOL && s != NULL {
		return false
	}

	return true
}

func (p *Parser) Parse() error {
	err := p.parseTokens()
	if err != nil {
		return err
	}
	err = p.scanTokens()

	if err != nil {
		return err
	}
	// check if tokens are correct
	if len(p.bracketStack) != 0 {
		return errors.New("Invalid parenthesis")
	}
	return nil
}
