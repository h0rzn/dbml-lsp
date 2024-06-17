package parser

import (
	"fmt"
)

type RelationshipParser struct {
	*Parser
}

func (r *RelationshipParser) Parse(inline bool) (*Relationship, error) {
	relationship := &Relationship{}
	var item LexItem

	if !inline {
		item = r.scanWithoutWhitespace()
		// catch optional name
		if item.IsToken(IDENT) {
			relationship.Name = item.value
			item = r.scanWithoutWhitespace()
		}

		if item.IsToken(BRACE_OPEN) {
			item = r.scanWithoutWhitespace()
			if (item.token & LINEBR) == 0 {
				return nil, fmt.Errorf("found %q, expected linebr after '{' for long relationsip declaration %d", item.value, item.position.line)
			}
		} else {
			if item.IsToken(COLON) {
				fmt.Println("??")
			}
		}
		relationship, err := r.parseLong()
		return relationship, err
	}

	item, exists := r.expect(COLON)
	if !exists {
		return nil, fmt.Errorf("found %q, expected ':' (after 'ref')", item.value)
	}

	item = r.scanWithoutWhitespace()
	if !item.IsToken(G_RELATION_TYPE) {
		return nil, fmt.Errorf("found %q, expected relationship declaration", item.value)
	}

	sideRight, err := r.parseSide()
	if err != nil {
		return nil, err
	}
	if len(sideRight) > 2 {
		relationship.SchemeB = sideRight[0].value
		relationship.TableB = sideRight[1].value
		relationship.ColumnB = sideRight[2].value
	} else {
		relationship.TableB = sideRight[0].value
		relationship.ColumnB = sideRight[1].value
	}

	return relationship, nil
}

func (r *RelationshipParser) parseLong() (*Relationship, error) {
	relationship := &Relationship{}
	sideLeft, err := r.parseSide()
	if err != nil {
		return nil, err
	}
	if len(sideLeft) > 2 {
		relationship.SchemeA = sideLeft[0].value
		relationship.TableA = sideLeft[1].value
		relationship.ColumnA = sideLeft[2].value
	} else {
		relationship.TableA = sideLeft[0].value
		relationship.ColumnA = sideLeft[1].value
	}

	item := r.scanWithoutWhitespace()
	if !item.IsToken(G_RELATION_TYPE) {
		return nil, fmt.Errorf("found %q, expected relationship declaration", item.value)
	}
	relationship.Type = item.value

	sideRight, err := r.parseSide()
	if err != nil {
		return nil, err
	}

	if len(sideRight) > 2 {
		relationship.SchemeB = sideRight[0].value
		relationship.TableB = sideRight[1].value
		relationship.ColumnB = sideRight[2].value
	} else {
		relationship.TableB = sideRight[0].value
		relationship.ColumnB = sideRight[1].value
	}

	return relationship, nil
}

func (r *RelationshipParser) parseSide() ([]LexItem, error) {
	var relationSide []LexItem

	// minimum requirement is: tableA.columnA
	items, exists := r.expectSequence(IDENT, DOT, IDENT)
	if !exists {
		return relationSide, fmt.Errorf("found %v, expected table.column or scheme.table.column for relationship declaration", items)
	}
	relationSide = append(relationSide, items[0], items[2])

	// if one more ident attached,
	// then the first ident is scheme not table
	// -> schemeA.tableA.columnA
	item := r.scanWithoutWhitespace()
	if item.IsToken(DOT) {
		item = r.scanWithoutWhitespace()
		if item.IsToken(IDENT) {
			relationSide = append(relationSide, item)
		}
	} else {
		r.unscan()
	}
	return relationSide, nil
}
