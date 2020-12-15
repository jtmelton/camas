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
			log.Fatalf("Could not construct .prof file for profiling: %v", err)
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

	configuration := config.ParseConfig(_options.ConfigFile)

	findings := processing.Walk(*inputDirectory, _options, configuration)

	reporting.WriteReport(findings, _options)

	/*
		TODO:
		- do a CI setup
			https://github.com/jandelgado/golang-ci-template-github-actions/blob/master/.github/workflows/test.yml

		- add a test for "create user ... identified by $&*Q#*@#(*" in a sql file
	*/

}
