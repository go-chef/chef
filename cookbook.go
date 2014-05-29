package chef

import (
	"github.com/mitchellh/mapstructure"
)

// Cookbook is the chef-cookbook container
type Cookbook struct {
	*Reader
	*nativeCookbook
}

// Each Cookbook lists it's files as a cookbookItem. This structure captures those and makes it easier to work with cook data
type CookbookItem struct {
	Url         string `mapstructure:"url"`
	Path        string `mapstructure:"path"`
	Name        string `mapstructure:"name"`
	Checksum    string `mapstructure:"checksum"`
	Specificity string `mapstructure:"specificity"`
}

// CookbookMeta represents a Golang version of cookbook metadata
type CookbookMeta struct {
	Name            string                 `mapstructure:"cookbook_name"`
	Version         string                 `mapstructure:"version"`
	Description     string                 `mapstructure:"description"`
	LongDescription string                 `mapstructure:"long_description"`
	Maintainer      string                 `mapstructure:"maintainer"`
	MaintainerEmail string                 `mapstructure:"maintainer_email"`
	License         string                 `mapstructure:"license"`
	Platforms       map[string]string      `mapstructure:"platforms"`
	Depends         map[string]string      `mapstructure:"dependencies"`
	Reccomends      map[string]string      `mapstructure:"recommendations"`
	Suggests        map[string]string      `mapstructure:"suggestions"`
	Conflicts       map[string]string      `mapstructure:"conflicting"`
	Provides        map[string]string      `mapstructure:"providing"`
	Replaces        map[string]string      `mapstructure:"replacing"`
	attributes      map[string]interface{} `mapstructure:"attributes"` // this has a format as well that could be typed, but blargh https://github.com/lob/chef/blob/master/cookbooks/apache2/metadata.json
	groupings       map[string]interface{} `mapstructure:"groupings"`  // never actually seen this used.. looks like it should be map[string]map[string]string, but not sure http://docs.opscode.com/essentials_cookbook_metadata.html
	recipes         map[string]string      `mapstructure:"recipes"`
}

// NativeNode represents the native Go version of the deserialized cookbook
type nativeCookbook struct {
	Name         string         `mapstructure:"name"`
	Version      string         `mapstructure:"version"`
	ChefType     string         `mapstructure:"chef_type"`
	Frozen       bool           `mapstructure:"frozen?"`
	JsonClass    string         `mapstructure:"json_class"`
	CookbookName string         `mapstructure:"cookbook_name"` // yes the json can have this :\
	Files        []CookbookItem `mapstructure:"files"`
	Templates    []CookbookItem `mapstructure:"Templates"`
	Attributes   []CookbookItem `mapstructure:"attributes"`
	Recipes      []CookbookItem `mapstructure:"recipes"`
	Definitions  []CookbookItem `mapstructure:"definitions"`
	Libraries    []CookbookItem `mapstructure:"libraries"`
	Providers    []CookbookItem `mapstructure:"Providers"`
	Resources    []CookbookItem `mapstructure:"Resources"`
	RootFiles    []CookbookItem `mapstructure:"Templates"`
	Metadata     CookbookMeta   `mapstructure:"Metadata"`
}

// NewCookbook is used to create a cookbook from a Reader
func NewCookbook(reader *Reader) (*Cookbook, error) {
	cook := Cookbook{reader, &nativeCookbook{}}
	if err := mapstructure.Decode(reader, cook.nativeCookbook); err != nil {
		return nil, err
	}
	return &cook, nil
}
