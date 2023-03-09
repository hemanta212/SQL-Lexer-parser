package main

import (
	"fmt"

	"github.com/hemanta212/parser-lexer-go/parser"
)

func main() {
	// log.SetOutput(ioutil.Discard)
	// reader := bytes.NewBufferString("SELECT name, color FROM apples")
	// "CREATE TABLE some_table(\nname string,\ncast varchar,\nother int)")
	// CREATE TABLE oranges(
	//         id integer primary key autoincrement,
	//         name text,
	//         description text
	// )
	// `)

	parser := parser.NewParser("Select count( hello , hi hello  )  from appels")
	results, err := parser.Parse()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%v\n", results)
}
