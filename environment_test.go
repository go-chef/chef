package chef_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/go-chef/chef"
)

var (
	testEnvironmentJSON = "test/environment.json"
)

// BUG(fujin): re-do with goconvey
func TestEnvironmentFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testEnvironmentJSON); err != nil {
		t.Error("unexpected error", err, "during os.Open on", testEnvironmentJSON)
	} else {
		dec := json.NewDecoder(file)
		var e chef.Environment
		if err := dec.Decode(&e); err == io.EOF {
			log.Println(e)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func TestEnvironmentsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/environments", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"_default":"blah", "development":"blah"}`)
	})

	environments, err := client.Environments.List()
	if err != nil {
		t.Errorf("Environments.List returned error: %v", err)
	}

	want := &chef.EnvironmentListResult{"_default": "blah", "development": "blah"}

	if !reflect.DeepEqual(environments, want) {
		t.Errorf("Environments.List returned %+v, want %+v", environments, want)
	}
}

func TestEnvironmentsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/environments/testenvironment", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
	    "name": "testenvironment",
	    "json_class": "Chef::Environment",
	    "chef_type": "environment"
		}`)
	})

	environments, err := client.Environments.Get("testenvironment")
	if err != nil {
		t.Errorf("Environments.Get returned error: %v", err)
	}

	want := &chef.Environment{
		Name:      "testenvironment",
		JsonClass: "Chef::Environment",
		ChefType:  "environment",
	}

	if !reflect.DeepEqual(environments, want) {
		t.Errorf("Environments.Get returned %+v, want %+v", environments, want)
	}
}

func TestEnvironmentsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/environments", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "uri": "http://localhost:4000/environments/dev" }`)
	})

	role := &chef.Environment{
		Name:             "dev",
		ChefType:         "environment",
		JsonClass:        "Chef::Environment",
		Attributes:       "",
		Description:      "",
		CookbookVersions: map[string]string{},
	}

	uri, err := client.Environments.Create(role)
	if err != nil {
		t.Errorf("Environments.Create returned error: %v", err)
	}

	want := &chef.EnvironmentCreateResult{"uri": "http://localhost:4000/environments/dev"}

	if !reflect.DeepEqual(uri, want) {
		t.Errorf("Environments.Create returned %+v, want %+v", uri, want)
	}
}

func TestEnvironmentsService_Put(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/environments/dev", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
		  "name": "dev",
		  "json_class": "Chef::Environment",
		  "description": "The Dev Environment",
		  "cookbook_versions": {},
		  "chef_type": "environment"
		}`)
	})

	environment := &chef.Environment{
		Name:             "dev",
		ChefType:         "environment",
		JsonClass:        "Chef::Environment",
		Description:      "The Dev Environment",
		CookbookVersions: map[string]string{},
	}

	updatedEnvironment, err := client.Environments.Put(environment)
	if err != nil {
		t.Errorf("Environments.Put returned error: %v", err)
	}

	if !reflect.DeepEqual(updatedEnvironment, environment) {
		t.Errorf("Environments.Put returned %+v, want %+v", updatedEnvironment, environment)
	}
}

func TestEnvironmentsService_EnvironmentListResultString(t *testing.T) {
	e := &chef.EnvironmentListResult{
		"_default":  "https://api.opscode.com/organizations/org_name/environments/_default",
		"webserver": "https://api.opscode.com/organizations/org_name/environments/webserver"}
	want := "_default => https://api.opscode.com/organizations/org_name/environments/_default\nwebserver => https://api.opscode.com/organizations/org_name/environments/webserver\n"
	if e.String() != want {
		t.Errorf("EnvironmentListResult.String returned %+v, want %+v", e.String(), want)
	}
}

func TestEnvironmentsService_EnvironmentCreateResultString(t *testing.T) {
	e := &chef.EnvironmentCreateResult{"uri": "http://localhost:4000/environments/dev"}
	want := "uri => http://localhost:4000/environments/dev\n"
	if e.String() != want {
		t.Errorf("EnvironmentCreateResult.String returned %+v, want %+v", e.String(), want)
	}
}
