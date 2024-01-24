package chef

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const _sourceManifestCookbookPath = "test/cookbooks/testdeps"
const _sourceComplexManifestCookbookPath = "test/cookbooks/testcomplex"

func TestCookbookMetaManifest(t *testing.T) {
	meta, err := ReadMetaData(_sourceManifestCookbookPath)
	assert.Nil(t, err)

	metaManifest := meta.Manifest()

	assert.Equal(t, &CookbookManifestMeta{
		Name:            "testdeps",
		Version:         "0.1.0",
		Description:     "Installs/Configures testdeps",
		Maintainer:      "The Authors",
		MaintainerEmail: "you@example.com",
		License:         "All Rights Reserved",
		Platforms: map[string]interface{}{
			"redhat": ">= 0.0.0",
			"ubuntu": ">= 20.04",
		},
		Depends: map[string]string{
			"lvm":     "~> 6.1",
			"vagrant": ">= 4.0.14",
		},
		ChefVersions: [][]string{
			{">= 18.0"},
		},
		OhaiVersions: [][]string{},
		Gems: [][]string{
			{"json", ">1.0.0"},
		},
		EagerLoadLibraries: true,
	}, metaManifest)
}

func TestCookbookManifestV0(t *testing.T) {
	cookbook, err := NewCookbookFromPath(_sourceManifestCookbookPath)
	assert.Nil(t, err)

	manifest := cookbook.ManifestV0()

	assert.Equal(t, &CookbookManifestV0{
		CookbookName: "testdeps",
		Name:         "testdeps-0.1.0",
		Version:      "0.1.0",
		ChefType:     "cookbook_version",
		Frozen:       false,
		JsonClass:    "Chef::CookbookVersion",
		Files:        nil,
		Templates:    nil,
		Attributes: []CookbookItem{{
			Name:        "default.rb",
			Path:        "attributes/default.rb",
			Checksum:    "553637b4fba46b5148f88d6dd3877e2f",
			Specificity: "default",
		}},
		Recipes: []CookbookItem{{
			Name:        "default.rb",
			Path:        "recipes/default.rb",
			Checksum:    "4e15b1e5593d717685323c5dac86b99e",
			Specificity: "default",
		}},
		Definitions: nil,
		Libraries:   nil,
		Providers:   nil,
		Resources:   nil,
		RootFiles: []CookbookItem{{
			Name:        "metadata.rb",
			Path:        "metadata.rb",
			Checksum:    "ba208d0ffc0dd8cbe9c71fb40fb207b2",
			Specificity: "default",
		}},
		Metadata: *cookbook.Metadata.Manifest(),
	}, manifest)

	var expectedManifest *CookbookManifestV0
	expectedMetadataBytes, err := os.ReadFile(filepath.Join(_sourceComplexManifestCookbookPath, "manifests", "chefv0.json"))
	assert.Nil(t, err)
	json.Unmarshal(expectedMetadataBytes, &expectedManifest)

	cookbook, err = NewCookbookFromPath(_sourceComplexManifestCookbookPath)
	assert.Nil(t, err)

	parsedManifest := cookbook.ManifestV0()

	assert.Equal(t, expectedManifest, parsedManifest)
}

func TestCookbookManifestV2(t *testing.T) {
	cookbook, err := NewCookbookFromPath(_sourceManifestCookbookPath)
	assert.Nil(t, err)

	manifest := cookbook.ManifestV2()

	assert.Equal(t, &CookbookManifestV2{
		CookbookName: "testdeps",
		Name:         "testdeps-0.1.0",
		Version:      "0.1.0",
		ChefType:     "cookbook_version",
		Frozen:       false,
		JsonClass:    "Chef::CookbookVersion",
		AllFiles: []CookbookItem{
			{
				Name:        "attributes/default.rb",
				Path:        "attributes/default.rb",
				Checksum:    "553637b4fba46b5148f88d6dd3877e2f",
				Specificity: "default",
			},
			{
				Name:        "recipes/default.rb",
				Path:        "recipes/default.rb",
				Checksum:    "4e15b1e5593d717685323c5dac86b99e",
				Specificity: "default",
			},
			{
				Name:        "root_files/metadata.rb",
				Path:        "metadata.rb",
				Checksum:    "ba208d0ffc0dd8cbe9c71fb40fb207b2",
				Specificity: "default",
			},
		},
		Metadata: *cookbook.Metadata.Manifest(),
	}, manifest)

	var expectedManifest *CookbookManifestV2
	expectedMetadataBytes, err := os.ReadFile(filepath.Join(_sourceComplexManifestCookbookPath, "manifests", "chefv2.json"))
	assert.Nil(t, err)
	json.Unmarshal(expectedMetadataBytes, &expectedManifest)

	cookbook, err = NewCookbookFromPath(_sourceComplexManifestCookbookPath)
	assert.Nil(t, err)

	parsedManifest := cookbook.ManifestV2()

	assert.Equal(t, expectedManifest.ChefType, parsedManifest.ChefType)
	assert.Equal(t, expectedManifest.CookbookName, parsedManifest.CookbookName)
	assert.Equal(t, expectedManifest.Frozen, parsedManifest.Frozen)
	assert.Equal(t, expectedManifest.JsonClass, parsedManifest.JsonClass)
	assert.Equal(t, expectedManifest.Name, parsedManifest.Name)
	assert.Equal(t, expectedManifest.Version, parsedManifest.Version)
	assert.Equal(t, expectedManifest.Metadata, parsedManifest.Metadata)

	assert.ElementsMatch(t, expectedManifest.AllFiles, parsedManifest.AllFiles)
}

func TestCookbookManifestJsonForApi(t *testing.T) {
	cookbook, err := NewCookbookFromPath(_sourceManifestCookbookPath)
	assert.Nil(t, err)

	// Ensure API versions 0-1 use the same manifest format
	v0ManifestJsonReader, err := cookbook.ManifestJsonForApi("0")
	assert.Nil(t, err)
	v1ManifestJsonReader, err := cookbook.ManifestJsonForApi("1")
	assert.Nil(t, err)

	v0ManifestJson, err := io.ReadAll(v0ManifestJsonReader)
	assert.Nil(t, err)
	v1ManifestJson, err := io.ReadAll(v1ManifestJsonReader)
	assert.Nil(t, err)
	assert.Equal(t, v0ManifestJson, v1ManifestJson)

	// Ensure API versions 2+ use the same manifest format
	v2ManifestJsonReader, err := cookbook.ManifestJsonForApi("2")
	assert.Nil(t, err)
	v3ManifestJsonReader, err := cookbook.ManifestJsonForApi("3")
	assert.Nil(t, err)

	v2ManifestJson, err := io.ReadAll(v2ManifestJsonReader)
	assert.Nil(t, err)
	v3ManifestJson, err := io.ReadAll(v3ManifestJsonReader)
	assert.Nil(t, err)
	assert.Equal(t, v2ManifestJson, v3ManifestJson)

	// Verify JSON content
	var parsedV0Json CookbookManifestV0
	err = json.Unmarshal(v0ManifestJson, &parsedV0Json)
	assert.Nil(t, err)

	assert.Equal(t, CookbookManifestV0{
		CookbookName: "testdeps",
		Name:         "testdeps-0.1.0",
		Version:      "0.1.0",
		ChefType:     "cookbook_version",
		Frozen:       false,
		JsonClass:    "Chef::CookbookVersion",
		Files:        nil,
		Templates:    nil,
		Attributes: []CookbookItem{{
			Name:        "default.rb",
			Path:        "attributes/default.rb",
			Checksum:    "553637b4fba46b5148f88d6dd3877e2f",
			Specificity: "default",
		}},
		Recipes: []CookbookItem{{
			Name:        "default.rb",
			Path:        "recipes/default.rb",
			Checksum:    "4e15b1e5593d717685323c5dac86b99e",
			Specificity: "default",
		}},
		Definitions: nil,
		Libraries:   nil,
		Providers:   nil,
		Resources:   nil,
		RootFiles: []CookbookItem{{
			Name:        "metadata.rb",
			Path:        "metadata.rb",
			Checksum:    "ba208d0ffc0dd8cbe9c71fb40fb207b2",
			Specificity: "default",
		}},
		Metadata: *cookbook.Metadata.Manifest(),
	}, parsedV0Json)

	var parsedV2Json CookbookManifestV2
	err = json.Unmarshal(v2ManifestJson, &parsedV2Json)
	assert.Nil(t, err)

	assert.Equal(t, CookbookManifestV2{
		CookbookName: "testdeps",
		Name:         "testdeps-0.1.0",
		Version:      "0.1.0",
		ChefType:     "cookbook_version",
		Frozen:       false,
		JsonClass:    "Chef::CookbookVersion",
		AllFiles: []CookbookItem{
			{
				Name:        "attributes/default.rb",
				Path:        "attributes/default.rb",
				Checksum:    "553637b4fba46b5148f88d6dd3877e2f",
				Specificity: "default",
			},
			{
				Name:        "recipes/default.rb",
				Path:        "recipes/default.rb",
				Checksum:    "4e15b1e5593d717685323c5dac86b99e",
				Specificity: "default",
			},
			{
				Name:        "root_files/metadata.rb",
				Path:        "metadata.rb",
				Checksum:    "ba208d0ffc0dd8cbe9c71fb40fb207b2",
				Specificity: "default",
			},
		},
		Metadata: *cookbook.Metadata.Manifest(),
	}, parsedV2Json)
}
