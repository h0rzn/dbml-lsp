package parser

import (
	"errors"
	"fmt"
)

type KeyConstraint struct {
	Key   string
	Value string
}

type ConstraintParser struct {
	*Parser
}

func (c *ConstraintParser) Parse() ([]string, []*Relationship, error) {
	var constraints []string
	var relations []*Relationship
	var lastToken int
	for {
		constraintItem := c.scanWithoutWhitespace()
		switch constraintItem.token {
		case SQUARE_CLOSE:
			if len(constraints) == 0 {
				return nil, nil, errors.New("empty constraints declaration")
			}
			return constraints, relations, nil
		case COMMA:
			// TODO: handle first token: ';'
			if lastToken == COMMA {
				return nil, nil, fmt.Errorf("found %q, expected constraint delimiter", constraintItem.value)
			}
		case CONS_PK:
			constraints = append(constraints, constraintItem.value)
		case CONS_PRIMARY:
			item, found := c.expect(CONS_KEY)
			if !found {
				return nil, nil, fmt.Errorf("found %q, expected 'key' after 'primary'", item.value)
			}
			constraints = append(constraints, "primary key")
		case CONS_INCREMENT:
			fallthrough
		case CONS_UNIQUE:
			constraints = append(constraints, constraintItem.value)

		case NOTE:
			item, err := c.parseKeyConstraint(NOTE)
			if err != nil {
				fmt.Println(err)
			}
			keyedConstraintValue := item.Key + ":" + item.Value
			constraints = append(constraints, keyedConstraintValue)
		case CONS_NOT:
			item, found := c.expect(CONS_NULL)
			if !found {
				return nil, nil, fmt.Errorf("found %q, expected 'null' (not null)", item.value)
			}
			constraints = append(constraints, "not null")

		case REF_LOW:
			rel, err := c.parseRelationship(true)
			if err != nil {
				return nil, nil, err
			}
			// TODO: dont take [] as empty constraint declaration
			constraints = append(constraints, rel.String())
			relations = append(relations, rel)
		case UNKOWN:

		default:
			// error unkown token
			if constraintItem.token == IDENT {
				return nil, nil, fmt.Errorf("found %q, expected contraint", constraintItem.value)
			}
			return nil, nil, fmt.Errorf("unhandled non-ident item %q (%d), last: %d", constraintItem.value, constraintItem.token, lastToken)
		}
	}
}

func (c *ConstraintParser) parseKeyConstraint(keyToken Token) (constraint KeyConstraint, err error) {
	if keyToken != NOTE {
		return constraint, fmt.Errorf("unexpected keyed constraint: %q", keyToken)
	}
	constraint.Key = "note"
	item, found := c.expect(COLON)
	if !found {
		return constraint, fmt.Errorf("found %q, expected ':' (key-value-delimiter missing)", item.value)
	}

	item, found = c.expect(QUOTATION)
	if !found {
		return constraint, fmt.Errorf("found %q, expected /\" (enclosing quotation start missing)", item.value)
	}

	item = c.scanner.ScanComposite('"')

	constraint.Value = item.value
	return
}
