package models

import (
	"time"
)

type Config struct {
	Parallelism  int
	Delay        time.Duration
	VercelBypass string
}

type Report struct {
	HTTPErrors    []Issue     `json:"httpErrors"`
	BrokenLinks   []Issue     `json:"brokenLinks"`
	Redirects     []Redirects `json:"redirects"`
	HTMLStructure []Issue     `json:"htmlStructureIssues"`
	SEO           []Issue     `json:"seoIssues"`
	Accessibility []Issue     `json:"accessibilityIssues"`
	Debug         []Issue     `json:"debug"`
}

type Issue struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Src     string `json:"src,omitempty"`
}

type Redirects struct {
	StatusCode int    `json:"statusCode"`
	Referer    string `json:"referer"`
	Path       string `json:"path"`
}
