package main

import (
	"fmt"
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
	err = parser.Parse()
	if err != nil {
		panic(err)
	}
	fmt.Println(parser.Symbols.Info())
	for i, rel := range parser.Symbols.Relations() {
		fmt.Println(i, rel)
	}
	fmt.Println("---")
	for name, tbl := range parser.Symbols.Tables() {
		fmt.Println(name, tbl)
	}
}
