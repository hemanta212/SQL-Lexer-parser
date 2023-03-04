package main

import (
	"bytes"
	"fmt"

	"github.com/hemanta212/parser-lexer-go/parser"
)

func main() {
	// log.SetOutput(ioutil.Discard)
	reader := bytes.NewBufferString("SELECT COUNT(*) FROM apples")
	parser := parser.NewParser(reader)
	stmt, err := parser.Parse()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n", stmt)
}
