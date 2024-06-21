package explicitparser

import (
	"errors"
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
	projectItem, found := p.expect(tokens.PROJECT)
	if !found {
		return nil, fmt.Errorf("found %q, expected 'Project'", projectItem.value)
	}
	project.Position = projectItem.position

	// find name declaration
	nameItem, found := p.expect(tokens.IDENT)
	if !found {
		return nil, fmt.Errorf("found %q, expected table name declaration", nameItem.value)
	}
	project.Name = nameItem.value

	// find opening brace and linebreak
	_, found = p.expectSequence(tokens.BRACE_OPEN, tokens.LINEBR)
	if !found {
		return nil, errors.New("found ?, expected delimiter '{' for project definition head end")
	}

	for {
		keyItem := p.scanWithoutWhitespace()
		if keyItem.IsToken(tokens.BRACE_CLOSE | tokens.LINEBR) {
			return project, nil
		} else if keyItem.IsToken(tokens.G_PROJECT_OPTS) {
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
		} else {
			return nil, fmt.Errorf("found %q, expected project option key", keyItem.value)
		}
	}
}
