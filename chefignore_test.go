package chef_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/go-chef/chef"
	"github.com/stretchr/testify/assert"
)

var (
	testChefignore  = "test/chefignore"
	expectedIgnores = []string{
		".DS_Store",
		"ehthumbs.db",
		"Icon?",
		"nohup.out",
		"Thumbs.db",
		".envrc",
		".#*",
		".project",
		".settings",
		"*_flymake",
		"*_flymake.*",
		"*.bak",
		"*.sw[a-z]",
		"*.tmproj",
		"*~",
		"\\#*",
		"REVISION",
		"TAGS*",
		"tmtags",
		".vscode",
		".editorconfig",
		"*.class",
		"*.com",
		"*.dll",
		"*.exe",
		"*.o",
		"*.pyc",
		"*.so",
		"*/rdoc/",
		"a.out",
		"mkmf.log",
		".circleci/*",
		".codeclimate.yml",
		".delivery/*",
		".foodcritic",
		".kitchen*",
		".mdlrc",
		".overcommit.yml",
		".rspec",
		".rubocop.yml",
		".travis.yml",
		".watchr",
		".yamllint",
		"azure-pipelines.yml",
		"Dangerfile",
		"examples/*",
		"features/*",
		"Guardfile",
		"kitchen.yml*",
		"mlc_config.json",
		"Procfile",
		"Rakefile",
		"spec/*",
		"test/*",
		".git",
		".gitattributes",
		".gitconfig",
		".github/*",
		".gitignore",
		".gitkeep",
		".gitmodules",
		".svn",
		"*/.bzr/*",
		"*/.git",
		"*/.hg/*",
		"*/.svn/*",
		"Berksfile",
		"Berksfile.lock",
		"cookbooks/*",
		"tmp",
		"vendor/*",
		"Gemfile",
		"Gemfile.lock",
		"Policyfile.rb",
		"Policyfile.lock.json",
		"CODE_OF_CONDUCT*",
		"CONTRIBUTING*",
		"documentation/*",
		"TESTING*",
		"UPGRADING*",
		".vagrant",
		"Vagrantfile",
	}
	unignoredPaths = []string{
		"metadata.rb",
		"metadata.json",
		"recipes/default.rb",
		"attributes/default.rb",
		"resources/custom.rb",
		"providers/custom.rb",
		"templates/default/temp.txt.erb",
	}
)

func TestChefignoreNew(t *testing.T) {
	c := chef.NewChefignore(testChefignore)

	assert.Equal(t, testChefignore, c.Path)
	assert.Equal(t, expectedIgnores, c.Ignores)
}

func TestChefignoreIgnore(t *testing.T) {
	c := chef.NewChefignore(testChefignore)

	// Ensure appropriate files will be ignored
	for _, ignore := range c.Ignores {
		ignoreStr := strings.ReplaceAll(ignore, "*", "test")

		// Replace contents of [*] matchers
		re := regexp.MustCompile(`\[.*\]`)
		ignoreStr = re.ReplaceAllString(ignoreStr, "a")

		// Handle \#
		ignoreStr = strings.ReplaceAll(ignoreStr, "\\#", "#")
		assert.Equal(t, true, c.Ignore(ignoreStr), "Did not ignore ", ignoreStr, " as expected")
	}

	for _, path := range unignoredPaths {
		assert.Equal(t, false, c.Ignore(path))
	}
}
