package parser

import (
	"fmt"
)

type Parser struct {
	l      *lexer
	buffer struct {
		token   itemType
		literal string
		n       int
	}
}

func NewParser(input string) *Parser {
	return &Parser{l: lex("Test", input)}
}

func (p *Parser) Parse() ([]item, error) {
	parsedItems := []item{}

	// we wanna see fields now
	for {
		token := p.l.nextItem()
		if token.typ == itemEOF {
			break
		} else if token.typ == itemError {
			return nil, fmt.Errorf("%s at line %d %d", token.val, token.line, token.pos)
		}
		parsedItems = append(parsedItems, token)
	}
	return parsedItems, nil
}

// func (p *Parser) scan() (itemType, string) {
// 	// if we have sth in buffer, return it instead of reading new
// 	if p.buffer.n != 0 {
// 		p.buffer.n = 0
// 		return p.buffer.token, p.buffer.literal
// 	}
// 	log.Println("Parser: scanning..")
// 	// otherwise scan into the buffer
// 	token, literal := p.s.scan()
// 	// save to the buffer for the unscan functionality
// 	p.buffer.token, p.buffer.literal = token, literal
// 	return token, literal
// }

// func (p *Parser) unscan() { p.buffer.n = 1 }

// func (p *Parser) scanIgnoreWhitespace() (itemType, string) {
// 	token, literal := p.scan()
// 	log.Println("Parser: sIW token", token, literal)
// 	if token == itemWS {
// 		token, literal = p.scan()
// 	}
// 	log.Printf("Parser: ScanIgnoreWhitespace: scanned %d %q", token, literal)
// 	return token, literal
// }
