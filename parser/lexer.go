package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

var eofChar = rune(0)

type LexItem struct {
	value    string
	token    Token
	position Position
}

type Position struct {
	line   uint32
	offset uint32
	len    uint32
}

func (p *Position) String() string {
	return fmt.Sprintf("[%d:%d-%d]", p.line, p.offset, p.offset+p.len)
}

type Scanner struct {
	reader *bufio.Reader
	// current focused line
	line uint32
	// current offset in focused line
	offset uint32
}

func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(reader),
		line:   0,
		offset: 0,
	}
}

// Scan fetches next token and literal value
func (s *Scanner) Scan() LexItem {
	char := s.read()
	if isWhitespace(char) {
		s.unread()
		return s.scanWhitespace()
	}
	if isLetter(char) {
		s.unread()
		return s.scanIdent()
	}

	token := mapSpecialChar(char)
	if token == LINEBR {
		s.line += 1
		s.offset = 0
	}
	item := LexItem{
		value: string(char),
		token: token,
		position: Position{
			line:   s.line,
			offset: s.offset,
			len:    1,
		},
	}

	// return mapSpecialChar(char), string(char)
	// return token, string(char)
	return item

}

func (s *Scanner) ScanComposite(endChar rune) LexItem {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var length uint32 = 1
	for {
		char := s.read()
		if char == eofChar || char == endChar {
			break
		} else {
			length += 1
			buf.WriteRune(char)
		}
	}
	literal := buf.String()
	item := LexItem{
		value: literal,
		token: mapIdentLiteral(literal),
		position: Position{
			line:   s.line,
			offset: s.offset - length,
			len:    length,
		},
	}
	return item

}

// Scan unfiltered (with whitespace, ...)
func (s *Scanner) Scan2(includeWhitespace bool) LexItem {
	char := s.read()
	fmt.Printf("scan2 %q\n", char)
	if includeWhitespace {
		if isWhitespace(char) {
			s.unread()
			return s.scanWhitespace()
		}
	}
	if isLetter(char) {
		fmt.Println("scan scan ident")
		s.unread()
		return s.scanIdent()
	}

	token := mapSpecialChar(char)
	if token == LINEBR {
		s.line += 1
		s.offset = 0
	}
	item := LexItem{
		value: string(char),
		token: token,
		position: Position{
			line:   s.line,
			offset: s.offset,
			len:    1,
		},
	}

	// return mapSpecialChar(char), string(char)
	// return token, string(char)
	return item
}

// read reads the next rune (char) from the (buffered) reader.
// Returns the rune(0) if an error occurs (or eofChar is returned).
func (s *Scanner) read() rune {
	char, _, err := s.reader.ReadRune()
	if err != nil {
		return eofChar
	}
	s.offset += 1
	return char
}

// unread puts the last read rune back on the reader
func (s *Scanner) unread() error {
	s.offset -= 1
	return s.reader.UnreadRune()
}

//
// scanning helpers
//

// scanWhitespace consumes current rune and contigous whitespace
func (s *Scanner) scanWhitespace() LexItem {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var length uint32 = 1
	for {
		char := s.read()
		if char == eofChar {
			break
		} else if !isWhitespace(char) {
			s.unread()
			break
		} else {
			length += 1
			buf.WriteRune(char)
		}
	}
	item := LexItem{
		value: "",
		token: WHITESPACE,
		position: Position{
			line:   s.line,
			offset: s.offset - length,
			len:    length,
		},
	}
	return item
	// return WHITESPACE, buf.String()
}

// scanWhitespace consumes current rune and contigous whitespace
func (s *Scanner) scanIdent() LexItem {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var length uint32 = 1
	for {
		char := s.read()
		if char == eofChar {
			break
		} else if !isLetter(char) && !isDigit(char) {
			s.unread()
			break
		} else {
			length += 1
			buf.WriteRune(char)

		}
	}
	literal := buf.String()
	item := LexItem{
		value: literal,
		token: mapIdentLiteral(literal),
		position: Position{
			line:   s.line,
			offset: s.offset - length,
			len:    length,
		},
	}
	return item
}

func mapIdentLiteral(literal string) Token {
	switch literal {
	case "Table":
		return TABLE
	case "enum":
		return ENUM
	case "pk":
		return CONS_PK
	case "primary":
		return CONS_PRIMARY
	case "key":
		return CONS_KEY
	case "null":
		return CONS_NULL
	case "not":
		return CONS_NOT
	case "increment":
		return CONS_INCREMENT
	case "unique":
		return CONS_UNIQUE
	case "note":
		return NOTE
	default:

	}
	return IDENT
}

func mapSpecialChar(char rune) Token {
	switch char {
	case '\n':
		return LINEBR
	case eofChar:
		return EOF
	case '/':
		return SLASH
	case '{':
		return BRACE_OPEN
	case '}':
		return BRACE_CLOSE
	case '[':
		return SQUARE_OPEN
	case ']':
		return SQUARE_CLOSE
	case '"':
		return QUOTATION
	case ',':
		return COMMA
	case ':':
		return COLON
	// case '`':
	// return BACKTICK
	// ...handle other chars
	default:
	}
	return UNKOWN
}

//
// character classes
//

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' // || ch == '\n'
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
	// return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
