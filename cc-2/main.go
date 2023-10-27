package main

import "fmt"

func main() {
	myJson := `
{
  "key": "value",
  "keyn": 101,
  "keyo": {
    "inner key": "inner value",
	"newkey": {
		"nested": "value",
		"boolean": true
	}
  }
}
`
	l := NewLexer(myJson)
	tokens, err := l.Read()
	if err != nil {
		fmt.Printf("Error while tokenizing: %s\n", err)
		return
	}

	p := newParser(tokens)
	err = p.Parse()
	if err != nil {
		fmt.Printf("Error while parsing: %s\n", err)
		return
	}

	fmt.Println("VALID")
}
