package explicitparser

import (
	"github.com/h0rzn/dbml-lsp/parser/symbols"
	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type TableParser struct {
	*Parser
}

func (t *TableParser) Parse() (*symbols.Table, error) {
	statement := &symbols.Table{}
	position, scheme, name, err := t.ParseDefinitionHead(tokens.TABLE)
	if err != nil {
		return nil, err
	}
	statement.Position = position
	statement.Scheme = scheme
	statement.Name = name

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
