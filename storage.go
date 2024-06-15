package main

import (
	"sync"

	"github.com/h0rzn/dbml-lsp/parser"
)

type Storage struct {
	*sync.Mutex
	tables map[string]*parser.TableStatement
}

func NewStorage() *Storage {
	return &Storage{
		&sync.Mutex{},
		make(map[string]*parser.TableStatement),
	}

}

func (s *Storage) GetTableByName(name string) (*parser.TableStatement, bool) {
	table, exists := s.tables[name]
	return table, exists
}

// override if exists
func (s *Storage) PutTable(table *parser.TableStatement) {
	s.Lock()
	s.tables[table.Name] = table
	s.Unlock()
}

func (s *Storage) DropTable(table *parser.TableStatement) {
	s.Lock()
	delete(s.tables, table.Name)
	s.Unlock()
}

func (s *Storage) DropTableByName(name string) {
	s.Lock()
	delete(s.tables, name)
	s.Unlock()
}

func (s *Storage) ColumnsByTableName(name string) []*parser.ColumnStatement {
	s.Lock()
	columns := make([]*parser.ColumnStatement, 0)
	if table, exists := s.GetTableByName(name); exists {
		columns = table.Columns
	}
	s.Unlock()
	return columns
}
