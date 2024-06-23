package symbols

import (
	"errors"
	"fmt"
	"sync"
)

// Storage stores:
// - Project Info
// - Tables
// - Relations (Ref)
type Storage struct {
	*sync.Mutex
	project   *Project
	tables    map[string]*Table
	relations []*Relationship
}

func NewStorage() *Storage {
	return &Storage{
		&sync.Mutex{},
		&Project{},
		make(map[string]*Table),
		make([]*Relationship, 0),
	}
}

// Project
func (s *Storage) SetProject(project *Project) {
	s.project = project
}

func (s *Storage) GetProject() *Project {
	return s.project
}

// Table
func (s *Storage) TableByName(name string) (*Table, bool) {
	table, exists := s.tables[name]
	return table, exists
}

func (s *Storage) Tables() map[string]*Table {
	return s.tables
}

func (s *Storage) PutTable(table *Table) {
	s.Lock()
	s.tables[table.Name] = table
	s.Unlock()
}

func (s *Storage) UpdateTable(tableName string, updatedTable *Table) error {
	_, exists := s.TableByName(tableName)
	if !exists {
		return fmt.Errorf("failed to find table %q", tableName)
	}
	s.PutTable(updatedTable)

	return nil
}

func (s *Storage) DropTableByName(name string) {
	s.Lock()
	delete(s.tables, name)
	s.Unlock()
}

func (s *Storage) FindTable(line uint32) (*Table, error) {
	s.Lock()
	for _, table := range s.tables {
		if table.Position.Line == line {
			return table, nil
		}
	}
	s.Unlock()
	return nil, errors.New("not found")
}

// Column
func (s *Storage) ColumnsByTableName(name string) []*Column {
	s.Lock()
	columns := make([]*Column, 0)
	if table, exists := s.TableByName(name); exists {
		columns = table.Columns
	}
	s.Unlock()
	return columns
}

// Relation
func (s *Storage) PutRelation(relation *Relationship) {
	s.Lock()
	s.relations = append(s.relations, relation)
	s.Unlock()
}

func (s *Storage) Relations() []*Relationship {
	return s.relations
}

// Misc
func (s *Storage) Clear() {
	s.Lock()
	clear(s.tables)
	s.relations = nil
	s.Unlock()
}

func (s *Storage) Info() string {
	return fmt.Sprintf("Symbol Storage: [project defined: %t], %d Tables, %d Relations", s.project != nil, len(s.tables), len(s.relations))
}
