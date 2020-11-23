package processing

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"jtmelton/camas/config"
	"jtmelton/camas/domain"
)

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func doesRuleMatch(content string, rule config.Rule, path string, lineNum int) (domain.Finding, bool) {

	if rule.Type == "simple" {
		if strings.Contains(content, rule.Content) {
			return domain.Finding{
				RuleName:   rule.Name,
				FilePath:   path,
				LineNumber: lineNum,
				Content:    content,
				Noise:      rule.Noise}, true
		}
	} else if rule.Type == "regex" {

		if rule.Regex.MatchString(content) {
			return domain.Finding{
				RuleName:   rule.Name,
				FilePath:   path,
				LineNumber: lineNum,
				Content:    content,
				Noise:      rule.Noise}, true
		}
	}

	return domain.Finding{}, false
}

func addFindingIfMatched(content string, rule config.Rule, path string, lineNum int, findings chan domain.Finding) {
	finding, match := doesRuleMatch(content, rule, path, lineNum)
	if match {
		findings <- finding
	}
	// fmt.Println(rule)
}

func worker(paths []string, _options domain.Options, config config.Config, jobs <-chan int, results chan<- int, findings chan domain.Finding) {
	for j := range jobs {

		path := paths[j-1]
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		filename := filepath.Base(path)
		extension := filepath.Ext(path)
		for _, rule := range config.Rules {
			if rule.Noise < _options.NoiseLevel {
				continue
			}

			if "extension" == rule.AnalysisLayer {
				// if the rule contains a leading '.', strip that away
				if strings.HasPrefix(extension, ".") {
					extension = trimLeftChar(extension)
				}
				addFindingIfMatched(extension, rule, path, 1, findings)
			} else if "filename" == rule.AnalysisLayer {
				addFindingIfMatched(filename, rule, path, 1, findings)
			} else if "path" == rule.AnalysisLayer {
				addFindingIfMatched(path, rule, path, 1, findings)
			}
		}

		scanner := bufio.NewScanner(file)
		// read file into enumerated file
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()

			for _, rule := range config.Rules {
				if "contents" == rule.AnalysisLayer {
					if rule.Noise < _options.NoiseLevel {
						continue
					}

					addFindingIfMatched(line, rule, path, lineNum, findings)
				}
			}
		}

		file.Close()
		results <- j
	}

}

// Walk walks the directory tree, forking a new worker goroutine for the number of workers
func Walk(rootPath string, _options domain.Options, config config.Config) []domain.Finding {

	// if num workers is not set as a cli config, default to # of cpus
	var defaultNumWorkers = runtime.NumCPU()
	var numWorkers = defaultNumWorkers
	if _options.NumWorkers != 0 {
		numWorkers = _options.NumWorkers
	}

	var paths []string
	filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}

		if info.IsDir() {
			return nil
		}

		for _, ignoredPath := range config.IgnoredPaths {
			if ignoredPath.Type == "simple" {
				if strings.Contains(path, ignoredPath.Content) {
					return nil
				}
			} else if ignoredPath.Type == "regex" {
				var regexTest, _ = regexp.Compile(ignoredPath.Content)

				if regexTest.MatchString(path) {
					return nil
				}

			}
		}

		paths = append(paths, path)

		return nil
	})

	jobs := make(chan int, len(paths))
	results := make(chan int, len(paths))

	fchan := make(chan domain.Finding)
	var findings []domain.Finding

	// launch a goroutine to catch findings as they come in
	go func() {
		for {
			f := <-fchan
			findings = append(findings, f)
		}
	}()

	// launch NUM_WORKERS # of jobs
	for w := 1; w <= numWorkers; w++ {
		go worker(paths, _options, config, jobs, results, fchan)
	}

	// queue up all the jobs
	for j := 1; j <= len(paths); j++ {
		jobs <- j
	}

	// wait for all results to complete
	for a := 1; a <= len(paths); a++ {
		<-results
	}

	return findings
}
