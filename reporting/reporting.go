package reporting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"jtmelton.com/camas/domain"
)

func check(e error) {
	if e != nil {
		log.Fatalf(e.Error())
	}
}

// Reporter is a marker interface
type Reporter interface {
	write(findings []domain.Finding, _options domain.Options)
}

// SimpleTextReporter is an interface impl for text reporting
type SimpleTextReporter struct{}

func (r SimpleTextReporter) write(findings []domain.Finding, _options domain.Options) {

	var results string
	const NEWLINE = "\n"

	if findings == nil || len(findings) == 0 {
		results = "No Findings."
	} else {
		var sb strings.Builder

		sb.WriteString("Findings:" + NEWLINE)

		for _, finding := range findings {
			sb.WriteString("[" + finding.FilePath + ":" + strconv.Itoa(finding.LineNumber) + " (" + finding.RuleName + ") \"" + finding.Content + "\"]" + NEWLINE)
		}

		results = sb.String()
	}

	if _options.OutputFile == "" {
		// no output file specified ... write results to stdout
		fmt.Println(results)
	} else {
		err := ioutil.WriteFile(_options.OutputFile, []byte(results), 0644)
		check(err)
	}
}

// JSONReporter is an interface impl for text reporting
type JSONReporter struct{}

func (r JSONReporter) write(findings []domain.Finding, _options domain.Options) {

	findingsJSON, jsonErr := json.Marshal(findings)
	check(jsonErr)

	var results []byte

	// populate results with either empty json array or the actual findings
	if findings == nil || len(findings) == 0 {
		results = []byte("[]")
	} else {
		results = findingsJSON
	}

	if _options.OutputFile == "" {
		// no output file specified ... write results to stdout
		fmt.Println(string(results))
	} else {
		err := ioutil.WriteFile(_options.OutputFile, results, 0644)
		check(err)
	}
}

// WriteReport is the generic function called by an outside class to serialize the results.
// It picks internally the right format based on the options passed in
func WriteReport(findings []domain.Finding, _options domain.Options) {
	var reporter Reporter

	if _options.OutputFormat == "" || _options.OutputFormat == "txt" {
		// default is txt
		reporter = SimpleTextReporter{}
	} else if _options.OutputFormat == "json" {
		// use json
		reporter = JSONReporter{}
	} else {
		log.Fatal("Unexpected reporting format.")
	}

	reporter.write(findings, _options)
}
