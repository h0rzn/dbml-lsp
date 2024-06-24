package parser

import (
	"github.com/h0rzn/dbml-lsp/parser/strategy"
	"github.com/h0rzn/dbml-lsp/parser/symbols"
)

type Parser struct {
	strategy.Strategy
	Symbols *symbols.Storage
}

func NewParser(strategy strategy.Strategy) *Parser {
	return &Parser{
		strategy,
		symbols.NewStorage(),
	}
}

func (p *Parser) Init() error {
	p.SetSymbols(p.Symbols)
	return nil
}
