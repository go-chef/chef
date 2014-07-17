package chef

import "fmt"

type CookbookService struct {
	client *Client
}

// Each Cookbook lists it's files as a cookbookItem. This structure captures those and makes it easier to work with cook data
type CookbookItem struct {
	Url         string `json:"url"`
	Path        string `json:"path"`
	Name        string `json:"name"`
	Checksum    string `json:"checksum"`
	Specificity string `json:"specificity"`
}

// CookbookListResult is the summary info returned by chef-api when listing
// http://docs.opscode.com/api_chef_server.html#cookbooks
// {
//  "apache2" => {
//    "url" => "http://localhost:4000/cookbooks/apache2",
//    "versions" => [
//      {"url" => "http://localhost:4000/cookbooks/apache2/5.1.0",
//       "version" => "5.1.0"},
//      {"url" => "http://localhost:4000/cookbooks/apache2/4.2.0",
//       "version" => "4.2.0"}
//    ]
//  }
//}
type CookbookListResult map[string]CookbookVersion

type CookbookVersion struct {
	Url      string                     `json:"url"`
	Versions map[string]CookbookVersion `json:"version"`
}

// CookbookMeta represents a Golang version of cookbook metadata
type CookbookMeta struct {
	Name            string                 `json:"cookbook_name"`
	Version         string                 `json:"version"`
	Description     string                 `json:"description,omitempty"`
	LongDescription string                 `json:"long_description,omitempty"`
	Maintainer      string                 `json:"maintainer,omitempty"`
	MaintainerEmail string                 `json:"maintainer_email"`
	License         string                 `json:"license"`
	Platforms       map[string]string      `json:"platforms"`
	Depends         map[string]string      `json:"dependencies"`
	Reccomends      map[string]string      `json:"recommendations"`
	Suggests        map[string]string      `json:"suggestions"`
	Conflicts       map[string]string      `json:"conflicting"`
	Provides        map[string]string      `json:"providing"`
	Replaces        map[string]string      `json:"replacing"`
	attributes      map[string]interface{} `json:"attributes"` // this has a format as well that could be typed, but blargh https://github.com/lob/chef/blob/master/cookbooks/apache2/metadata.json
	groupings       map[string]interface{} `json:"groupings"`  // never actually seen this used.. looks like it should be map[string]map[string]string, but not sure http://docs.opscode.com/essentials_cookbook_metadata.html
	recipes         map[string]string      `json:"recipes"`
}

// NativeNode represents the native Go version of the deserialized cookbook
type Cookbook struct {
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	ChefType     string         `json:"chef_type"`
	Frozen       bool           `json:"frozen?"`
	JsonClass    string         `json:"json_class"`
	CookbookName string         `json:"cookbook_name"` // yes the json can have this :\
	Files        []CookbookItem `json:"files"`
	Templates    []CookbookItem `json:"Templates"`
	Attributes   []CookbookItem `json:"attributes"`
	Recipes      []CookbookItem `json:"recipes"`
	Definitions  []CookbookItem `json:"definitions"`
	Libraries    []CookbookItem `json:"libraries"`
	Providers    []CookbookItem `json:"Providers"`
	Resources    []CookbookItem `json:"Resources"`
	RootFiles    []CookbookItem `json:"Templates"`
	Metadata     CookbookMeta   `json:"Metadata"`
}

// Get - Grabs a cookbook from a chef api and populates the object
func (c *CookbookService) Get(name string, numVersions int) (CookbookVersion, error) {
	return CookbookVersion{}, nil
}

// GetVersion
// /cookbook/foo/1.2.3
// /cookbook/foo/_latest
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id5
func (c *CookbookService) GetVersion(name string, version string) (data *Cookbook, err error) {
	url := fmt.Sprintf("cookbooks/%s/%s", name, version)
	req, err := c.client.MakeRequest("GET", url, nil)
	if err != nil {
		return
	}

	_, err = c.client.Do(req, &data)
	if err != nil {
		return
	}

	return
}

// List lists the cookbooks on the Chef server.
// TODO: Support num_versions parameter
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id2
func (c *CookbookService) List() (data *CookbookListResult, err error) {
	// num_versions=n
	req, err := c.client.MakeRequest("GET", "cookbooks", nil)
	if err != nil {
		return
	}

	_, err = c.client.Do(req, &data)
	if err != nil {
		return
	}

	return
}
