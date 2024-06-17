package parser

import (
	"fmt"
	"io"
)

var debug = true

type Parser struct {
	scanner *Scanner
	Symbols *Storage
	buffer  struct {
		current LexItem
		size    int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: NewScanner(r),
		Symbols: NewStorage(),
	}
}

func (p *Parser) Parse() error {
	for {
		item := p.scanWithoutWhitespace()
		if item.IsToken(EOF) || item.IsToken(BRACE_CLOSE) {
			break
		}
		if item.IsToken(LINEBR) {
			continue
		}

		switch item.token {
		case TABLE:
			p.unscan()
			table, err := p.parseTableDefinition()
			if err != nil {
				return err
			}
			p.Symbols.PutTable(table)
		case REF_CAP:
			// explicit pass of declaration type,
			// introducing token is not expected
			rel, err := p.parseRelationship(false)
			if err != nil {
				return err
			}
			fmt.Printf("relationship: %+v\n", rel)
			p.Symbols.PutRelation(rel)
		default:
			return fmt.Errorf("unexpected: %q", item.value)
		}
	}
	return nil
}

// scan returns next token from scanner.
// if token has been unscanned then read that instead.
func (p *Parser) scan() LexItem {
	// return current buffer token if exists
	if p.buffer.size > 0 {
		p.buffer.size = 0
		return p.buffer.current
	}

	// read next token from scanner
	// token, literal = p.scanner.Scan()
	item := p.scanner.Scan()

	// save to buffer
	p.buffer.current = item

	return item
}

// unscan puts last read token back to buffer
func (p *Parser) unscan() {
	p.buffer.size = 1
}

// scanWithoutWhitespace scans next token ignoring whitespace
func (p *Parser) scanWithoutWhitespace() LexItem {
	item := p.scan()
	if item.IsToken(WHITESPACE) {
		item = p.scan()
	}
	return item
}

func (p *Parser) parseTableDefinition() (*TableStatement, error) {
	parser := &TableParser{p}
	return parser.Parse()
}

func (p *Parser) parseRelationship(inline bool) (*Relationship, error) {
	parser := &RelationshipParser{p}
	return parser.Parse(inline)
}

// parseColumnDefinition parses a column definition.
// e.g. id integer [pk, unique]
// returns a column statement and error
func (p *Parser) parseColumnDefinition() (*ColumnStatement, error) {
	parser := &ColumnParser{p}
	return parser.Parse()
}

// parseConstraints returns a list of constraints.
// this function expects the opening square bracket '[' to be already read
// a, b, c]
// ^ starting position
func (p *Parser) parseConstraints() ([]string, []*Relationship, error) {
	parser := &ConstraintParser{p}
	return parser.Parse()
}

func (p *Parser) handleComment() {
	_, found := p.expectSequence(SLASH, SLASH)
	if !found {
		p.unscan()
		return
	}
	p.jumpLineEnd()
}

func (p *Parser) jumpLineEnd() {
	for {
		item := p.scan()
		if item.IsToken(LINEBR) {
			return
		}
	}
}

func (p *Parser) expect(expected Token) (item LexItem, found bool) {
	item = p.scanWithoutWhitespace()
	if !item.IsToken(expected) {
		return item, false
	}
	return item, true
}

func (p *Parser) expectAlternative(expected ...Token) (item LexItem, found bool) {
	item = p.scanWithoutWhitespace()
	for _, foundToken := range expected {
		if item.IsToken(foundToken) {
			return item, true
		}
	}
	return item, false
}

func (p *Parser) expectSequence(expected ...Token) (item []LexItem, found bool) {
	var items []LexItem
	for _, token := range expected {
		item, found := p.expect(token)
		if !found {
			return items, false
		}
		items = append(items, item)
	}
	return items, true
}
