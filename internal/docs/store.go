package docs

import (
	"sort"
	"strings"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	entries []DocEntry
	icons   []IconEntry
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Reload(entries []DocEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = entries
}

func (s *Store) ReloadIcons(icons []IconEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.icons = icons
}

func (s *Store) Search(query string) []DocEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(query)
	var results []DocEntry

	for _, entry := range s.entries {
		if strings.Contains(strings.ToLower(entry.Name), q) ||
			strings.Contains(strings.ToLower(entry.Description), q) ||
			strings.Contains(strings.ToLower(entry.Content), q) {
			results = append(results, entry)
		}
	}

	return results
}

func (s *Store) GetByName(name string) *DocEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := strings.ToLower(name)

	for _, entry := range s.entries {
		if strings.ToLower(entry.Name) == n {
			return &entry
		}
	}

	return nil
}

func (s *Store) List(category string) []DocEntrySummary {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []DocEntrySummary
	cat := strings.ToLower(category)

	for _, entry := range s.entries {
		if cat != "" && strings.ToLower(string(entry.Category)) != cat {
			continue
		}

		results = append(results, DocEntrySummary{
			Name:        entry.Name,
			Category:    entry.Category,
			Description: entry.Description,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Category != results[j].Category {
			return results[i].Category < results[j].Category
		}
		return results[i].Name < results[j].Name
	})

	return results
}

func (s *Store) SearchIcons(query string) []IconEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	q := strings.ToLower(query)
	var results []IconEntry

	for _, icon := range s.icons {
		if strings.Contains(strings.ToLower(icon.Name), q) ||
			strings.Contains(strings.ToLower(icon.Description), q) ||
			strings.Contains(strings.ToLower(icon.Category), q) {
			results = append(results, icon)
		}
	}

	return results
}
