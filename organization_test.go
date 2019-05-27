package chef

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
)

var (
	testOrganizationJSON = "test/organization.json"
)

func TestOrganizationFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testOrganizationJSON); err != nil {
		t.Error("Unexpected error '", err, "' during os.Open on", testOrganizationJSON)
	} else {
		dec := json.NewDecoder(file)
		var g Organization
		if err := dec.Decode(&g); err == io.EOF {
			log.Fatal(g)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func TestOrganizationslist(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{ "org_name1": "https://url/for/org_name1", "org_name2": "https://url/for/org_name2"}`)
		case r.Method == "POST":
			fmt.Fprintf(w, `{
				 "name": "organization1"
				"full_name": "This Organization"
			 }`)
		}
	})

	// Test list
	organizations, err := client.Organizations.List()
	if err != nil {
		t.Errorf("Organizations.List returned error: %v", err)
	}
	listWant := map[string]string{"org_name1": "https://url/for/org_name1", "org_name2": "https://url/for/org_name2"}
	if !reflect.DeepEqual(organizations, listWant) {
		t.Errorf("Organizations.List returned %+v, want %+v", organizations, listWant)
	}
}

func TestOrganizationsCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{ "org_name1": "https://url/for/org_name1", "org_name2": "https://url/for/org_name2"}`)
		case r.Method == "POST":
			fmt.Fprintf(w, `{ "clientname": "organization3", "private_key": "mine", "uri": "https://url/for/organization3" }`)
		}
	})
	var wantOrganization Organization
	wantOrganization.Name = "organization3"
	wantOrganization.FullName = "Chef Software"

	res, err := client.Organizations.Create(wantOrganization)
	if err != nil {
		t.Errorf("Organizations.Create returned error: %s", err.Error())
	}
	createResult := OrganizationResult{ClientName: "organization3", PrivateKey: "mine", Uri: "https://url/for/organization3"}
	if !reflect.DeepEqual(createResult, res) {
		t.Errorf("Organizations.Post returned %+v, want %+v", res, createResult)
	}
}

func TestOrganizationsUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations/organization3", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{
                                "name": "organization3",
                                "full_name": "Chef Software",
                                "guid": "f980d1asdfda0331235s00ff36862"
                        }`)
		case r.Method == "PUT":
			fmt.Fprintf(w, `{
                                "name": "organization3",
                                "full_name": "Updated Software",
                                "guid": "f980d1asdfda0331235s00ff36862"
                        }`)
		}
	})
	// test Update
	organization, err := client.Organizations.Get("organization3")
	organization.FullName = "Updated Software"
	updateRes, err := client.Organizations.Update(organization)
	if err != nil {
		t.Errorf("Organizations.Update returned error: %v", err)
	}
	if !reflect.DeepEqual(updateRes, organization) {
		t.Errorf("Organizations.Update returned %+v, want %+v", updateRes, organization)
	}
}

func TestOrganizationsDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/organizations/organization3", func(w http.ResponseWriter, r *http.Request) {
		switch {
		// TODO: Add true test for PUT, updating an existing value.
		case r.Method == "GET" || r.Method == "DELETE":
			fmt.Fprintf(w, `{
                                "name": "organization3",
                                "full_name": "Chef Software",
                                "guid": "f980d1asdfda0331235s00ff36862"
                        }`)
		}
	})

	// test Delete
	err := client.Organizations.Delete("organization3")
	if err != nil {
		t.Errorf("Organizations.Delete returned error: %v", err)
	}
}

func TestOrganizationsGet(t *testing.T) {
	// TODO: We should return HTTP response codes as defined by the Chef API so we can test for them.
	setup()
	defer teardown()

	mux.HandleFunc("/organizations/organization3", func(w http.ResponseWriter, r *http.Request) {
		switch {
		// TODO: Add true test for PUT, updating an existing value.
		case r.Method == "GET" || r.Method == "PUT" || r.Method == "DELETE":
			fmt.Fprintf(w, `{
                		"name": "organization3",
                		"full_name": "Chef Software",
                		"guid": "f980d1asdfda0331235s00ff36862"
            		}`)
		}
	})

	// test Get
	organization, err := client.Organizations.Get("organization3")
	if err != nil {
		t.Errorf("Organizations.Get returned error: %v", err)
	}

	var wantOrganization Organization
	wantOrganization.Name = "organization3"
	wantOrganization.FullName = "Chef Software"
	wantOrganization.Guid = "f980d1asdfda0331235s00ff36862"
	if !reflect.DeepEqual(organization, wantOrganization) {
		t.Errorf("Organizations.Get returned %+v, want %+v", organization, wantOrganization)
	}
}
