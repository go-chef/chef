package chef

import (
	"github.com/mitchellh/mapstructure"
)

// Cookbook is the chef-cookbook container
type CookBook struct {
	*Reader
	*nativeCookBook
}

// Each Cookbook lists it's files as a cookbookItem. This structure captures those and makes it easier to work with cook data
type CookBookItem struct {
	Url         string `mapstructure:"url"`
	Path        string `mapstructure:"path"`
	Name        string `mapstructure:"name"`
	Checksum    string `mapstructure:"checksum"`
	Specificity string `mapstructure:"specificity"`
}

// CookBookMeta represents a Golang version of cookbook metadata
type CookBookMeta struct {
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
type nativeCookBook struct {
	Name         string         `mapstructure:"name"`
	Version      string         `mapstructure:"version"`
	ChefType     string         `mapstructure:"chef_type"`
	Frozen       bool           `mapstructure:"frozen?"`
	JsonClass    string         `mapstructure:"json_class"`
	CookbookName string         `mapstructure:"cookbook_name"` // yes the json can have this :\
	Files        []CookBookItem `mapstructure:"files"`
	Templates    []CookBookItem `mapstructure:"Templates"`
	Attributes   []CookBookItem `mapstructure:"attributes"`
	Recipes      []CookBookItem `mapstructure:"recipes"`
	Definitions  []CookBookItem `mapstructure:"definitions"`
	Libraries    []CookBookItem `mapstructure:"libraries"`
	Providers    []CookBookItem `mapstructure:"Providers"`
	Resources    []CookBookItem `mapstructure:"Resources"`
	RootFiles    []CookBookItem `mapstructure:"Templates"`
	Metadata     CookBookMeta   `mapstructure:"Metadata"`
}

// NewCookBook is used to create a cookbook from a Reader
func NewCookBook(reader *Reader) (*CookBook, error) {
	cook := CookBook{reader, &nativeCookBook{}}
	if err := mapstructure.Decode(reader, cook.nativeCookBook); err != nil {
		return nil, err
	}
	return &cook, nil
}
