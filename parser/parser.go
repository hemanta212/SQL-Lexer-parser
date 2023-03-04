package parser

import (
	"fmt"
	"io"
	"log"
)

type Parser struct {
	s      *Scanner
	buffer struct {
		token   Token
		literal string
		n       int
	}
}

type SelectStatement struct {
	Fields    []string
	TableName string
	Behaviour string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (*SelectStatement, error) {
	if token, literal := p.scanIgnoreWhitespace(); token != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", literal)
	}
	stmt := &SelectStatement{}

	// we wanna see fields now
	for {
		token, literal := p.scanIgnoreWhitespace()
		if token != IDENT && token != ASTERISK {
			return nil, fmt.Errorf("found %q, expected fields", literal)
		}
		stmt.Fields = append(stmt.Fields, literal)

		// if next token is not comma then break loop
		if token, _ := p.scanIgnoreWhitespace(); token != COMMA {
			p.unscan()
			break
		}
	}

	// we wanna see FROM now
	if token, literal := p.scanIgnoreWhitespace(); token != FROM {
		return nil, fmt.Errorf("Found %q, expected FROM", literal)
	}

	if token, literal := p.scanIgnoreWhitespace(); token == IDENT {
		stmt.TableName = literal
	} else {
		return nil, fmt.Errorf("Found %q, expected table name", literal)
	}

	return stmt, nil
}

func (p *Parser) scan() (Token, string) {
	// if we have sth in buffer, return it instead of reading new
	if p.buffer.n != 0 {
		p.buffer.n = 0
		return p.buffer.token, p.buffer.literal
	}
	log.Println("Parser: scanning..")
	// otherwise scan into the buffer
	token, literal := p.s.scan()
	// save to the buffer for the unscan functionality
	p.buffer.token, p.buffer.literal = token, literal
	return token, literal
}

func (p *Parser) unscan() { p.buffer.n = 1 }

func (p *Parser) scanIgnoreWhitespace() (Token, string) {
	token, literal := p.scan()
	log.Println("Parser: sIW token", token, literal)
	if token == WS {
		token, literal = p.scan()
	}
	log.Printf("Parser: ScanIgnoreWhitespace: scanned %d %q", token, literal)
	return token, literal
}
