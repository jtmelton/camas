package config

import (
	"bufio"
	"encoding/json"

	"io/ioutil"
	"log"
	"os"
	"regexp"
)

const (
	regexErrorMessage = `Invalid regex rule in configuration!
    Key: %s
    Rule: %s
    Error: %s
`
)

// Rule represents a rule from the config file
type Rule struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	Noise         int    `json:"noiseLevel"`
	AnalysisLayer string `json:"analysisLayer"`
	Source        string `json:"source"`
	Regex         *regexp.Regexp
}

// IgnoredPath represents an ignored path from the config file
type IgnoredPath struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// Config struct holds all config from the given JSON file.
type Config struct {
	Rules        []Rule        `json:"rules"`
	IgnoredPaths []IgnoredPath `json:"ignored-paths"`
}

// ParseConfig function takes the configuration file path as input and returns the Config struct as output
func ParseConfig(configFilePath string) Config {
	// Read contents of JSON file
	f, err := os.Open(configFilePath)
	if err != nil {
		log.Fatalf("Unable to open file %s: %s", configFilePath, err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Unable to read file %s: %s", configFilePath, err)
	}

	// Parse JSON file and compile regex rules
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf("Unable to unmarshal configuration %s: %s", configFilePath, err)
	}

	configUpdatedRules := config
	configUpdatedRules.Rules = nil // empty out rules so we can recreate
	for _, rule := range config.Rules {
		if rule.Type != "regex" {
			configUpdatedRules.Rules = append(configUpdatedRules.Rules, rule)
			continue
		}

		// pre-vet rules to make sure they all compile
		re, err := regexp.Compile(rule.Content)
		if err != nil {
			log.Fatalf(regexErrorMessage, rule.Name, rule.Content, err)
		}

		// cache the regex after validation so it doesn't have to be compiled again
		rule.Regex = re
		configUpdatedRules.Rules = append(configUpdatedRules.Rules, rule)
	}

	return configUpdatedRules
}
