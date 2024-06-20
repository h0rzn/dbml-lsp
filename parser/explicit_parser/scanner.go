package explicitparser

import (
	"bufio"
	"bytes"
	"io"
	"unicode"

	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type LexItem struct {
	value    string
	token    tokens.Token
	position tokens.Position
}

func (l *LexItem) IsToken(expected tokens.Token) bool {
	return (l.token & expected) != 0
}

// type Position struct {
// 	line   uint32
// 	offset uint32
// 	len    uint32
// }
//
// func (p *Position) String() string {
// 	return fmt.Sprintf("[%d:%d-%d]", p.line, p.offset, p.offset+p.len)
// }

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

	current := tokens.MapChar(char)
	if current == tokens.LINEBR {
		s.line += 1
		s.offset = 0
	} else if (current & tokens.REL_1TM) != 0 {
		// could potentially be <> and not just <
		if s.read() != '>' {
			s.unread()
		} else {
			s.offset += 1
			return LexItem{
				value: "<>",
				token: tokens.REL_MTN,
				position: tokens.Position{
					Line:   s.line,
					Offset: s.offset,
					Len:    2,
				},
			}
		}
	}

	item := LexItem{
		value: string(char),
		token: current,
		position: tokens.Position{
			Line:   s.line,
			Offset: s.offset,
			Len:    1,
		},
	}

	return item

}

func (s *Scanner) ScanComposite(endChar rune) LexItem {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var length uint32 = 1
	for {
		char := s.read()
		if char == tokens.EOFChar || char == endChar {
			break
		} else {
			length += 1
			buf.WriteRune(char)
		}
	}
	literal := buf.String()
	item := LexItem{
		value: literal,
		token: tokens.MapLiteral(literal),
		position: tokens.Position{
			Line:   s.line,
			Offset: s.offset - length,
			Len:    length,
		},
	}
	return item

}

// read reads the next rune (char) from the (buffered) reader.
// Returns the rune(0) if an error occurs (or eofChar is returned).
func (s *Scanner) read() rune {
	char, _, err := s.reader.ReadRune()
	if err != nil {
		return tokens.EOFChar
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
		if char == tokens.EOFChar {
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
		token: tokens.WHITESPACE,
		position: tokens.Position{
			Line:   s.line,
			Offset: s.offset - length,
			Len:    length,
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
		if char == tokens.EOFChar {
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
		token: tokens.MapLiteral(literal),
		position: tokens.Position{
			Line:   s.line,
			Offset: s.offset - length,
			Len:    length,
		},
	}
	return item
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
