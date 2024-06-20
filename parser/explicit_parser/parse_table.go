package explicitparser

import (
	"errors"
	"fmt"

	"github.com/h0rzn/dbml-lsp/parser/symbols"
	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type TableParser struct {
	*Parser
}

func (t *TableParser) Parse() (*symbols.Table, error) {
	statement := &symbols.Table{}
	tableItem, found := t.expect(tokens.TABLE)
	if !found {
		return nil, fmt.Errorf("found %q, expected 'Table'", tableItem.value)
	}
	statement.Position = tableItem.position

	// find name declaration
	nameItem, found := t.expect(tokens.IDENT)
	if !found {
		return nil, fmt.Errorf("found %q, expected table name declaration", nameItem.value)
	}
	statement.Name = nameItem.value

	// find opening brace and linebreak
	_, found = t.expectSequence(tokens.BRACE_OPEN, tokens.LINEBR)
	if !found {
		return nil, errors.New("found ?, expected delimiter '{' for table head end")
	}

	// column definitions
	for {
		columnItem := t.scanWithoutWhitespace()
		switch columnItem.token {
		case tokens.LINEBR:
			continue
		case tokens.BRACE_CLOSE:
			return statement, nil
		default:
			t.unscan()
			column, err := t.parseColumnDefinition()
			if err != nil {
				return nil, err
			}
			statement.Columns = append(statement.Columns, column)
		}
	}
}
