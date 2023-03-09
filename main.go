package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/hemanta212/parser-lexer-go/parser"
)

func main() {
	log.SetOutput(ioutil.Discard)
	// reader := bytes.NewBufferString("SELECT name, color FROM apples")
	// "CREATE TABLE some_table(\nname string,\ncast varchar,\nother int)")
	// CREATE TABLE oranges(
	//         id integer primary key autoincrement,
	//         name text,
	//         description text
	// )
	// `)

	parser := parser.NewParser("SELECT count(name id) , FROM oranges,  apples ")
	results, err := parser.Parse()
	if err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("\nParsed Stmt:\n%v\n", results)
	}
}
