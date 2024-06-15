package parser

import (
	"errors"
	"fmt"
	"io"
)

var debug = true

type Parser struct {
	scanner *Scanner
	buffer  struct {
		// literal string
		current LexItem
		// token   Token
		size int
	}
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: NewScanner(r),
	}
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
	if item.token == WHITESPACE {
		item = p.scan()
	}
	return item
}

// func (p *Parser) scanEnquoted() LexItem {
// 	item
//
// }

func (p *Parser) Parse() ([]*TableStatement, error) {
	var tables []*TableStatement
	for {
		item := p.scanWithoutWhitespace()
		if item.token == EOF || item.token == BRACE_CLOSE {
			break
		}
		if item.token == LINEBR {
			continue
		}
		p.unscan()

		table, err := p.parseTableDefinition()
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (p *Parser) ParseDEBUG() {
	for {
		item := p.scanWithoutWhitespace()
		if item.token == EOF {
			break
		}
		if item.token == LINEBR {
			fmt.Println("linebr")
		}
	}
}

func (p *Parser) parseTableDefinition() (*TableStatement, error) {
	statement := &TableStatement{}
	tableItem, found := p.expect(TABLE)
	if !found {
		return nil, fmt.Errorf("found %q, expected 'Table'", tableItem.value)
	}
	statement.Position = tableItem.position

	// find name declaration
	nameItem, found := p.expect(IDENT)
	if !found {
		return nil, fmt.Errorf("found %q, expected table name declaration", nameItem.value)
	}
	statement.Name = nameItem.value

	// find opening brace and linebreak
	_, found = p.expectSequence(BRACE_OPEN, LINEBR)
	if !found {
		return nil, errors.New("found ?, expected delimiter '{' for table head end")
	}

	// column definitions
	for {
		columnItem := p.scanWithoutWhitespace()
		switch columnItem.token {
		case LINEBR:
			continue
		case BRACE_CLOSE:
			return statement, nil
		default:
			p.unscan()
			column, err := p.parseColumnDefinition()
			if err != nil {
				return nil, err
			}
			statement.Columns = append(statement.Columns, column)
		}
	}
}

// parseColumnDefinition parses a column definition.
// e.g. id integer [pk, unique]
// returns a column statement and error
func (p *Parser) parseColumnDefinition() (*ColumnStatement, error) {
	statement := &ColumnStatement{}

	// colum name
	nameItem, found := p.expect(IDENT)
	if !found {
		return nil, fmt.Errorf("found %q, expected column name", nameItem.value)
	}
	statement.Name = nameItem.value
	statement.Position = nameItem.position

	// column type
	typeItem, found := p.expect(IDENT)
	if !found {
		return nil, fmt.Errorf("found %q, expected column type", typeItem.value)
	}
	statement.Type = typeItem.value

	// look for constraints
	item, found := p.expect(SQUARE_OPEN)
	if !found {
		if item.token != LINEBR {
			return nil, fmt.Errorf("found %q, expected column definition stop", item.value)
		}
	} else {
		// constraints definition found
		constraints, err := p.parseConstraints()
		if err != nil {
			return nil, fmt.Errorf("incorrect constraint declaration: %s", err.Error())
		}
		statement.Constraints = constraints
	}

	return statement, nil
}

// parseConstraints returns a list of constraints.
// this function expects the opening square bracket '[' to be already read
// a, b, c]
// ^ starting position
func (p *Parser) parseConstraints() ([]string, error) {
	var constraints []string
	var lastToken int
	for {
		constraintItem := p.scanWithoutWhitespace()
		switch constraintItem.token {
		case SQUARE_CLOSE:
			if len(constraints) == 0 {
				return nil, errors.New("empty constraints declaration")
			}
			return constraints, nil
		case COMMA:
			// TODO: handle first token: ';'
			if lastToken == COMMA {
				return nil, fmt.Errorf("found %q, expected constraint delimiter", constraintItem.value)
			}
		case CONS_PK:
			constraints = append(constraints, constraintItem.value)
		case CONS_PRIMARY:
			item, found := p.expect(CONS_KEY)
			if !found {
				return nil, fmt.Errorf("found %q, expected 'key' after 'primary'", item.value)
			}
			constraints = append(constraints, "primary key")
		case CONS_INCREMENT:
			fallthrough
		case CONS_UNIQUE:
			constraints = append(constraints, constraintItem.value)

		case NOTE:
			item, err := p.parseKeyConstraint(NOTE)
			if err != nil {
				fmt.Println(err)
			}
			keyedConstraintValue := item.Key + ":" + item.Value
			constraints = append(constraints, keyedConstraintValue)
		case CONS_NOT:
			item, found := p.expect(CONS_NULL)
			if !found {
				return nil, fmt.Errorf("found %q, expected 'null' (not null)", item.value)
			}
			constraints = append(constraints, "not null")

		case UNKOWN:

		default:
			// error unkown token
			if constraintItem.token == IDENT {
				return nil, fmt.Errorf("found %q, expected contraint", constraintItem.value)
			}
			return nil, fmt.Errorf("unhandled non-ident item %q (%d), last: %d", constraintItem.value, constraintItem.token, lastToken)
		}
	}
}

type KeyConstraint struct {
	Key   string
	Value string
}

func (p *Parser) parseKeyConstraint(keyToken Token) (constraint KeyConstraint, err error) {
	if keyToken != NOTE {
		return constraint, fmt.Errorf("unexpected keyed constraint: %q", keyToken)
	}
	constraint.Key = "note"
	item, found := p.expect(COLON)
	if !found {
		return constraint, fmt.Errorf("found %q, expected ':' (key-value-delimiter missing)", item.value)
	}

	item, found = p.expect(QUOTATION)
	if !found {
		return constraint, fmt.Errorf("found %q, expected /\" (enclosing quotation start missing)", item.value)
	}

	item = p.scanner.ScanComposite('"')

	constraint.Value = item.value
	return
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
		if item.token == LINEBR {
			return
		}
	}
}

type scanWhileFunc func(LexItem) bool

func (p *Parser) scanWhile(whileFunc scanWhileFunc) string {
	var out string
	for {
		item := p.scan()
		if item.token == LINEBR || !whileFunc(item) {
			fmt.Printf("returning %q\n", out)
			return out
		}
		out += item.value
	}
	return "???"
}

func (p *Parser) expect(expected Token) (item LexItem, found bool) {
	item = p.scanWithoutWhitespace()
	if item.token != expected {
		return item, false
	}
	return item, true
}

func (p *Parser) expectAlternative(expected ...Token) (item LexItem, found bool) {
	item = p.scanWithoutWhitespace()
	for _, foundToken := range expected {
		if item.token == foundToken {
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
