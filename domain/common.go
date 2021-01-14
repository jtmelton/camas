package domain

// Options struct represents the cli options passed to the camas tool
type Options struct {
	InputDirectory string
	ConfigFile     string
	OutputFile     string
	OutputFormat   string
	NoiseLevel     int
	NumWorkers     int
}

// Finding struct represents a finding from the secrets analysis
type Finding struct {
	RuleName   string `json:"rule-name"`
	FilePath   string `json:"absolute-file-path"`
	LineNumber int    `json:"line-number"`
	Content    string `json:"content"`
	Noise      int    `json:"noise-level"`
}

// Findings struct represents an array of findings
type Findings struct {
	Findings []Finding `json:"findings"`
}
