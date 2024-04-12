package chef

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCBADownloadThatDoesNotExist(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/cookbook_artifacts/seven_zip/xyz", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	err := client.CookbookArtifacts.DownloadTo("seven_zip", "xyz", "")
	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "404")
	}
}

func TestCBADownloadTo(t *testing.T) {
	setup()
	defer teardown()

	mockedCBAResponseFile := cbaData()
	tempDir, err := os.MkdirTemp("", "seven_zip-cookbook")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tempDir) // clean up

	mux.HandleFunc("/cookbook_artifacts/seven_zip/0e1fed3b56aa5e84205e330d92aca22d8704a014", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(mockedCBAResponseFile))
	})
	mux.HandleFunc("/bookshelf/seven_zip/metadata_rb", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "name 'foo'")
	})
	mux.HandleFunc("/bookshelf/seven_zip/default_rb", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "log 'this is a resource'")
	})

	err = client.CookbookArtifacts.DownloadTo("seven_zip", "0e1fed3b56aa5e84205e330d92aca22d8704a014", tempDir)
	assert.Nil(t, err)

	var (
		cookbookPath = path.Join(tempDir, "seven_zip-0e1fed3b56aa5e84205e")
		metadataPath = path.Join(cookbookPath, "metadata.rb")
		recipesPath  = path.Join(cookbookPath, "recipes")
		defaultPath  = path.Join(recipesPath, "default.rb")
	)
	assert.DirExistsf(t, cookbookPath, "the cookbook directory should exist")
	assert.DirExistsf(t, recipesPath, "the recipes directory should exist")
	if assert.FileExistsf(t, metadataPath, "a metadata.rb file should exist") {
		metadataBytes, err := os.ReadFile(metadataPath)
		assert.Nil(t, err)
		assert.Equal(t, "name 'foo'", string(metadataBytes))
	}
	if assert.FileExistsf(t, defaultPath, "the default.rb recipes should exist") {
		recipeBytes, err := os.ReadFile(defaultPath)
		assert.Nil(t, err)
		assert.Equal(t, "log 'this is a resource'", string(recipeBytes))
	}

}

func cbaData() string {
	return `

{
	"version": "3.1.2",
	"name": "seven_zip",
	"identifier": "0e1fed3b56aa5e84205e330d92aca22d8704a014",
	"frozen?": false,
  	"chef_type": "cookbook_version",
	"root_files": [
      {
		"name": "metadata.rb",
		"path": "metadata.rb",
		"checksum": "6607f3131919e82dc4ba4c026fcfee9f",
		"specificity": "default",
		"url": "` + server.URL + `/bookshelf/seven_zip/metadata_rb"
	  }
	],
  	"attributes": [],
  	"definitions": [],
  	"files": [],
  	"libraries": [],
	"providers": [],
	"recipes": [
      {
      	"name": "default.rb",
      	"path": "recipes/default.rb",
      	"checksum": "8e751ed8663cb9b97499956b6a20b0de",
      	"specificity": "default",
      	"url": "` + server.URL + `/bookshelf/seven_zip/default_rb"
      }
  	],
   "resources": [],
   "templates": [],
   "metadata": {},
   "access": {}
} `
}
