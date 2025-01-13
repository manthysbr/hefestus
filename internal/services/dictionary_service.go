package services

import (
	"encoding/json"
	"hefestus-api/internal/models"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type DictionaryService struct {
	dictionaries map[string]*models.ErrorDictionary
	mu           sync.RWMutex
}

func NewDictionaryService() (*DictionaryService, error) {
	domains := []string{"kubernetes", "github", "argocd"}
	dictionaries := make(map[string]*models.ErrorDictionary)

	for _, domain := range domains {
		dict, err := loadDomainDictionary(domain)
		if err != nil {
			// Log warning but continue loading other dictionaries
			log.Printf("Warning: couldn't load dictionary for domain %s: %v", domain, err)
			continue
		}
		dictionaries[domain] = dict
	}

	return &DictionaryService{
		dictionaries: dictionaries,
	}, nil
}

func (s *DictionaryService) FindMatches(domain, errorText string) []models.ErrorPattern {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dict, ok := s.dictionaries[domain]
	if !ok {
		return nil
	}

	var matches []models.ErrorPattern
	for _, pattern := range dict.Patterns {
		if matched, _ := regexp.MatchString(pattern.Pattern, errorText); matched {
			matches = append(matches, pattern)
		}
	}
	return matches
}

func loadDomainDictionary(domain string) (*models.ErrorDictionary, error) {
	path := filepath.Join("data", "patterns", domain+"_errors.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var dict models.ErrorDictionary
	if err := json.Unmarshal(data, &dict); err != nil {
		return nil, err
	}

	return &dict, nil
}
