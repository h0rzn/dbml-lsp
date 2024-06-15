package parser

import (
	"errors"
	"fmt"
)

type TableParser struct {
	*Parser
}

func (t *TableParser) Parse() (*TableStatement, error) {
	statement := &TableStatement{}
	tableItem, found := t.expect(TABLE)
	if !found {
		return nil, fmt.Errorf("found %q, expected 'Table'", tableItem.value)
	}
	statement.Position = tableItem.position

	// find name declaration
	nameItem, found := t.expect(IDENT)
	if !found {
		return nil, fmt.Errorf("found %q, expected table name declaration", nameItem.value)
	}
	statement.Name = nameItem.value

	// find opening brace and linebreak
	_, found = t.expectSequence(BRACE_OPEN, LINEBR)
	if !found {
		return nil, errors.New("found ?, expected delimiter '{' for table head end")
	}

	// column definitions
	for {
		columnItem := t.scanWithoutWhitespace()
		switch columnItem.token {
		case LINEBR:
			continue
		case BRACE_CLOSE:
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
