package parser

import (
	"fmt"
	"log"
	"strings"
)

type Parser struct {
	l      *lexer
	buffer struct {
		token item
		n     int
	}
	lastItem item
}

func NewParser(input string) *Parser {
	return &Parser{l: lex("Test", input)}
}

type SelectStmt struct {
	command     string
	columnNames []string
	tableNames  []string
	functions   []SqlFunction
}

type SqlFunction struct {
	name      string
	arguments []string
}

type CreateStmt struct {
	command   string
	tableName string
	function  SqlCreateFunction
}

type SqlCreateFunction struct {
	name      string
	arguments []DatatypeArgument
}

type DatatypeArgument struct {
	name     string
	datatype []string
}

func (p *Parser) Parse() (any, error) {
	parsedItems := []item{}

	fmt.Println("Scanning SQL:\n", p.l.input)

	// we wanna see fields now
	for {
		token := p.l.nextItem()
		if token.typ == itemError {
			return nil, fmt.Errorf("%s at line %d %d", token.val, token.line, token.pos)
		} else if token.typ == itemWs {
			continue
		}
		token.val = strings.ToLower(token.val)
		parsedItems = append(parsedItems, token)
		if token.typ == itemEOF {
			break
		}
	}
	printTokens(parsedItems)

	command := parsedItems[0].val
	switch command {
	case "select":
		return p.prepareSelectStatement(parsedItems[1:])
	case "create":
		return p.prepareCreateStatement(parsedItems[1:])
	default:
		return nil, fmt.Errorf("Error while parsing, Unknown sql command %s", command)
	}
}

func (p *Parser) prepareSelectStatement(items []item) (*SelectStmt, error) {
	columnItems, index, err := p.scanColumnItems(items)
	log.Println("Parser: prepareSelectStmt: columnItems", columnItems)
	if err != nil {
		return nil, err
	}

	functions, err := p.scanFunctions(columnItems)
	log.Println("Parser: prepareSelectStmt: functions", functions)
	if err != nil {
		return nil, err
	}

	columnNames := filterColumnNames(columnItems)
	log.Println("Parser: prepareSelectStmt: columnNames", columnNames)

	tableNames, index, err := p.scanTableNames(items[index:])
	if err != nil {
		return nil, err
	}

	return &SelectStmt{
		command:     "select",
		columnNames: columnNames,
		tableNames:  tableNames,
		functions:   functions,
	}, nil
}

func (p *Parser) prepareCreateStatement(items []item) (*CreateStmt, error) {
	return &CreateStmt{
		command:   "create",
		tableName: "",
		function:  SqlCreateFunction{},
	}, nil
}

// scans all the keywords as column names till 'from' keyword
// returns the scanned list, current index position of items or error
// ignores the *valid* function and args for further processing by other funcs
// Test cases:
// select col1, col2 from names  | colnames: [col1, col2]
// select col1 col2 from names  | colnames: [col1]
// select col1 col2, col3 col4 from names  | colnames: [col1, col3]
// select col1, from names  | colnames: [], error: yes
func (p *Parser) scanColumnItems(items []item) ([]item, int, error) {
	results := []item{}
	items = p.labelFunctions(items)
	count := 0
	// from above eg: only columns followed by 'select' or comma are valid
	canAcceptColumn := true
	// every arg must have func but this is here to know if that func is valid (after comma)
	// so if we are ignoring certain func because of no comma then we ignore its args too
	isInsideFunction := false
	for i, ident := range items {
		count, p.lastItem = i, ident
		if ident.typ == itemFunction && canAcceptColumn {
			isInsideFunction = true
			results = append(results, ident)
			canAcceptColumn = false
		} else if (ident.typ == itemFunctionArg || ident.typ == itemAsterisk) && isInsideFunction {
			results = append(results, ident)
		} else if ident.val == "from" {
			count += 1
			break
		} else if ident.typ == itemComma && !isInsideFunction {
			// this comma is from column separator
			canAcceptColumn = true
		} else if ident.typ == itemComma && isInsideFunction {
			// this comma is from function args separator
			results = append(results, ident)
		} else if ident.typ == itemFunctionClose {
			results = append(results, ident)
			isInsideFunction = false
		} else if ident.typ == itemEOF {
			return nil, 0, fmt.Errorf(
				"Parsing error: Unexpected EOF near %q at line %d pos %d", p.lastItem.val, ident.line, ident.pos)
		} else if (ident.typ == itemIdent || ident.typ == itemAsterisk) && canAcceptColumn {
			results = append(results, ident)
			canAcceptColumn = false
		}
	}
	if len(results) == 0 {
		return nil, 0, fmt.Errorf(
			"Parsing error: Insufficient column names near %q at line %d pos %d", p.lastItem.val, p.lastItem.line, p.lastItem.pos)
	} else if canAcceptColumn {
		// if canAccpetColumn is still true, then extra trailing comma is there
		return nil, 0, fmt.Errorf(
			"Parsing error: Unexpected trailing comma near %q at line %d pos %d", p.lastItem.val, p.lastItem.line, p.lastItem.pos)
	}
	return results, count, nil
}

// similar to scanColumnames with little sensible differences
func (p *Parser) scanTableNames(items []item) ([]string, int, error) {
	results := []string{}
	count := 0
	canAcceptTable := true
	for i, ident := range items {
		count, p.lastItem = i, ident
		if ident.val == "where" {
			count += 1
			break
		} else if ident.typ == itemComma {
			canAcceptTable = true
		} else if ident.typ == itemEOF && canAcceptTable {
			return nil, 0, fmt.Errorf(
				"Parsing error: Unexpected eof near %q at line %d pos %d", p.lastItem.val, ident.line, ident.pos)
		} else if ident.typ == itemEOF {
			break
		} else if canAcceptTable {
			results = append(results, ident.val)
			canAcceptTable = false
		}
	}
	if len(results) == 0 {
		return nil, 0, fmt.Errorf(
			"Parsing error: Insufficient table names near %q at line %d pos %d", p.lastItem.val, p.lastItem.line, p.lastItem.pos)
	} else if canAcceptTable {
		return nil, 0, fmt.Errorf(
			"Parsing error: Unexpected trailing comma near %q at line %d pos %d", p.lastItem.val, p.lastItem.line, p.lastItem.pos)
	}

	return results, count, nil
}

func (p *Parser) scanFunctions(items []item) ([]SqlFunction, error) {
	//filter only the function and args
	results := []SqlFunction{}
	currFunc := SqlFunction{}
	isinsideFunc := false
	// for comma control, if Func(a b) given; then raise error
	canAcceptArg := false
	for _, it := range items {
		p.lastItem = it
		if it.typ == itemFunction {
			if isinsideFunc {
				results = append(results, currFunc)
			}
			currFunc = SqlFunction{name: it.val}
			isinsideFunc, canAcceptArg = true, true
		} else if (it.typ == itemFunctionArg || it.typ == itemAsterisk) && canAcceptArg {
			currFunc.arguments = append(currFunc.arguments, it.val)
			canAcceptArg = false
		} else if it.typ == itemComma && isinsideFunc {
			canAcceptArg = true
		} else if it.typ == itemFunctionClose {
			results = append(results, currFunc)
			isinsideFunc = false
		} else if canAcceptArg && !isinsideFunc {
			return nil, fmt.Errorf(
				"Parsing error: Unexpected token %q near %q expected comma at line %d pos %d", it.val, p.lastItem.val, it.line, it.pos)
		}
	}
	return results, nil
}

func (p *Parser) labelFunctions(items []item) []item {
	functionIndexes := []int{}
	// find all those itemIdent with () (i.e itemFuncs)
	for index, it := range items {
		if it.typ == itemFunctionOpen {
			functionIndexes = append(functionIndexes, index-1)
		}
	}
	log.Println("Parser: labelFunctions: funcIndexes", functionIndexes)
	// label the item as itemFunction instead of itemIdent
	for _, index := range functionIndexes {
		it := items[index]
		it.typ = itemFunction
		items[index] = it
	}
	log.Println("Parser: labelFunctions: final items", items)
	return items
}

// removes function and args type from column items
// returns the string representation of column items i.e item.val
func filterColumnNames(items []item) []string {
	results := []string{}
	for _, it := range items {
		if it.typ == itemIdent || it.typ == itemAsterisk {
			results = append(results, it.val)
		}
	}
	return results
}

func printTokens(tokens []item) {
	for _, token := range tokens {
		typeS := map[itemType]string{
			0: "itemError",
			1: "itemEOF",
			2: "itemWs",
			3: "itemIdent",
			4: "itemCommand",
			5: "itemFunctionArg",
			6: "itemAsterisk",
			7: "itemComma",
		}
		if token.typ == itemWs {
			continue
		}
		log.Printf("{type: %s,\tvalue: %q}\n", typeS[token.typ], token.val)
	}
}
