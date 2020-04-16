package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const CBAListResponseFile = "test/cookbook_artifacts_list_response.json"
const CBAGetResponseFile = "test/cookbook_artifacts_get_response.json"
const CBAGetVersionResponseFile = "test/cookbook_artifacts_getversion_response.json"

func TestListCBA(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(CBAListResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/cookbook_artifacts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.CookbookArtifacts.List()
	if err != nil {
		t.Error(err)
	}

	if data == nil {
		t.Fatal("We should have some data")
	}

	if len(data) != 2 {
		t.Error("Mismatch in expected policies count. Expected 2, Got: ", len(data))
	}

	if _, ok := data["oc-hec-postfix"]; !ok {
		t.Error("oc-hec-postfix policy should be listed")
	}

	if _, ok := data["grafana"]; !ok {
		t.Error("grafana policy should be listed")
	}

}

func TestGetCBA(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(CBAGetResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/cookbook_artifacts/seven_zip", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})
	mux.HandleFunc("/cookbook_artifacts/bad", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	data, err := client.CookbookArtifacts.Get("seven_zip")
	if err != nil {
		t.Error(err)
	}

	if data["seven_zip"].Url != "https://localhost/organizations/utility/cookbook_artifacts/seven_zip" {
		t.Errorf("URL mismatch, expected '%s', got '%s'\n", "https://api.chef.io/organizations/chef-utility/cookbook_artifacts/seven_zip", data["seven_zip"].Url)
	}

	if len(data["seven_zip"].CBAVersions) != 7 {
		t.Errorf("Expected 7 versions of this cookbook artifact, received %d", len(data["seven_zip"].CBAVersions))
	}
	_, err = client.CookbookArtifacts.Get("bad")
	if err == nil {
		t.Error("We expected this bad request to error", err)
	}
}

func TestGetVersionCBA(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(CBAGetVersionResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/cookbook_artifacts/seven_zip/0e1fed3b56aa5e84205e330d92aca22d8704a014", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})
	mux.HandleFunc("/cookbook_artifacts/seven_zip/bad", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	cookbookArtifact, err := client.CookbookArtifacts.GetVersion("seven_zip", "0e1fed3b56aa5e84205e330d92aca22d8704a014")
	if err != nil {
		t.Error(err)
	}

	if assert.NotNil(t, cookbookArtifact) {
		assert.Equal(t, "seven_zip", cookbookArtifact.Name)
		assert.Equal(t, "3.1.2", cookbookArtifact.Version)
		assert.Equal(t, "0e1fed3b56aa5e84205e330d92aca22d8704a014", cookbookArtifact.Identifier)
		assert.Equal(t, 9, len(cookbookArtifact.RootFiles))
		assert.Equal(t, 1, len(cookbookArtifact.Providers))
		assert.Equal(t, 2, len(cookbookArtifact.Resources))
		assert.Equal(t, 1, len(cookbookArtifact.Libraries))
		assert.Equal(t, 1, len(cookbookArtifact.Attributes))
		assert.Equal(t, "https://s3-external-1/amazonaws.com:443/test-s3-url", cookbookArtifact.RootFiles[0].Url)
		assert.Equal(t, "d07d9934f97445c5187cf86abf107b7b", cookbookArtifact.Providers[0].Checksum)
		assert.Equal(t, "archive.rb", cookbookArtifact.Resources[0].Name)
		assert.Equal(t, "matchers.rb", cookbookArtifact.Libraries[0].Name)
		assert.Equal(t, "https://s3-external-1/amazonaws.com:443/test-s3-url-attribs", cookbookArtifact.Attributes[0].Url)
		assert.Equal(t, "https://s3-external-1/amazonaws.com:443/test-s3-url-recipes", cookbookArtifact.Recipes[0].Url)
		assert.Equal(t, "seven_zip", cookbookArtifact.Metadata.Name)
		assert.Equal(t, "3.1.2", cookbookArtifact.Metadata.Version)
		assert.Equal(t, "Apache-2.0", cookbookArtifact.Metadata.License)
		assert.Equal(t, "Installs/Configures 7-Zip", cookbookArtifact.Metadata.Description)
		assert.Equal(t, ">= 0.0.0", cookbookArtifact.Metadata.Depends["windows"])
		assert.Equal(t, ">= 13.0", cookbookArtifact.Metadata.ChefVersions[0][0])
	}

	_, err = client.CookbookArtifacts.GetVersion("seven_zip", "bad")
	if err == nil {
		t.Error("We expected this bad request to error", err)
	}
}
