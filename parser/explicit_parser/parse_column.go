package explicitparser

import (
	"fmt"

	"github.com/h0rzn/dbml-lsp/parser/symbols"
	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type ColumnParser struct {
	*Parser
}

func (c *ColumnParser) Parse() (*symbols.Column, error) {
	statement := &symbols.Column{}

	// colum name
	nameItem, found := c.expect(tokens.IDENT)
	if !found {
		return nil, fmt.Errorf("found %q, expected column name", nameItem.value)
	}
	statement.Name = nameItem.value
	statement.Position = nameItem.position

	// column type
	typeItem, found := c.expect(tokens.IDENT)
	if !found {
		return nil, fmt.Errorf("found %q, expected column type", typeItem.value)
	}
	statement.Type = typeItem.value

	// look for constraints
	item, found := c.expect(tokens.SQUARE_OPEN)
	if !found {
		if item.token != tokens.LINEBR {
			return nil, fmt.Errorf("found %q, expected column definition stop", item.value)
		}
	} else {
		// constraints definition found
		constraints, relations, err := c.parseConstraints()
		if err != nil {
			return nil, fmt.Errorf("incorrect constraint declaration: %s", err.Error())
		}
		statement.Constraints = constraints
		for _, relation := range relations {
			relation.ColumnA = nameItem.value
			c.Symbols.PutRelation(relation)
		}
	}
	return statement, nil
}
