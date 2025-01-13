package models

type ErrorPattern struct {
	Pattern    string   `json:"pattern"`
	Category   string   `json:"category"`
	Solutions  []string `json:"solutions"`
	References []string `json:"references"`
}

type ErrorDictionary struct {
	Patterns map[string]ErrorPattern `json:"patterns"`
}
