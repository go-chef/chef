package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cookbookListResponseFile = "test/cookbooks_response.json"
const cookbookResponseFile = "test/cookbook.json"
const _IssueUrl = "https://github.com/<insert_org_here>/apache/issues"
const _Name = "apache"
const _Maintainer = "The Authors"
const _MaintainerEmail = "you@example.com"
const _SourceUrl = "https://github.com/<insert_org_here>/apache"
const _License = "All Rights Reserved"
const _Version = "0.1.0"
const _ChefVersion = ">= 15.0"
const _Description = "Installs/Configures apache"

func TestGetVersion(t *testing.T) {
	setup()
	defer teardown()

	cbookResp, err := ioutil.ReadFile(cookbookResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/cookbooks/foo/_latest", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(cbookResp))
	})

	cookbook, err := client.Cookbooks.GetVersion("foo", "_latest")
	assert.Nil(t, err)
	if assert.NotNil(t, cookbook) {
		assert.Equal(t, "foo", cookbook.CookbookName)
		assert.Equal(t, "0.3.0", cookbook.Version)
		assert.Equal(t, "foo-0.3.0", cookbook.Name)
		assert.Equal(t, "cookbook_version", cookbook.ChefType)
		assert.Equal(t, false, cookbook.Frozen)
		assert.Equal(t, "Chef::CookbookVersion", cookbook.JsonClass)
		assert.Equal(t, 0, len(cookbook.Files))
		assert.Equal(t, 0, len(cookbook.Templates))
		assert.Equal(t, 0, len(cookbook.Attributes))
		assert.Equal(t, 0, len(cookbook.Definitions))
		assert.Equal(t, 0, len(cookbook.Libraries))
		assert.Equal(t, 0, len(cookbook.Providers))
		assert.Equal(t, 0, len(cookbook.Resources))
		// Assert Recipes (verify only one field)
		assert.Equal(t, 1, len(cookbook.Recipes))
		assert.Equal(t, "default.rb", cookbook.Recipes[0].Name)
		assert.Equal(t, "recipes/default.rb", cookbook.Recipes[0].Path)
		assert.Equal(t, "4e855dcab35b481ee56518db164b501d", cookbook.Recipes[0].Checksum)
		assert.Equal(t, "default", cookbook.Recipes[0].Specificity)
		// Check partial string just for convenience
		assert.Contains(t, cookbook.Recipes[0].Url, "https://localhost:443/bookshelf/organization-")
		// Assert RootFiles
		assert.Equal(t, 8, len(cookbook.RootFiles))
		// Assert CookbookMeta struct
		assert.Equal(t, "foo", cookbook.Metadata.Name)
		assert.Equal(t, "0.3.0", cookbook.Metadata.Version)
		assert.Equal(t, "The Authors", cookbook.Metadata.Maintainer)
		assert.Equal(t, "you@example.com", cookbook.Metadata.MaintainerEmail)
		assert.Equal(t, "Installs/Configures foo", cookbook.Metadata.Description)
		assert.Equal(t, "All Rights Reserved", cookbook.Metadata.License)
		// Assert CookbookAccess struct
		assert.Equal(t, true, cookbook.Access.Read)
		assert.Equal(t, true, cookbook.Access.Create)
		assert.Equal(t, true, cookbook.Access.Grant)
		assert.Equal(t, true, cookbook.Access.Update)
		assert.Equal(t, true, cookbook.Access.Delete)
	}
}

func TestCookbookList(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(cookbookListResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/cookbooks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.Cookbooks.List()
	if err != nil {
		t.Error(err)
	}

	if data == nil {
		t.Fatal("WTF we should have some data")
	}
	fmt.Println(data)

	_, err = client.Cookbooks.ListAvailableVersions("3")
	if err != nil {
		t.Error(err)
	}

	_, err = client.Cookbooks.ListAvailableVersions("0")
	if err != nil {
		t.Error(err)
	}
}

func TestCookbookListAvailableVersions_0(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/cookbooks", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "BAD FUCKING REQUEST", 503)
	})

	_, err := client.Cookbooks.ListAvailableVersions("2")
	if err == nil {
		t.Error("We expected this bad request to error", err)
	}
}

func TestCookBookDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/cookbooks/good/1.1.1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "")
	})
	mux.HandleFunc("/cookbooks/bad/1.1.1", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	err := client.Cookbooks.Delete("bad", "1.1.1")
	if err == nil {
		t.Error("We expected this bad request to error", err)
	}

	err = client.Cookbooks.Delete("good", "1.1.1")
	if err != nil {
		t.Error(err)
	}
}

func TestCookBookGet(t *testing.T) {
	setup()
	defer teardown()

	cookbookVerionJSON := `{"url": "http://localhost:4000/cookbooks/apache2/5.1.0", "version": "5.1.0"}`
	mux.HandleFunc("/cookbooks/good", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, cookbookVerionJSON)
	})
	mux.HandleFunc("/cookbooks/bad", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	data, err := client.Cookbooks.Get("good")
	if err != nil {
		t.Error(err)
	}

	if data.Version != "5.1.0" {
		t.Errorf("We expected '5.1.0' and got '%s'\n", data.Version)
	}

	_, err = client.Cookbooks.Get("bad")
	if err == nil {
		t.Error("We expected this bad request to error", err)
	}
}

func TestCookBookGetAvailableVersions(t *testing.T) {
	setup()
	defer teardown()

	cookbookVerionsJSON := `
	{	"apache2": {
    "url": "http://localhost:4000/cookbooks/apache2",
    "versions": [
      {"url": "http://localhost:4000/cookbooks/apache2/5.1.0",
       "version": "5.1.0"},
      {"url": "http://localhost:4000/cookbooks/apache2/4.2.0",
       "version": "4.2.0"}
    ]
	}}`

	mux.HandleFunc("/cookbooks/good", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, cookbookVerionsJSON)
	})
	mux.HandleFunc("/cookbooks/bad", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	data, err := client.Cookbooks.GetAvailableVersions("good", "3")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}

func TestCookBookListAllRecipes(t *testing.T) {
	setup()
	defer teardown()

	cookbookRecipesJSON := `
	[
	  "apache2",
	  "apache2::mod_access_compat",
	  "apache2::mod_actions",
	  "apache2::mod_alias"
	]`

	mux.HandleFunc("/cookbooks/_recipes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, cookbookRecipesJSON)
	})

	data, err := client.Cookbooks.ListAllRecipes()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)
}

func TestNewCookbookMeta(t *testing.T) {
	data := `   name 'apache'
				maintainer 'The Authors'
				maintainer_email 'you@example.com'
				license 'All Rights Reserved'
				description 'Installs/Configures apache'
				version '0.1.0'
				chef_version '>= 15.0'
				
				#issues_url points to the location where issues for this cookbook are
				# tracked.  A View Issues link will be displayed on this cookbook's page when
				# uploaded to a Supermarket.
				#
				issues_url 'https://github.com/<insert_org_here>/apache/issues'
				
				source_url 'https://github.com/<insert_org_here>/apache'`
	md, err := NewMetaData(data)
	if err != nil {
		t.Error("invalid metadata.rb contain please validate it", err)
	}
	validateCookbookMetaData(md, t, "TestNewMetaData")

}

func TestNewCookbookMetaFromJson(t *testing.T) {
	data := `{"name":"apache","description":"Installs/Configures apache","long_description":"","maintainer":"The Authors","maintainer_email":"you@example.com","license":"All Rights Reserved","platforms":{},"dependencies":{},"providing":null,"recipes":null,"version":"0.1.0","source_url":"https://github.com/\u003cinsert_org_here\u003e/apache","issues_url":"https://github.com/\u003cinsert_org_here\u003e/apache/issues","ChefVersion":"\u003e= 15.0","OhaiVersion":"","gems":null,"eager_load_libraries":false,"privacy":false}`
	md, err := NewMetaDataFromJson([]byte(data))
	if err != nil {
		t.Error("invalid metadata.rb contain please validate it", err)
	}
	validateCookbookMetaData(md, t, "TestNewMetaDataFromJson")
}
func TestReadCookbookMeta(t *testing.T) {
	file, err := os.Create("/tmp/metadata.rb")
	if err != nil {
		t.Error("unable to create to metadata.rb", err)
	}
	defer file.Close()
	data := `   name 'apache'
				maintainer 'The Authors'
				maintainer_email 'you@example.com'
				license 'All Rights Reserved'
				description 'Installs/Configures apache'
				version '0.1.0'
				chef_version '>= 15.0'
				
				#issues_url points to the location where issues for this cookbook are
				# tracked.  A View Issues link will be displayed on this cookbook's page when
				# uploaded to a Supermarket.
				#
				issues_url 'https://github.com/<insert_org_here>/apache/issues'
				
				source_url 'https://github.com/<insert_org_here>/apache'`
	_, err = file.WriteString(data)
	if err != nil {
		t.Error("error in creating tmp file for metadata.rb", err)
	}
	md, err := ReadMetaData("/tmp")
	if err != nil {
		t.Error("error in reading tmp file for metadata.rb", err)
	}
	validateCookbookMetaData(md, t, "TestReadMetaData")
	os.Remove("/tmp/metadata.rb")

}
func TestReadCookbookMeta2(t *testing.T) {
	file, err := os.Create("/tmp/metadata.json")
	if err != nil {
		t.Error("unable to create to metadata.rb", err)
	}
	defer file.Close()
	data := `{"name":"apache","description":"Installs/Configures apache","long_description":"","maintainer":"The Authors","maintainer_email":"you@example.com","license":"All Rights Reserved","platforms":{},"dependencies":{},"providing":null,"recipes":null,"version":"0.1.0","source_url":"https://github.com/\u003cinsert_org_here\u003e/apache","issues_url":"https://github.com/\u003cinsert_org_here\u003e/apache/issues","ChefVersion":"\u003e= 15.0","OhaiVersion":"","gems":null,"eager_load_libraries":false,"privacy":false}`
	_, err = file.WriteString(data)
	if err != nil {
		t.Error("error in creating tmp file for metadata.json", err)
	}
	md, err := ReadMetaData("/tmp")
	if err != nil {
		t.Error("error in reading tmp file for metadata.json", err)
	}
	validateCookbookMetaData(md, t, "TestReadMetaData")
	os.Remove("/tmp/metadata.json")
}
func validateCookbookMetaData(md CookbookMeta, t *testing.T, funcName string) {
	assert.Equal(t, _Description, md.Description)
	assert.Equal(t, _IssueUrl, md.IssueUrl)
	assert.Equal(t, _Name, md.Name)
	assert.Equal(t, _Maintainer, md.Maintainer)
	assert.Equal(t, _MaintainerEmail, md.MaintainerEmail)
	assert.Equal(t, _SourceUrl, md.SourceUrl)
	assert.Equal(t, _License, md.License)
	assert.Equal(t, _Version, md.Version)
	assert.Equal(t, _ChefVersion, md.ChefVersion)

}
