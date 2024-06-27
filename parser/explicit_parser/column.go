package explicitparser

import (
	"fmt"

	"github.com/h0rzn/dbml-lsp/parser/symbols"
	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type ColumnParser struct {
	*Parser
}

func (c *ColumnParser) Parse() (*symbols.Column, []*symbols.Relationship, error) {
	statement := &symbols.Column{}
	relations := make([]*symbols.Relationship, 0)

	// colum name
	nameItem, found := c.expect(tokens.IDENT)
	if !found {
		return nil, relations, fmt.Errorf("found %q, expected column name", nameItem.value)
	}
	statement.Name = nameItem.value
	statement.Position = nameItem.position

	// column type
	typeItem, found := c.expect(tokens.IDENT)
	if !found {
		return nil, relations, fmt.Errorf("found %q, expected column type", typeItem.value)
	}
	statement.Type = typeItem.value

	// look for constraints
	item, found := c.expect(tokens.SQUARE_OPEN)
	if !found {
		if item.token != tokens.LINEBR {
			return nil, relations, fmt.Errorf("found %q, expected column definition stop", item.value)
		}
	} else {
		// constraints definition found
		constraints, rels, err := c.parseConstraints()
		if err != nil {
			return nil, relations, fmt.Errorf("incorrect constraint declaration: %s", err.Error())
		}
		statement.Constraints = constraints
		table := c.GetTableCtx()
		if table != nil {
			for _, relation := range rels {
				relation.SchemeA = table.Scheme
				relation.TableA = table.Name
				relation.ColumnA = nameItem.value
			}
		}
		relations = rels
	}
	return statement, relations, nil
}
