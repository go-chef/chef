package chef

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Matches all lines that only contain comments or regexp
const commentsAndWhitespaceRegex = `^\s*(?:#.*)?$`

// chefignore object
type Chefignore struct {
	Path    string
	Ignores []string
}

// Creates a new Chefignore from a path
func NewChefignore(path string) Chefignore {
	chefignore := Chefignore{
		Path:    path,
		Ignores: make([]string, 0),
	}

	// Parse ignored lines, respecting comments
	chefignore.Ignores = parseChefignoreContent(path)
	return chefignore
}

// Returns true if the given path should be ignored, false otherwise
func (chefignore *Chefignore) Ignore(path string) bool {
	for _, ignorePattern := range chefignore.Ignores {
		// TODO: Handle antd style wildcards (**) since these are supported in the spec
		//       but not by filepath.Match
		matched, err := filepath.Match(ignorePattern, path)
		if err != nil {
			// Skip malformed patterns
			continue
		}

		if matched {
			return true
		}
	}
	return false
}

// Parses .chefignore content
func parseChefignoreContent(path string) []string {
	skipLineRegex := regexp.MustCompile(commentsAndWhitespaceRegex)

	// Read content line-by-line
	file, err := os.Open(path)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	ignores := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !skipLineRegex.MatchString(line) {
			ignores = append(ignores, strings.TrimSpace(line))
		}
	}

	return ignores
}
