package parser

import (
	"bufio"
	"bytes"
	"io"
	"log"
)

type Token int

const (
	// special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // fields, table_name

	// MISC Chars
	ASTERISK // *
	COMMA    // ,

	// keywords
	SELECT
	FROM
)

type Scanner struct {
	reader *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(r),
	}
}

// scan returns the next token and literal value
func (s *Scanner) scan() (token Token, literal string) {
	// read the ch rune
	ch := s.read()
	log.Println("Lexer: Scanned rune", ch)

	if isWhitespace(ch) {
		// if we see whitespace consumes all whitespace
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		// if we see letter we consume it as ident or reserved word
		log.Println("Lexer: detected letter", ch)
		s.unread()
		return s.scanIdent()
	}

	// otherwise read individual charectars
	switch ch {
	case eof:
		return EOF, ""
	case '*':
		return ASTERISK, string(ch)
	case ',':
		return COMMA, string(ch)
	}
	return ILLEGAL, string(ch)
}

// read reads the next rune from the buffered reader
// returns the rune(0) if and error occurs (or io.EOF is returned)
func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		if err != io.EOF {
			log.Printf("Warning: err reading rune, sending EOF %v\n", err)
		}
		return eof
	}
	log.Println("Lexer: read rune", ch)
	return ch
}

// unread places the previously read rune back on the reader
func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
}

// consumes all contiguous whitespaces
func (s *Scanner) scanWhitespace() (token Token, literal string) {
	var buf bytes.Buffer
	for read := s.read(); isWhitespace(read); read = s.read() {
		buf.WriteRune(read)
	}
	s.unread()
	log.Printf("Lexer: ScanWhitespace: Scanned %q", buf.String())
	return WS, buf.String()
}

// consumes all contiguous ident runes
func (s *Scanner) scanIdent() (token Token, literal string) {
	var buf bytes.Buffer
	for read := s.read(); isIdentRune(read); read = s.read() {
		buf.WriteRune(read)
	}
	s.unread()

	literal = buf.String()
	log.Println("Lexer: ScanIndent: Scanned", literal)
	switch literal {
	case "SELECT":
		return SELECT, literal
	case "FROM":
		return FROM, literal
	}
	return IDENT, literal
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9')
}

func isIdentRune(ch rune) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_'
}

var eof = rune(0)
