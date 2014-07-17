package chef

import "errors"
import "fmt"

type CookbookService struct {
	client *Client
}

// Each Cookbook lists it's files as a cookbookItem. This structure captures those and makes it easier to work with cook data
type CookbookItem struct {
	Url         string `json:"url,omitempty"`
	Path        string `json:"path,omitempty"`
	Name        string `json:"name,omitempty"`
	Checksum    string `json:"checksum,omitempty"`
	Specificity string `json:"specificity,omitempty"`
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
	Url     string `json:"url,omitempty"`
	Version string `json:"version,omitempty"`
}

// totally looks better
type CookbookVersions struct {
	Url      string            `json:"url,omitempty"`
	Versions []CookbookVersion `json:"versions,omitempty"`
}

// CookbookMeta represents a Golang version of cookbook metadata
type CookbookMeta struct {
	Name            string                 `json:"cookbook_name,omitempty"`
	Version         string                 `json:"version,omitempty"`
	Description     string                 `json:"description,omitempty"`
	LongDescription string                 `json:"long_description,omitempty"`
	Maintainer      string                 `json:"maintainer,omitempty"`
	MaintainerEmail string                 `json:"maintainer_email,omitempty"`
	License         string                 `json:"license,omitempty"`
	Platforms       map[string]string      `json:"platforms,omitempty"`
	Depends         map[string]string      `json:"dependencies,omitempty"`
	Reccomends      map[string]string      `json:"recommendations,omitempty"`
	Suggests        map[string]string      `json:"suggestions,omitempty"`
	Conflicts       map[string]string      `json:"conflicting,omitempty"`
	Provides        map[string]string      `json:"providing,omitempty"`
	Replaces        map[string]string      `json:"replacing,omitempty"`
	attributes      map[string]interface{} `json:"attributes,omitempty"` // this has a format as well that could be typed, but blargh https://github.com/lob/chef/blob/master/cookbooks/apache2/metadata.json
	groupings       map[string]interface{} `json:"groupings,omitempty"`  // never actually seen this used.. looks like it should be map[string]map[string]string, but not sure http://docs.opscode.com/essentials_cookbook_metadata.html
	recipes         map[string]string      `json:"recipes,omitempty"`
}

// NativeNode represents the native Go version of the deserialized cookbook
type Cookbook struct {
	Name         string         `json:"name"`
	Version      string         `json:"version,omitempty"`
	ChefType     string         `json:"chef_type,omitempty"`
	Frozen       bool           `json:"frozen?,omitempty"`
	JsonClass    string         `json:"json_class,omitempty"`
	CookbookName string         `json:"cookbook_name,omitempty"` // yes the json can have this :\
	Files        []CookbookItem `json:"files,omitempty"`
	Templates    []CookbookItem `json:"Templates,omitempty"`
	Attributes   []CookbookItem `json:"attributes,omitempty"`
	Recipes      []CookbookItem `json:"recipes,omitempty"`
	Definitions  []CookbookItem `json:"definitions,omitempty"`
	Libraries    []CookbookItem `json:"libraries,omitempty"`
	Providers    []CookbookItem `json:"Providers,omitempty"`
	Resources    []CookbookItem `json:"Resources,omitempty"`
	RootFiles    []CookbookItem `json:"Templates,omitempty"`
	Metadata     CookbookMeta   `json:"Metadata,omitempty"`
}

// Get - Grabs a cookbook from a chef api and populates the object
//  GET /cookbooks/name
func (c *CookbookService) Get(name string, numVersions int) (CookbookVersion, error) {
	return CookbookVersion{}, nil
}

// GetVersion
// GET /cookbook/foo/1.2.3
// GET /cookbook/foo/_latest
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

// ListVersions lists the cookbooks available on the server limited to numVersions
// Chef API docs: http://docs.opscode.com/api_chef_server.html#id2
func (c *CookbookService) ListVersions(numVersions string) (data *CookbookListResult, err error) {
	if numVersions == "0" {
		return nil, errors.New("numVersions of 0 are not allowed by the chef server api")
	}
	// need to optionally add numVersion args to the request
	u := "cookbooks"
	if len(numVersions) > 0 {
		u = fmt.Sprintf("%s?num_versions=%s", u, numVersions)
	}

	req, err := c.client.MakeRequest("GET", u, nil)
	if err != nil {
		return
	}

	_, err = c.client.Do(req, &data)
	if err != nil {
		return
	}

	return
}

// List returns a CookbookListResult with the latest versions of cookbooks available on the server
func (c *CookbookService) List() (*CookbookListResult, error) {
	return c.ListVersions("")
}
