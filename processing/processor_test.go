package processing

import (
	"log"
	"regexp"
	"testing"

	"jtmelton.com/camas/config"
	"jtmelton.com/camas/domain"
)

func getRegex(regex string) *regexp.Regexp {
	re, err := regexp.Compile(regex)
	if err != nil {
		log.Fatalf("Could not compile regex %s", regex)
	}
	return re
}

func TestDoesRuleMatch(t *testing.T) {
	tables := []struct {
		content       string
		rule          config.Rule
		path          string
		lineNum       int
		resultMatch   bool
		resultFinding domain.Finding
	}{
		{
			"$access_key = 'abc',",
			config.Rule{Name: "AWS Access Key ID Value", Type: "regex", Content: "(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}", Noise: 1, AnalysisLayer: "contents", Regex: getRegex("(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}")},
			"/path/1",
			1,
			false,
			domain.Finding{},
		},
		{
			"$access_key = 'AKIAIGL7R2LDUXGQK4IA',",
			config.Rule{Name: "AWS Access Key ID Value", Type: "regex", Content: "(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}", Noise: 1, AnalysisLayer: "contents", Regex: getRegex("(A3T[A-Z0-9]|AKIA|AGPA|AIDA|AROA|AIPA|ANPA|ANVA|ASIA)[A-Z0-9]{16}")},
			"/path/2",
			1,
			true,
			domain.Finding{RuleName: "AWS Access Key ID Value", FilePath: "/path/2", LineNumber: 1, Content: "$access_key = 'AKIAIGL7R2LDUXGQK4IA',", Noise: 1},
		},
		{
			"my.file",
			config.Rule{Name: "Private SSH key", Type: "regex", Content: "^.*_rsa$", Noise: 1, AnalysisLayer: "filename", Regex: getRegex("^.*_rsa$")},
			"/path/3",
			1,
			false,
			domain.Finding{},
		},
		{
			".id_rsa",
			config.Rule{Name: "Private SSH key", Type: "regex", Content: "^.*_rsa$", Noise: 1, AnalysisLayer: "filename", Regex: getRegex("^.*_rsa$")},
			"/path/4",
			1,
			true,
			domain.Finding{RuleName: "Private SSH key", FilePath: "/path/4", LineNumber: 1, Content: ".id_rsa", Noise: 1},
		},
		{
			"/my/path/to/my.file",
			config.Rule{Name: "SSH configuration file", Type: "regex", Content: "\\.?ssh/config$", Noise: 1, AnalysisLayer: "path", Regex: getRegex("\\.?ssh/config$")},
			"/path/5",
			1,
			false,
			domain.Finding{},
		},
		{
			".ssh/config",
			config.Rule{Name: "SSH configuration file", Type: "regex", Content: "\\.?ssh/config$", Noise: 1, AnalysisLayer: "path", Regex: getRegex("\\.?ssh/config$")},
			"/path/6",
			1,
			true,
			domain.Finding{RuleName: "SSH configuration file", FilePath: "/path/6", LineNumber: 1, Content: ".ssh/config", Noise: 1},
		},
		{
			"/my/path/to/my.file",
			config.Rule{Name: "Potential cryptographic key bundle", Type: "simple", Content: "p12", Noise: 1, AnalysisLayer: "extension"},
			"/path/7",
			1,
			false,
			domain.Finding{},
		},
		{
			"/my/path/to/somesecretfile.p12",
			config.Rule{Name: "Potential cryptographic key bundle", Type: "simple", Content: "p12", Noise: 1, AnalysisLayer: "extension"},
			"/path/8",
			1,
			true,
			domain.Finding{RuleName: "Potential cryptographic key bundle", FilePath: "/path/8", LineNumber: 1, Content: "/my/path/to/somesecretfile.p12", Noise: 1},
		},
	}

	for _, table := range tables {
		finding, match := doesRuleMatch(table.content, table.rule, table.path, table.lineNum)
		if table.resultMatch != match {
			// do something
			t.Errorf("Finding for '%s' with rule '%v' that was expected as '%t' (received '%t') has failed.", table.content, table.rule, table.resultMatch, match)
		} else if table.resultFinding != finding {
			// log.Printf("Finding received: %v", finding)
			// log.Printf("Finding expected: %v", table.resultFinding)
			t.Errorf("Finding for '%s' with rule '%v' that was expected with finding '%v' (received '%v') has failed.", table.content, table.rule, table.resultFinding, finding)

		}

	}
}
