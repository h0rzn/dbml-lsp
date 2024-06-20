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
	parser.Parse()
	// fmt.Println(parser.Symbols.Info())
	fmt.Println("---")
	fmt.Println(parser.Symbols.ColumnsByTableName("tableA"))
}
