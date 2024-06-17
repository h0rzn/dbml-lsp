package parser

import (
	"fmt"
	"sync"
)

type Storage struct {
	*sync.Mutex
	tables    map[string]*TableStatement
	relations []*Relationship
}

func NewStorage() *Storage {
	return &Storage{
		&sync.Mutex{},
		make(map[string]*TableStatement),
		make([]*Relationship, 0),
	}
}

func (s *Storage) TableByName(name string) (*TableStatement, bool) {
	table, exists := s.tables[name]
	return table, exists
}

func (s *Storage) Tables() map[string]*TableStatement {
	return s.tables
}

// override if exists
func (s *Storage) PutTable(table *TableStatement) {
	s.Lock()
	s.tables[table.Name] = table
	s.Unlock()
}

func (s *Storage) DropTable(table *TableStatement) {
	s.Lock()
	delete(s.tables, table.Name)
	s.Unlock()
}

func (s *Storage) DropTableByName(name string) {
	s.Lock()
	delete(s.tables, name)
	s.Unlock()
}

func (s *Storage) ColumnsByTableName(name string) []*ColumnStatement {
	s.Lock()
	columns := make([]*ColumnStatement, 0)
	if table, exists := s.TableByName(name); exists {
		columns = table.Columns
	}
	s.Unlock()
	return columns
}

func (s *Storage) PutRelation(relation *Relationship) {
	s.Lock()
	s.relations = append(s.relations, relation)
	s.Unlock()
}

func (s *Storage) Relations() []*Relationship {
	return s.relations
}

func (s *Storage) Info() string {
	tableLen := len(s.tables)
	relationLen := len(s.relations)
	return fmt.Sprintf("Symbol Storage: %d Tables, %d Relations", tableLen, relationLen)
}