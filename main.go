package main

import (
	"os"

	"github.com/h0rzn/dbml-lsp/parser"
)

func main() {
	file, err := os.Open("test.dbml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	parser := parser.NewParser(file)
	// parser.ParseDEBUG()
	tables, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	for _, table := range tables {
		// fmt.Println("---")
		// fmt.Println(table.Name, table.Position.String())
		// for _, column := range table.Columns {
		// 	fmt.Println("-- column", column.Name, column.Position.String())
		// 	for _, constraint := range column.Constraints {
		// 		fmt.Printf("\t %q", constraint)
		// 	}
		// 	fmt.Println()
		// }
		table.Print()

	}
	_ = tables
}
