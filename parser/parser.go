package parser

import (
	"github.com/h0rzn/dbml-lsp/parser/strategy"
	"github.com/h0rzn/dbml-lsp/parser/symbols"
)

type Parser struct {
	parser  strategy.Strategy
	Symbols *symbols.Storage
}

func NewParser(strategy strategy.Strategy) *Parser {
	return &Parser{
		parser:  strategy,
		Symbols: symbols.NewStorage(),
	}
}

func (p *Parser) Parse() error {
	p.parser.SetSymbols(p.Symbols)
	return p.parser.Parse()
}
