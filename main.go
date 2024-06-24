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
	err = parser.Init()
	if err != nil {
		panic(err)
	}
	err = parser.Parse()
	if err != nil {
		panic(err)
	}
	table, exists := parser.Symbols.TableByName("tableA")
	if !exists {
		panic("table does not exists")
	}
	fmt.Println(table)
	fmt.Println("---")
	err = parser.Parse()
	if err != nil {
		panic(err)
	}
	table, exists = parser.Symbols.TableByName("tableA")
	if !exists {
		panic("table does not exists")
	}
	fmt.Println(table)
}
