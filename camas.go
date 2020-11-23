// camas is a tool that finds potential secrets in source code

package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"jtmelton/camas/config"
	"jtmelton/camas/domain"
	"jtmelton/camas/processing"
	"jtmelton/camas/reporting"
)

var (
	inputDirectory *string
	configFile     *string
	outputFile     *string
	outputFormat   *string
	noiseLevel     *int
	numWorkers     *int
)

// // Options struct represents the cli options passed to the camas tool
// type Options struct {
// 	inputDirectory string
// 	configFile     string
// 	outputFile     string
// 	outputFormat   string
// 	noiseLevel     int
// 	numWorkers     int
// }

// // Finding struct represents a finding from the secrets analysis
// type Finding struct {
// 	RuleName   string `json:"rule-name"`
// 	FilePath   string `json:"absolute-file-path"`
// 	LineNumber int    `json:"line-number"`
// 	Content    string `json:"content"`
// 	Noise      int    `json:"noise-level"`
// }

func main() {
	inputDirectory = flag.String("inputDirectory", "", "Directory to analyze (Required)")
	configFile = flag.String("configFile", "", "Configuration File (Required)")
	outputFile = flag.String("outputFile", "", "Output File")
	outputFormat = flag.String("outputFormat", "", "Output Format [txt, json]")
	noiseLevel = flag.Int("noiseLevel", 0, "minimum noise level to report on")
	numWorkers = flag.Int("numWorkers", 0, "number of go workers to execute")
	var cpuProfile = flag.Bool("cpuProfile", false, "write cpu profile to file")

	flag.Parse()

	if *inputDirectory == "" || *configFile == "" {
		flag.PrintDefaults()

		os.Exit(1)
	}

	if *cpuProfile == true {
		f, err := os.Create("camas.prof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	_options := domain.Options{
		InputDirectory: *inputDirectory,
		ConfigFile:     *configFile,
		OutputFile:     *outputFile,
		OutputFormat:   *outputFormat,
		NoiseLevel:     *noiseLevel,
		NumWorkers:     *numWorkers,
	}

	config := config.ParseConfig(_options.ConfigFile)

	findings := processing.Walk(*inputDirectory, _options, config)

	reporting.WriteReport(findings, _options)

	/*
		TODO:
		- DONE docckerize (with from scratch)
		- DONE pull in rules from all the different projects
				https://github.com/awslabs/git-secrets/blob/master/git-secrets#L233
				* BSD3 https://gitlab.com/gitlab-com/gl-security/security-operations/gl-redteam/token-hunter/-/blob/master/regexes.json
				* Apache 2.0 https://github.com/newrelic/rusty-hog/blob/master/src/default_rules.json
				? MIT https://github.com/eth0izzle/shhgit/blob/master/config.yaml
				X GPL3 https://github.com/BishopFox/GitGot/blob/master/checks/default.list
				* MIT https://github.com/UKHomeOffice/repo-security-scanner/blob/master/rules/gitrob.json
				? MIT https://github.com/michenriksen/gitrob/blob/master/core/signatures.go
				? MIT https://github.com/zricethezav/gitleaks/blob/master/examples/leaky-repo.toml
					https://github.com/pelletier/go-toml
					$ curl https://raw.githubusercontent.com/zricethezav/gitleaks/master/examples/leaky-repo.toml > leaky-repo.toml
					docker run -v $PWD:/workdir pelletier/go-toml tomljson /workdir/example.toml

				https://github.com/nielsing/yar/blob/master/config/yarconfig.json
				https://github.com/dxa4481/truffleHogRegexes/blob/master/truffleHogRegexes/regexes.json
		- DONE make NUM_WORKERS a parameter, maybe defaulted to # of processors
		- DONE add noise level parameter (default might be "all" rules)
		- DONE start using regex_matcher code
		- DONE add excluded file types ... like https://github.com/securing/DumpsterDiver/blob/master/config.yaml#L4
		- DONE add ability to look at file name, file ext, file path, and contents of file
				"analysis-layer": "contents",
				"analysis-layer": "extension",
				"analysis-layer": "filename",
				"analysis-layer": "path",
		- DONE unit test everything
		- DONE extract reporting into its' own file
		- DONE use output file for reporting if parameter is present
		- DONE default reporting is "summary", options are "summary", "full", and "json"


		- use logging properly (use for everything except reporting)


		- do a CI setup
			https://github.com/jandelgado/golang-ci-template-github-actions/blob/master/.github/workflows/test.yml

		- make sure to have test for "create user ... identified by $&*Q#*@#(*" in a sql file
	*/

}
