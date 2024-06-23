package symbols

import (
	"fmt"
	"sync"
)

type Storage struct {
	*sync.Mutex
	project *Project
	tables  map[uint32]*Table
}

func NewStorage() *Storage {
	return &Storage{
		&sync.Mutex{},
		&Project{},
		make(map[uint32]*Table),
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
	// table, exists := s.tables[name]
	for _, table := range s.tables {
		if table.Name == name {
			return table, true
		}
	}

	return nil, false
}

func (s *Storage) Tables() map[uint32]*Table {
	return s.tables
}

func (s *Storage) PutTable(table *Table) {
	s.Lock()
	s.tables[table.Position.Line] = table
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
	var line uint32
	var found bool

	for l, table := range s.tables {
		if table.Name == name {
			line = l
			found = true
		}
	}
	if found {
		delete(s.tables, line)
	}
	s.Unlock()
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

// Misc
func (s *Storage) Clear() {
	s.Lock()
	clear(s.tables)
	s.Unlock()
}

func (s *Storage) Info() string {
	return fmt.Sprintf("Symbol Storage: [project defined: %t], %d Tables", s.project != nil, len(s.tables))
}
