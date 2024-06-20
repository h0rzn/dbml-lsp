package explicitparser

/*
import (
	"encoding/json"
	"fmt"
)

type TableStatement struct {
	Name     string
	Columns  []*ColumnStatement
	Position Position
}

func (t *TableStatement) String() string {
	// var out bytes.Buffer
	tableJson, err := json.Marshal(t)

	if err != nil {
		fmt.Println("table String() err", err.Error())
		return "{<error>}"
	}

	// err = json.Indent(&out, []byte(tableJson), "", "  ")
	// if err != nil {
	// 	fmt.Println("table String() err", err.Error())
	// 	return "{<error>}"
	// }
	// return out.String()
	return string(tableJson)
}

func (t *TableStatement) Print() {
	fmt.Println("===")
	fmt.Printf("Table '%s' @ %s\n", t.Name, t.Position.String())
	for _, column := range t.Columns {
		fmt.Printf("[%s] %s %s # %v (%d)\n", column.Name, column.Type, column.Position.String(), column.Constraints, len(column.Constraints))
	}

}

type EnumStatement struct {
	Name string
}

type ColumnStatement struct {
	Name        string
	Type        string
	Constraints []string
	Position    Position
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
*/