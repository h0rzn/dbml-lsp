package explicitparser

import (
	"fmt"

	"github.com/h0rzn/dbml-lsp/parser/symbols"
	"github.com/h0rzn/dbml-lsp/parser/tokens"
)

type ProjectParser struct {
	*Parser
}

func (p *ProjectParser) Parse() (*symbols.Project, error) {
	project := &symbols.Project{
		Options: make(map[string]string),
	}

	position, _, name, err := p.ParseDefinitionHead(tokens.PROJECT)
	if err != nil {
		return nil, err
	}
	project.Position = position
	project.Name = name

	for {
		keyItem := p.scanWithoutWhitespace()
		switch keyItem.token {
		case tokens.BRACE_CLOSE:
			return project, nil
		case tokens.LINEBR:
			continue
		default:
			if !keyItem.IsToken(tokens.G_PROJECT_OPTS) {
				return nil, fmt.Errorf("found %q, expected project option key", keyItem.value)
			}

			colonItem, found := p.expect(tokens.COLON)
			if !found {
				return nil, fmt.Errorf("found %q, expected ':' (key-value-delimiter missing)", colonItem.value)
			}

			quoteStartItem, found := p.expect(tokens.QUOTATION)
			if !found {
				return nil, fmt.Errorf("found %q, expected /\" (enclosing quotation start missing)", quoteStartItem.value)
			}

			valueItem := p.scanner.ScanComposite('"')
			project.Options[keyItem.value] = valueItem.value

		}
	}
}
