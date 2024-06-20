package strategy

import "github.com/h0rzn/dbml-lsp/parser/symbols"

// Strategy Inferface
// for different parser systems
type Strategy interface {
	Parse() error
	SetSymbols(*symbols.Storage)
}
