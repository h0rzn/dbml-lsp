package explicitparser

import (
	"errors"
	"fmt"

	"github.com/h0rzn/dbml-lsp/parser/symbols"
	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type KeyConstraint struct {
	Key   string
	Value string
}

type ConstraintParser struct {
	*Parser
}

func (c *ConstraintParser) Parse() ([]string, []*symbols.Relationship, error) {
	var constraints []string
	var relations []*symbols.Relationship
	var lastToken int
	for {
		constraintItem := c.scanWithoutWhitespace()
		switch constraintItem.token {
		case tokens.SQUARE_CLOSE:
			if len(constraints) == 0 {
				return nil, nil, errors.New("empty constraints declaration")
			}
			return constraints, relations, nil
		case tokens.COMMA:
			// TODO: handle first token: ';'
			if lastToken == tokens.COMMA {
				return nil, nil, fmt.Errorf("found %q, expected constraint delimiter", constraintItem.value)
			}
		case tokens.CONS_PK:
			constraints = append(constraints, constraintItem.value)
		case tokens.CONS_PRIMARY:
			item, found := c.expect(tokens.CONS_KEY)
			if !found {
				return nil, nil, fmt.Errorf("found %q, expected 'key' after 'primary'", item.value)
			}
			constraints = append(constraints, "primary key")
		case tokens.CONS_INCREMENT:
			fallthrough
		case tokens.CONS_UNIQUE:
			constraints = append(constraints, constraintItem.value)

		case tokens.NOTE:
			item, err := c.parseKeyConstraint(tokens.NOTE)
			if err != nil {
				fmt.Println(err)
			}
			keyedConstraintValue := item.Key + ":" + item.Value
			constraints = append(constraints, keyedConstraintValue)
		case tokens.CONS_NOT:
			item, found := c.expect(tokens.CONS_NULL)
			if !found {
				return nil, nil, fmt.Errorf("found %q, expected 'null' (not null)", item.value)
			}
			constraints = append(constraints, "not null")

		case tokens.REF_LOW:
			rel, err := c.parseRelationship(true)
			if err != nil {
				return nil, nil, err
			}
			// TODO: dont take [] as empty constraint declaration
			constraints = append(constraints, rel.String())
			relations = append(relations, rel)
		case tokens.UNKOWN:

		default:
			// error unkown token
			if constraintItem.token == tokens.IDENT {
				return nil, nil, fmt.Errorf("found %q, expected contraint", constraintItem.value)
			}
			return nil, nil, fmt.Errorf("unhandled non-ident item %q (%d), last: %d", constraintItem.value, constraintItem.token, lastToken)
		}
	}
}

func (c *ConstraintParser) parseKeyConstraint(keyToken tokens.Token) (constraint KeyConstraint, err error) {
	if keyToken != tokens.NOTE {
		return constraint, fmt.Errorf("unexpected keyed constraint: %q", keyToken)
	}
	constraint.Key = "note"
	item, found := c.expect(tokens.COLON)
	if !found {
		return constraint, fmt.Errorf("found %q, expected ':' (key-value-delimiter missing)", item.value)
	}

	item, found = c.expect(tokens.QUOTATION)
	if !found {
		return constraint, fmt.Errorf("found %q, expected /\" (enclosing quotation start missing)", item.value)
	}

	item = c.scanner.ScanComposite('"')

	constraint.Value = item.value
	return
}
