package chef

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
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
		var e Environment
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

	want := &EnvironmentListResult{"_default": "blah", "development": "blah"}

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

	want := &Environment{
		Name:      "testenvironment",
		JsonClass: "Chef::Environment",
		ChefType:  "environment",
	}

	if !reflect.DeepEqual(environments, want) {
		t.Errorf("Environments.Get returned %+v, want %+v", environments, want)
	}
}
