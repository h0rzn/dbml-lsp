package symbols

import (
	"encoding/json"
	"fmt"

	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type Table struct {
	Name     string
	Columns  []*Column
	Position tokens.Position
}

func (t *Table) String() string {
	tableJson, err := json.Marshal(t)
	if err != nil {
		fmt.Println("table String() err", err.Error())
		return "{<error>}"
	}

	return string(tableJson)
}

func (t *Table) Print() {
	fmt.Println("===")
	fmt.Printf("Table '%s' @ %s\n", t.Name, t.Position.String())
	for _, column := range t.Columns {
		fmt.Printf("[%s] %s %s # %v (%d)\n", column.Name, column.Type, column.Position.String(), column.Constraints, len(column.Constraints))
	}

}

type Column struct {
	Name        string
	Type        string
	Constraints []*Constraint
	Position    tokens.Position
}

func (c *Column) String() string {
	var out string
	out += fmt.Sprintf("%s [%s] ", c.Name, c.Type)
	for _, constraint := range c.Constraints {
		if len(constraint.Key) > 0 {
			out += fmt.Sprintf("[%q: %q] ", constraint.Key, constraint.Value)
		} else {
			out += constraint.Value
		}
	}

	return out
}

type Constraint struct {
	// Only for 'key: value' constraints
	// empty string otherwise
	Key string
	// Actual value like "pk"
	Value string
}

type Relationship struct {
	Name    string
	SchemeA string
	TableA  string
	ColumnA string
	SchemeB string
	TableB  string
	ColumnB string
	Type    string
}

func (r *Relationship) String() string {
	format := "%s %s.%s.%s %s %s.%s.%s"
	return fmt.Sprintf(format, r.Name, r.SchemeA, r.TableA, r.ColumnA, r.Type, r.SchemeB, r.TableB, r.ColumnB)

}
