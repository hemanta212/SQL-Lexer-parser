package main

import (
	"fmt"

	"github.com/hemanta212/parser-lexer-go/parser"
)

func main() {
	// log.SetOutput(ioutil.Discard)

	parser := parser.NewParser(`
	CREATE TABLE oranges
(
       id integer primary key autoincrement,
       name text,
       description text
)
	`)

	results, err := parser.Parse()
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("\nParsed Stmt:\n%v\n", results)
	}
}
