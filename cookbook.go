package chef

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const metaRbName = "metadata.rb"
const metaJsonName = "metadata.json"

type metaFunc func(s []string, m *CookbookMeta) error

var metaRegistry map[string]metaFunc

// CookbookService  is the service for interacting with chef server cookbooks endpoint
type CookbookService struct {
	client *Client
}

// CookbookItem represents a object of cookbook file data
type CookbookItem struct {
	Url         string `json:"url,omitempty"`
	Path        string `json:"path,omitempty"`
	Name        string `json:"name,omitempty"`
	Checksum    string `json:"checksum,omitempty"`
	Specificity string `json:"specificity,omitempty"`
}

// CookbookListResult is the summary info returned by chef-api when listing
// http://docs.opscode.com/api_chef_server.html#cookbooks
type CookbookListResult map[string]CookbookVersions

// CookbookRecipesResult is the summary info returned by chef-api when listing
// http://docs.opscode.com/api_chef_server.html#cookbooks-recipes
type CookbookRecipesResult []string

// CookbookVersions is the data container returned from the chef server when listing all cookbooks
type CookbookVersions struct {
	Url      string            `json:"url,omitempty"`
	Versions []CookbookVersion `json:"versions,omitempty"`
}

// CookbookVersion is the data for a specific cookbook version
type CookbookVersion struct {
	Url     string `json:"url,omitempty"`
	Version string `json:"version,omitempty"`
}

// CookbookMeta represents a Golang version of cookbook metadata
type CookbookMeta struct {
	Name               string                 `json:"name,omitempty"`
	Version            string                 `json:"version,omitempty"`
	Description        string                 `json:"description,omitempty"`
	LongDescription    string                 `json:"long_description,omitempty"`
	Maintainer         string                 `json:"maintainer,omitempty"`
	MaintainerEmail    string                 `json:"maintainer_email,omitempty"`
	License            string                 `json:"license,omitempty"`
	Platforms          map[string]interface{} `json:"platforms,omitempty"`
	Depends            map[string]string      `json:"dependencies,omitempty"`
	Reccomends         map[string]string      `json:"recommendations,omitempty"`
	Suggests           map[string]string      `json:"suggestions,omitempty"`
	Conflicts          map[string]string      `json:"conflicting,omitempty"`
	Provides           map[string]interface{} `json:"providing,omitempty"`
	Replaces           map[string]string      `json:"replacing,omitempty"`
	Attributes         map[string]interface{} `json:"attributes,omitempty"` // this has a format as well that could be typed, but blargh https://github.com/lob/chef/blob/master/cookbooks/apache2/metadata.json
	Groupings          map[string]interface{} `json:"groupings,omitempty"`  // never actually seen this used.. looks like it should be map[string]map[string]string, but not sure http://docs.opscode.com/essentials_cookbook_metadata.html
	Recipes            map[string]string      `json:"recipes,omitempty"`
	SourceUrl          string                 `json:"source_url"`
	IssueUrl           string                 `json:"issues_url"`
	ChefVersion        string
	OhaiVersion        string
	Gems               [][]string `json:"gems"`
	EagerLoadLibraries bool       `json:"eager_load_libraries"`
	Privacy            bool       `json:"privacy"`
}

// CookbookAccess represents the permissions on a Cookbook
type CookbookAccess struct {
	Read   bool `json:"read,omitempty"`
	Create bool `json:"create,omitempty"`
	Grant  bool `json:"grant,omitempty"`
	Update bool `json:"update,omitempty"`
	Delete bool `json:"delete,omitempty"`
}

// Cookbook represents the native Go version of the deserialized api cookbook
type Cookbook struct {
	CookbookName string         `json:"cookbook_name"`
	Name         string         `json:"name"`
	Version      string         `json:"version,omitempty"`
	ChefType     string         `json:"chef_type,omitempty"`
	Frozen       bool           `json:"frozen?,omitempty"`
	JsonClass    string         `json:"json_class,omitempty"`
	Files        []CookbookItem `json:"files,omitempty"`
	Templates    []CookbookItem `json:"templates,omitempty"`
	Attributes   []CookbookItem `json:"attributes,omitempty"`
	Recipes      []CookbookItem `json:"recipes,omitempty"`
	Definitions  []CookbookItem `json:"definitions,omitempty"`
	Libraries    []CookbookItem `json:"libraries,omitempty"`
	Providers    []CookbookItem `json:"providers,omitempty"`
	Resources    []CookbookItem `json:"resources,omitempty"`
	RootFiles    []CookbookItem `json:"root_files,omitempty"`
	otherFiles   []CookbookItem
	Metadata     CookbookMeta   `json:"metadata,omitempty"`
	Access       CookbookAccess `json:"access,omitempty"`
}

// Returns a slice of all items in the cookbook
func (c *Cookbook) AllItems() []CookbookItem {
	var allItems []CookbookItem

	allItems = append(allItems, c.Files...)
	allItems = append(allItems, c.Templates...)
	allItems = append(allItems, c.Attributes...)
	allItems = append(allItems, c.Recipes...)
	allItems = append(allItems, c.Definitions...)
	allItems = append(allItems, c.Libraries...)
	allItems = append(allItems, c.Providers...)
	allItems = append(allItems, c.Resources...)
	allItems = append(allItems, c.RootFiles...)
	allItems = append(allItems, c.otherFiles...)

	return allItems
}

// Returns a map of all items in the cookbook keyed by the item checksum
func (c *Cookbook) AllItemsByChecksum() map[string]CookbookItem {
	itemMap := make(map[string]CookbookItem)
	allItems := c.AllItems()

	for _, item := range allItems {
		itemMap[item.Checksum] = item
	}

	return itemMap
}

// String makes CookbookListResult implement the string result
func (c CookbookListResult) String() (out string) {
	for k, v := range c {
		out += fmt.Sprintf("%s => %s\n", k, v.Url)
		for _, i := range v.Versions {
			out += fmt.Sprintf(" * %s\n", i.Version)
		}
	}
	return out
}

// versionParams assembles a querystring for the chef api's  num_versions
// This is used to restrict the number of versions returned in the reponse
func versionParams(path, numVersions string) string {
	if numVersions == "0" {
		numVersions = "all"
	}

	// need to optionally add numVersion args to the request
	if len(numVersions) > 0 {
		path = fmt.Sprintf("%s?num_versions=%s", path, numVersions)
	}
	return path
}

// Get retruns a CookbookVersion for a specific cookbook
//
//	GET /cookbooks/name
func (c *CookbookService) Get(name string) (data CookbookVersion, err error) {
	path := fmt.Sprintf("cookbooks/%s", name)
	err = c.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// GetAvailable returns the versions of a coookbook available on a server
func (c *CookbookService) GetAvailableVersions(name, numVersions string) (data CookbookListResult, err error) {
	path := versionParams(fmt.Sprintf("cookbooks/%s", name), numVersions)
	err = c.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// GetVersion fetches a specific version of a cookbooks data from the server api
//
//	GET /cookbook/foo/1.2.3
//	GET /cookbook/foo/_latest
//	Chef API docs: https://docs.chef.io/api_chef_server.html#cookbooks-name-version
func (c *CookbookService) GetVersion(name, version string) (data Cookbook, err error) {
	url := fmt.Sprintf("cookbooks/%s/%s", name, version)
	err = c.client.magicRequestDecoder("GET", url, nil, &data)
	return
}

// ListVersions lists the cookbooks available on the server limited to numVersions
//
//	Chef API docs: https://docs.chef.io/api_chef_server.html#cookbooks-name
func (c *CookbookService) ListAvailableVersions(numVersions string) (data CookbookListResult, err error) {
	path := versionParams("cookbooks", numVersions)
	err = c.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// ListAllRecipes lists the names of all recipes in the most recent cookbook versions
//
//	Chef API docs: https://docs.chef.io/api_chef_server.html#cookbooks-recipes
func (c *CookbookService) ListAllRecipes() (data CookbookRecipesResult, err error) {
	path := "cookbooks/_recipes"
	err = c.client.magicRequestDecoder("GET", path, nil, &data)
	return
}

// List returns a CookbookListResult with the latest versions of cookbooks available on the server
func (c *CookbookService) List() (CookbookListResult, error) {
	return c.ListAvailableVersions("")
}

// DeleteVersion removes a version of a cook from a server
func (c *CookbookService) Delete(name, version string) (err error) {
	path := fmt.Sprintf("cookbooks/%s/%s", name, version)
	err = c.client.magicRequestDecoder("DELETE", path, nil, nil)
	return
}
func ReadMetaData(path string) (m CookbookMeta, err error) {
	fileName := filepath.Join(path, metaJsonName)
	jsonType := true
	if !isFileExists(fileName) {
		jsonType = false
		fileName = filepath.Join(path, metaRbName)

	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if jsonType {
		return NewMetaDataFromJson(file)
	} else {
		return NewMetaData(string(file))
	}

}
func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
func getKeyValue(str string) (string, []string) {
	c := strings.Split(str, " ")
	if len(c) == 0 {
		return "", nil
	}
	return strings.TrimSpace(c[0]), c[1:]
}
func isFileExists(name string) bool {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func clearWhiteSpace(s []string) (result []string) {
	for _, i := range s {
		if len(i) > 0 {
			result = append(result, i)
		}
	}
	return result
}

func clearComments(s []string) (result []string) {
	for _, i := range s {
		if len(i) > 0 {
			// once a comment is found, break out of line parsing
			if i == "#" {
				break
			}
			result = append(result, i)
		}
	}
	return result
}

func clearQuotesAndCommas(s []string) (result []string) {
	for _, i := range s {
		i = strings.Trim(i, ",")
		i = trimQuotes(i)

		result = append(result, i)
	}
	return result
}

// Creates a new Cookbook object from a given cookbook path.
// Parses in the cookbook's metadata and all cookbook files, respecting ChefIgnore.
func NewCookbookFromPath(cookbookPath string) (*Cookbook, error) {
	cookbook := Cookbook{
		JsonClass: "Chef::CookbookVersion",
		ChefType:  "cookbook_version",
	}

	// Parse cookbook metadata
	meta, err := ReadMetaData(cookbookPath)
	if err != nil {
		return nil, err
	}

	// Update cookbook with metadata information
	cookbook.Version = meta.Version
	cookbook.CookbookName = meta.Name
	cookbook.Name = fmt.Sprintf("%s-%s", cookbook.CookbookName, cookbook.Version)
	cookbook.Metadata = meta

	chefignore := NewChefignore(filepath.Join(cookbookPath, "chefignore"))

	// Find all files in the cookbook and classify them.
	// Allowed types:
	//
	// * Attributes (attributes/)
	// * Definitions (definitions/)
	// * Files (files/)
	// * Libraries (libraries/)
	// * Providers (providers/)
	// * Recipes (recipes/)
	// * Resources (resources/)
	// * RootFiles (/)
	// * Templates (templates/)
	//
	// Gather directories under the root of the cookbook and root files.
	rootPaths, _ := filepath.Glob(filepath.Join(cookbookPath, "*"))
	var rootDirs []string

	for _, rootPath := range rootPaths {
		fileInfo, err := os.Stat(rootPath)
		if err != nil {
			return nil, err
		}

		basePath := filepath.Base(rootPath)

		if chefignore.Ignore(basePath) {
			continue
		}

		if fileInfo.Mode().IsDir() {
			rootDirs = append(rootDirs, basePath)
		} else if fileInfo.Mode().IsRegular() {
			checksum, err := fileMD5Checksum(rootPath)
			if err != nil {
				return nil, err
			}

			// Add the file as a root file
			newItem := CookbookItem{
				Path:        basePath,
				Name:        "root_files/" + basePath,
				Specificity: "default",
				Checksum:    checksum,
			}
			cookbook.RootFiles = append(cookbook.RootFiles, newItem)
		}
	}

	for _, rootDir := range rootDirs {
		fullPath := filepath.Join(cookbookPath, rootDir)
		items, err := walkCookbookDir(fullPath, &chefignore)

		if err != nil {
			return nil, err
		}

		switch rootDir {
		case "attributes":
			cookbook.Attributes = items
		case "definitions":
			cookbook.Definitions = items
		case "files":
			cookbook.Files = items
		case "libraries":
			cookbook.Libraries = items
		case "providers":
			cookbook.Providers = items
		case "recipes":
			cookbook.Recipes = items
		case "resources":
			cookbook.Resources = items
		case "templates":
			cookbook.Templates = items
		default:
			// Store other non-standard files for use with V2 manifests
			cookbook.otherFiles = append(cookbook.otherFiles, items...)
		}

		if err != nil {
			return nil, err
		}
	}

	return &cookbook, nil
}

// Walks a cookbook directory to parse cookbook items from a given path
func walkCookbookDir(dir string, chefignore *Chefignore) ([]CookbookItem, error) {
	baseDir := filepath.Base(dir)
	parentDir := filepath.Dir(dir)
	var items []CookbookItem

	err := filepath.WalkDir(dir, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if file.Type().IsRegular() {
			// Path should be relative to the target directory
			chefPath, err := filepath.Rel(parentDir, path)
			if err != nil {
				return err
			}

			if chefignore.Ignore(chefPath) {
				// Return early for ignored files
				return nil
			}

			// Parse specificity from path if template or cookbook file
			specificity := "default"
			if baseDir == "templates" || baseDir == "files" {
				splitPath := splitCookbookDir(chefPath)
				if len(splitPath) == 2 {
					specificity = "root_default"
				} else if len(splitPath) > 2 {
					specificity = splitPath[1]
				}
			}

			itemChecksum, err := fileMD5Checksum(path)
			if err != nil {
				return err
			}

			newItem := CookbookItem{
				Path:        chefPath,
				Name:        baseDir + "/" + filepath.Base(path),
				Specificity: specificity,
				Checksum:    itemChecksum,
			}
			items = append(items, newItem)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil
}

// Splits the cookbook directory into a slice of strings by directory separator
func splitCookbookDir(path string) []string {
	dir, last := filepath.Split(path)
	if dir == "" {
		return []string{last}
	}
	return append(splitCookbookDir(filepath.Clean(dir)), last)
}

func NewMetaData(data string) (m CookbookMeta, err error) {
	linesData := strings.Split(data, "\n")
	if len(linesData) < 3 {
		return m, errors.New("not much info")
	}
	m.Depends = make(map[string]string, 1)
	m.Platforms = make(map[string]interface{}, 1)
	for _, i := range linesData {
		key, value := getKeyValue(strings.TrimSpace(i))
		if fn, ok := metaRegistry[key]; ok {
			err = fn(value, &m)
			if err != nil {
				return
			}
		}
	}
	return m, err
}

func NewMetaDataFromJson(data []byte) (m CookbookMeta, err error) {
	err = json.Unmarshal(data, &m)
	return m, err
}

func StringParserForMeta(s []string) string {
	str := strings.Join(s, " ")
	return trimQuotes(strings.TrimSpace(str))
}
func metaNameParser(s []string, m *CookbookMeta) error {
	m.Name = StringParserForMeta(s)
	return nil
}
func metaMaintainerParser(s []string, m *CookbookMeta) error {
	m.Maintainer = StringParserForMeta(s)
	return nil
}
func metaMaintainerMailParser(s []string, m *CookbookMeta) error {
	m.MaintainerEmail = StringParserForMeta(s)
	return nil
}
func metaLicenseParser(s []string, m *CookbookMeta) error {
	m.License = StringParserForMeta(s)
	return nil
}
func metaDescriptionParser(s []string, m *CookbookMeta) error {
	m.Description = StringParserForMeta(s)
	return nil
}
func metaLongDescriptionParser(s []string, m *CookbookMeta) error {
	m.LongDescription = StringParserForMeta(s)
	return nil
}
func metaIssueUrlParser(s []string, m *CookbookMeta) error {
	m.IssueUrl = StringParserForMeta(s)
	return nil
}
func metaSourceUrlParser(s []string, m *CookbookMeta) error {
	m.SourceUrl = StringParserForMeta(s)
	return nil
}
func metaGemParser(s []string, m *CookbookMeta) error {
	s = clearWhiteSpace(s)
	s = clearComments(s)
	s = clearQuotesAndCommas(s)

	m.Gems = append(m.Gems, s)
	return nil
}

func metaVersionParser(s []string, m *CookbookMeta) error {
	m.Version = StringParserForMeta(s)
	return nil
}
func metaOhaiVersionParser(s []string, m *CookbookMeta) error {
	m.OhaiVersion = StringParserForMeta(s)
	return nil
}
func metaChefVersionParser(s []string, m *CookbookMeta) error {
	m.ChefVersion = StringParserForMeta(s)
	return nil
}
func metaPrivacyParser(s []string, m *CookbookMeta) error {
	if s[0] == "true" {
		m.Privacy = true
	}
	return nil
}
func metaSupportsParser(s []string, m *CookbookMeta) error {
	s = clearWhiteSpace(s)
	s = clearComments(s)

	// Remove surrounding spaces, commas, and quotes from keys
	k := strings.TrimSpace(s[0])
	k = strings.Trim(k, ",")
	k = trimQuotes(k)

	switch len(s) {
	case 1:
		if s[0] != "os" {
			m.Platforms[k] = ">= 0.0.0"
		}
	case 2:
		m.Platforms[k] = s[1]
	case 3:
		v := trimQuotes(s[1] + " " + s[2])
		m.Platforms[k] = v

	}
	if len(s) > 3 {
		return errors.New(`<<~OBSOLETED
		The supports specification syntax you are using is no longer valid. You may not
		specify more than one version constraint for a particular supported platform.
			Consult https://docs.chef.io/config_rb_metadata/ for the updated syntax.`)
	}
	return nil
}
func metaDependsParser(s []string, m *CookbookMeta) error {
	s = clearWhiteSpace(s)
	s = clearComments(s)

	// Remove surrounding spaces, commas, and quotes from keys
	k := strings.TrimSpace(s[0])
	k = strings.Trim(k, ",")
	k = trimQuotes(k)

	switch len(s) {
	case 1:
		m.Depends[k] = ">= 0.0.0"
	case 2:
		m.Depends[k] = s[1]

	case 3:
		v := trimQuotes(s[1] + " " + s[2])
		m.Depends[k] = v

	}
	if len(s) > 3 {
		return errors.New(`<<~OBSOLETED
		The dependency specification syntax you are using is no longer valid. You may not
		specify more than one version constraint for a particular cookbook.
			Consult https://docs.chef.io/config_rb_metadata/ for the updated syntax.`)
	}
	return nil
}

func metaSupportsRubyParser(s []string, m *CookbookMeta) error {
	if len(s) > 1 {
		for _, i := range s {
			switch i {
			case ").each":
				continue
			case "do":
				continue
			case "|os|":
				continue
			default:
				m.Platforms[strings.TrimSpace(s[0])] = ">= 0.0.0"
			}
		}
	}
	return nil
}

func init() {
	metaRegistry = make(map[string]metaFunc, 15)
	metaRegistry["name"] = metaNameParser
	metaRegistry["maintainer"] = metaMaintainerParser
	metaRegistry["maintainer_email"] = metaMaintainerMailParser
	metaRegistry["license"] = metaLicenseParser
	metaRegistry["description"] = metaDescriptionParser
	metaRegistry["long_description"] = metaLongDescriptionParser
	metaRegistry["source_url"] = metaSourceUrlParser
	metaRegistry["issues_url"] = metaIssueUrlParser
	metaRegistry["platforms"] = metaSupportsParser
	metaRegistry["supports"] = metaSupportsParser
	metaRegistry["%w("] = metaSupportsRubyParser
	metaRegistry["privacy"] = metaPrivacyParser
	metaRegistry["depends"] = metaDependsParser
	metaRegistry["version"] = metaVersionParser
	metaRegistry["chef_version"] = metaChefVersionParser
	metaRegistry["ohai_version"] = metaOhaiVersionParser
	metaRegistry["gem"] = metaGemParser
}
