package explicitparser

import (
	"errors"
	"fmt"
	"io"

	"github.com/h0rzn/dbml-lsp/parser/symbols"
	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type Parser struct {
	scanner  *Scanner
	Symbols  *symbols.Storage
	tableCtx *symbols.Table
	buffer   struct {
		current LexItem
		size    int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: NewScanner(r),
	}
}

func (p *Parser) SetSymbols(storage *symbols.Storage) {
	p.Symbols = storage
}

func (p *Parser) SetTableCtx(table *symbols.Table) {
	p.tableCtx = table
}

func (p *Parser) GetTableCtx() *symbols.Table {
	return p.tableCtx
}

func (p *Parser) Parse() error {
	for {
		item := p.scanWithoutWhitespace()
		if item.IsToken(tokens.EOF | tokens.BRACE_CLOSE) {
			break
		}
		if item.IsToken(tokens.LINEBR) {
			continue
		}

		switch item.token {
		case tokens.PROJECT:
			p.unscan()
			project, err := p.parseProjectDefinition()
			if err != nil {
				return err
			}
			p.Symbols.SetProject(project)

		case tokens.TABLE:
			p.unscan()
			table, err := p.parseTableDefinition()
			if err != nil {
				return err
			}
			p.Symbols.PutTable(table)

		case tokens.REF_CAP:
			// explicit pass of declaration type,
			// introducing token is not expected
			rel, err := p.parseRelationship(false)
			if err != nil {
				return err
			}
			table, exists := p.Symbols.TableByName(rel.TableA)
			if !exists {
				return fmt.Errorf("reference: host table %q does not exist", rel.TableA)
			}
			table.References = append(table.References, rel)
			err = p.Symbols.UpdateTable(table.Name, table)
			if err != nil {
				return err
			}

		case tokens.SLASH:
			item, exists := p.expect(tokens.SLASH)
			if !exists {
				return fmt.Errorf("found %q, expected trailing '/' (comment)", item.value)
			}

			for {
				item := p.scan()
				if item.IsToken(tokens.LINEBR) {
					break
				}
			}

		default:
			return fmt.Errorf("unexpected: %q", item.value)
		}
	}

	return nil
}

func (p *Parser) ParseDefinitionHead(startToken tokens.Token) (position tokens.Position, scheme string, name string, err error) {
	startItem, found := p.expect(startToken)
	if !found {
		return position, scheme, name, fmt.Errorf("found %q, expected definition type", startItem.value)
	}

	nameItem, found := p.expect(tokens.IDENT)
	if !found {
		return position, scheme, name, fmt.Errorf("found %q, expected definition name declaration", nameItem.value)
	}

	nextItem := p.scan()
	if nextItem.IsToken(tokens.DOT) {
		name2Item, found := p.expect(tokens.IDENT)
		if !found {
			return position, scheme, name, fmt.Errorf("found %q, expected name after '.'", name2Item.value)
		}
		name = name2Item.value
		scheme = nameItem.value
	} else if nextItem.IsToken(tokens.WHITESPACE) {
		name = nameItem.value
	} else {
		// unhandled token
		return position, scheme, name, fmt.Errorf("unexpected %q", nextItem.value)
	}

	_, found = p.expectSequence(tokens.BRACE_OPEN, tokens.LINEBR)
	if !found {
		return position, scheme, name, errors.New("found ?, expected delimiter '{' for definition head end")
	}
	return startItem.position, scheme, name, nil
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
	if item.IsToken(tokens.WHITESPACE) {
		item = p.scan()
	}
	return item
}

func (p *Parser) parseProjectDefinition() (*symbols.Project, error) {
	parser := &ProjectParser{p}
	return parser.Parse()
}

func (p *Parser) parseTableDefinition() (*symbols.Table, error) {
	parser := &TableParser{p}
	return parser.Parse()
}

func (p *Parser) parseRelationship(inline bool) (*symbols.Relationship, error) {
	parser := &RelationshipParser{p}
	return parser.Parse(inline)
}

// parseColumnDefinition parses a column definition.
// e.g. id integer [pk, unique]
// returns a column statement and error
func (p *Parser) parseColumnDefinition() (*symbols.Column, []*symbols.Relationship, error) {
	parser := &ColumnParser{p}
	return parser.Parse()
}

// parseConstraints returns a list of constraints.
// this function expects the opening square bracket '[' to be already read
// a, b, c]
// ^ starting position
func (p *Parser) parseConstraints() ([]*symbols.Constraint, []*symbols.Relationship, error) {
	parser := &ConstraintParser{p}
	return parser.Parse()
}

func (p *Parser) expect(expected tokens.Token) (item LexItem, found bool) {
	item = p.scanWithoutWhitespace()
	if !item.IsToken(expected) {
		return item, false
	}
	return item, true
}

func (p *Parser) expectAlternative(expected ...tokens.Token) (item LexItem, found bool) {
	item = p.scanWithoutWhitespace()
	for _, foundToken := range expected {
		if item.IsToken(foundToken) {
			return item, true
		}
	}
	return item, false
}

func (p *Parser) expectSequence(expected ...tokens.Token) (item []LexItem, found bool) {
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
