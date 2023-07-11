package internal

import (
	"errors"
	"log"
	"strings"
)

type noteService struct {
	repo NodeRepository
}

func NewNoteService(repo NodeRepository) *noteService {
	return &noteService{
		repo: repo,
	}
}

func (s *noteService) New(title, content string) error {
	note, err := NewNote(title, content)
	if err != nil {
		return err
	}
	return s.repo.Save(note)
}

func (s *noteService) ListAll() ([]Node, error) {
	notes, err := s.repo.GetNodes()
	if err != nil {
		return nil, errors.New("Error to consult notes")
	}
	return notes, nil
}

func (s *noteService) GetByTitle(title string) (Node, error) {
	notes, err := s.ListAll()
	if err != nil {
		return Node{}, err
	}
	found := -1
	for i := 0; i < len(notes); i++ {
		if notes[i].Title == title {
			found = i
			break
		}
	}
	if found == -1 {
		return Node{}, errors.New("Note note found")
	}

	return notes[found], nil
}

func filtersToMap(filters []Filter) map[string]string {
	fields := map[string]bool{"title": true, "tags": true}
	mapFilters := map[string]string{}
	for _, v := range filters {
		if v.Field != "" && fields[strings.ToLower(v.Field)] {
			if v.Value != "" {
				mapFilters[strings.ToLower(v.Field)] = v.Value
			}
		}
	}
	return mapFilters
}

func (s *noteService) Find(_filters []Filter) ([]Node, error) {
	notes, err := s.ListAll()
	if err != nil {
		return nil, err
	}
	filters := filtersToMap(_filters)

	var founds []Node
	for _, note := range notes {
		if IncludeNote(filters, note) {
			founds = append(founds, note)
		}
	}
	return founds, nil
}

func IncludeNote(filters map[string]string, note Node) bool {
	val, ok := filters["tags"]
	if ok {
		meta, err := ExtractMetadata(note.Content)
		if err != nil {
			log.Printf("Error extracting metadata from Note %v", note.Title)
			return false
		}
		if !meta.IncludeTags(val) {
			log.Printf("Tags %v no included in %v", val, note.Title)
			return false
		}
	}
	val, ok = filters["title"]
	if ok {
		source := strings.ToLower(note.Title)
		target := strings.ToLower(val)
		if !strings.Contains(source, target) {
			return false
		}
	}
	return true
}
