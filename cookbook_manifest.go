package chef

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
)

// Cookbook Manifest meta format for Chef API versions 0-2
type CookbookManifestMeta struct {
	Name               string                 `json:"name,omitempty"`
	Version            string                 `json:"version,omitempty"`
	Description        string                 `json:"description,omitempty"`
	LongDescription    string                 `json:"long_description"`
	Maintainer         string                 `json:"maintainer,omitempty"`
	MaintainerEmail    string                 `json:"maintainer_email,omitempty"`
	License            string                 `json:"license,omitempty"`
	Platforms          map[string]interface{} `json:"platforms"`
	Depends            map[string]string      `json:"dependencies"`
	Recommends         map[string]string      `json:"recommendations,omitempty"`
	Suggests           map[string]string      `json:"suggestions,omitempty"`
	Conflicts          map[string]string      `json:"conflicting,omitempty"`
	Provides           map[string]interface{} `json:"providing,omitempty"`
	Replaces           map[string]string      `json:"replacing,omitempty"`
	Attributes         map[string]interface{} `json:"attributes,omitempty"`
	Groupings          map[string]interface{} `json:"groupings,omitempty"`
	Recipes            map[string]string      `json:"recipes,omitempty"`
	SourceUrl          string                 `json:"source_url"`
	IssueUrl           string                 `json:"issues_url"`
	ChefVersions       [][]string             `json:"chef_versions"`
	OhaiVersions       [][]string             `json:"ohai_versions"`
	Gems               [][]string             `json:"gems"`
	EagerLoadLibraries bool                   `json:"eager_load_libraries"`
	Privacy            bool                   `json:"privacy"`
}

// Cookbook Manifest format for Chef API versions 0-1
type CookbookManifestV0 struct {
	CookbookName string               `json:"cookbook_name"`
	Name         string               `json:"name"`
	Version      string               `json:"version,omitempty"`
	ChefType     string               `json:"chef_type,omitempty"`
	Frozen       bool                 `json:"frozen?"`
	JsonClass    string               `json:"json_class,omitempty"`
	Files        []CookbookItem       `json:"files,omitempty"`
	Templates    []CookbookItem       `json:"templates,omitempty"`
	Attributes   []CookbookItem       `json:"attributes,omitempty"`
	Recipes      []CookbookItem       `json:"recipes,omitempty"`
	Definitions  []CookbookItem       `json:"definitions,omitempty"`
	Libraries    []CookbookItem       `json:"libraries,omitempty"`
	Providers    []CookbookItem       `json:"providers,omitempty"`
	Resources    []CookbookItem       `json:"resources,omitempty"`
	RootFiles    []CookbookItem       `json:"root_files,omitempty"`
	Metadata     CookbookManifestMeta `json:"metadata"`
}

// Cookbook Manifest format for Chef API version 2
type CookbookManifestV2 struct {
	CookbookName string               `json:"cookbook_name"`
	Name         string               `json:"name"`
	Version      string               `json:"version,omitempty"`
	ChefType     string               `json:"chef_type,omitempty"`
	Frozen       bool                 `json:"frozen?"`
	JsonClass    string               `json:"json_class,omitempty"`
	AllFiles     []CookbookItem       `json:"all_files"`
	Metadata     CookbookManifestMeta `json:"metadata"`
}

// Generate a CookbookManifestMeta from an existing CookbookMeta object
func (meta *CookbookMeta) Manifest() *CookbookManifestMeta {
	chefVersions := make([][]string, 0)
	ohaiVersions := make([][]string, 0)

	if meta.ChefVersion != "" {
		chefVersionList := []string{meta.ChefVersion}
		chefVersions = append(chefVersions, chefVersionList)
	}

	if meta.OhaiVersion != "" {
		ohaiVersionList := []string{meta.OhaiVersion}
		ohaiVersions = append(ohaiVersions, ohaiVersionList)
	}

	if meta.Gems == nil {
		// Ensure the gems array exists, nil will break the client
		meta.Gems = make([][]string, 0)
	}

	// Set eager load libraries to true since this is default behavior
	meta.EagerLoadLibraries = true

	manifestMeta := CookbookManifestMeta{
		Name:               meta.Name,
		Version:            meta.Version,
		Description:        meta.Description,
		LongDescription:    meta.LongDescription,
		Maintainer:         meta.Maintainer,
		MaintainerEmail:    meta.MaintainerEmail,
		License:            meta.License,
		Platforms:          meta.Platforms,
		Depends:            meta.Depends,
		Provides:           meta.Provides,
		Recipes:            meta.Recipes,
		Recommends:         meta.Reccomends,
		Suggests:           meta.Suggests,
		Conflicts:          meta.Conflicts,
		Replaces:           meta.Replaces,
		Attributes:         meta.Attributes,
		Groupings:          meta.Groupings,
		SourceUrl:          meta.SourceUrl,
		IssueUrl:           meta.IssueUrl,
		Gems:               meta.Gems,
		EagerLoadLibraries: meta.EagerLoadLibraries,
		Privacy:            meta.Privacy,
		ChefVersions:       chefVersions,
		OhaiVersions:       ohaiVersions,
	}

	return &manifestMeta
}

func (c *Cookbook) populateProvidesMetadata() {
	if c.Metadata.Provides == nil {
		c.Metadata.Provides = make(map[string]interface{})

		for _, recipe := range c.Recipes {
			recipeName := strings.Split(filepath.Base(recipe.Name), ".")[0]
			recipeKey := fmt.Sprintf("%s::%s", c.CookbookName, recipeName)

			if recipeName == "default" {
				recipeKey = c.CookbookName
			}

			c.Metadata.Provides[recipeKey] = ">= 0.0.0"
		}
	}
}

func (c *Cookbook) populateRecipesMetadata() {
	if c.Metadata.Recipes == nil {
		c.Metadata.Recipes = make(map[string]string)

		for _, recipe := range c.Recipes {
			recipeName := strings.Split(filepath.Base(recipe.Name), ".")[0]
			recipeKey := fmt.Sprintf("%s::%s", c.CookbookName, recipeName)

			if recipeName == "default" {
				recipeKey = c.CookbookName
			}

			c.Metadata.Recipes[recipeKey] = ""
		}
	}
}

// Generate a CookbookManifestV0 from an existing Cookbook object
func (c *Cookbook) ManifestV0() *CookbookManifestV0 {
	c.populateProvidesMetadata()
	c.populateRecipesMetadata()

	manifest := CookbookManifestV0{
		CookbookName: c.CookbookName,
		Name:         c.Name,
		Version:      c.Version,
		ChefType:     "cookbook_version",
		Frozen:       false,
		JsonClass:    "Chef::CookbookVersion",
		Metadata:     *c.Metadata.Manifest(),
	}

	for i := range c.Files {
		if manifest.Files == nil {
			manifest.Files = make([]CookbookItem, len(c.Files))
		}
		manifest.Files[i] = *manifestV0CookbookItem(&c.Files[i])
	}

	for i := range c.Templates {
		if manifest.Templates == nil {
			manifest.Templates = make([]CookbookItem, len(c.Templates))
		}
		manifest.Templates[i] = *manifestV0CookbookItem(&c.Templates[i])
	}

	for i := range c.Attributes {
		if manifest.Attributes == nil {
			manifest.Attributes = make([]CookbookItem, len(c.Attributes))
		}
		manifest.Attributes[i] = *manifestV0CookbookItem(&c.Attributes[i])
	}

	for i := range c.Recipes {
		if manifest.Recipes == nil {
			manifest.Recipes = make([]CookbookItem, len(c.Recipes))
		}
		manifest.Recipes[i] = *manifestV0CookbookItem(&c.Recipes[i])
	}

	for i := range c.Definitions {
		if manifest.Definitions == nil {
			manifest.Definitions = make([]CookbookItem, len(c.Definitions))
		}
		manifest.Definitions[i] = *manifestV0CookbookItem(&c.Definitions[i])
	}

	for i := range c.Libraries {
		if manifest.Libraries == nil {
			manifest.Libraries = make([]CookbookItem, len(c.Libraries))
		}
		manifest.Libraries[i] = *manifestV0CookbookItem(&c.Libraries[i])
	}

	for i := range c.Providers {
		if manifest.Providers == nil {
			manifest.Providers = make([]CookbookItem, len(c.Providers))
		}
		manifest.Providers[i] = *manifestV0CookbookItem(&c.Providers[i])
	}

	for i := range c.Resources {
		if manifest.Resources == nil {
			manifest.Resources = make([]CookbookItem, len(c.Resources))
		}
		manifest.Resources[i] = *manifestV0CookbookItem(&c.Resources[i])
	}

	for i := range c.RootFiles {
		if manifest.RootFiles == nil {
			manifest.RootFiles = make([]CookbookItem, len(c.RootFiles))
		}
		manifest.RootFiles[i] = *manifestV0CookbookItem(&c.RootFiles[i])
	}

	return &manifest
}

// Generate a new cookbook item that follows ManifestV0 naming syntax
func manifestV0CookbookItem(item *CookbookItem) *CookbookItem {
	return &CookbookItem{
		Name:        filepath.Base(item.Name),
		Path:        item.Path,
		Url:         item.Url,
		Checksum:    item.Checksum,
		Specificity: item.Specificity,
	}
}

// Generate a CookbookManifestV2 from an existing Cookbook object
func (c *Cookbook) ManifestV2() *CookbookManifestV2 {
	c.populateProvidesMetadata()
	c.populateRecipesMetadata()

	manifest := CookbookManifestV2{
		CookbookName: c.CookbookName,
		Name:         c.Name,
		Version:      c.Version,
		ChefType:     "cookbook_version",
		Frozen:       c.Frozen,
		JsonClass:    "Chef::CookbookVersion",
		AllFiles:     c.AllItems(),
		Metadata:     *c.Metadata.Manifest(),
	}

	return &manifest
}

func (c *Cookbook) ManifestJsonForApi(serverApiVersion string) (reader io.Reader, err error) {
	apiVersionInt, _ := strconv.Atoi(serverApiVersion)

	if apiVersionInt >= 2 {
		reader, err = JSONReader(c.ManifestV2())
	} else {
		reader, err = JSONReader(c.ManifestV0())
	}

	if err != nil {
		return nil, err
	}

	return reader, nil
}
