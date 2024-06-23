package main

import (
	"fmt"
	"os"

	"github.com/h0rzn/dbml-lsp/parser"
	explicitparser "github.com/h0rzn/dbml-lsp/parser/explicit_parser"
)

func main() {
	file, err := os.Open("test.dbml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	expliParser := explicitparser.NewParser(file)
	parser := parser.NewParser(expliParser)
	err = parser.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(parser.Symbols.Info())
	fmt.Println("---")
	table, err := parser.Symbols.FindTable(0)
	if err != nil {
		panic(err)
	}
	fmt.Println("found table:", table.Name, table.Position)
}
