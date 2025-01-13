package services

import (
	"encoding/json"
	"fmt"
	"hefestus-api/internal/models"
	"log"
	"os"
	"regexp"
	"sync"
)

type DictionaryService struct {
	dictionaries map[string]*models.ErrorDictionary
	domains      map[string]models.DomainConfig
	mu           sync.RWMutex
}

func NewDictionaryService() (*DictionaryService, error) {
	// Load domains configuration
	domainsConfig, err := loadDomainsConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load domains config: %w", err)
	}

	dictionaries := make(map[string]*models.ErrorDictionary)

	// Load dictionaries based on domain configurations
	for domain, config := range domainsConfig.Domains {
		dict, err := loadDomainDictionary(config.DictionaryPath)
		if err != nil {
			log.Printf("Warning: couldn't load dictionary for domain %s: %v", domain, err)
			continue
		}
		dictionaries[domain] = dict
	}

	return &DictionaryService{
		dictionaries: dictionaries,
		domains:      domainsConfig.Domains,
	}, nil
}

func loadDomainsConfig() (*models.DomainsConfig, error) {
	data, err := os.ReadFile("config/domains.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read domains config: %w", err)
	}

	var config models.DomainsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse domains config: %w", err)
	}

	return &config, nil
}

func (s *DictionaryService) GetDomainConfig(domain string) (models.DomainConfig, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	config, exists := s.domains[domain]
	return config, exists
}

func loadDomainDictionary(path string) (*models.ErrorDictionary, error) {
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

func (s *DictionaryService) FindMatches(domain string, errorText string) []models.ErrorPattern {
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
