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
	testGroupJSON = "test/group.json"
)

func TestGroupFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testGroupJSON); err != nil {
		t.Error("Unexpected error '", err, "' during os.Open on", testGroupJSON)
	} else {
		dec := json.NewDecoder(file)
		var g Group
		if err := dec.Decode(&g); err == io.EOF {
			log.Println(g)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

// TODO: Break out these method tests into separate test functions.
func TestGroupsService_Methods(t *testing.T) {
	setup()
	defer teardown()

	// Set up our HTTP routes.
	// FIXME: We should return HTTP response codes as defined by the Chef API so we can test for them.
	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{"group1": "https://url/for/group1", "group2": "https://url/for/group2"}`)
		case r.Method == "POST":
			fmt.Fprintf(w, `{ "uri": "http://localhost:4545/groups/group3" }`)
		}
	})

	mux.HandleFunc("/groups/group3", func(w http.ResponseWriter, r *http.Request) {
		switch {
		// TODO: Add true test for PUT, updating an existing value.
		case r.Method == "GET":
			fmt.Fprintf(w, `{
               		"name": "group3",
                	"groupname": "group3",
                	"orgname": "Test Org, LLC",
                	"actors": ["tester"],
                	"clients": ["tester"],
                	"groups": ["nested-group"]
            		}`)
		case r.Method == "PUT":
			fmt.Fprintf(w, `{
                	"name": "group3",
                	"groupname": "group3",
                	"actors": {
                		"clients": [],
                		"groups": [],
                		"users": ["tester2"]
			}
            		}`)
		case r.Method == "DELETE":
		}
	})

	// Test list
	groups, err := client.Groups.List()
	if err != nil {
		t.Errorf("Groups.List returned error: %v", err)
	}

	listWant := map[string]string{"group1": "https://url/for/group1", "group2": "https://url/for/group2"}

	if !reflect.DeepEqual(groups, listWant) {
		t.Errorf("Groups.List returned %+v, want %+v", groups, listWant)
	}

	// test Get
	group, err := client.Groups.Get("group3")
	if err != nil {
		t.Errorf("Groups.Get returned error: %v", err)
	}

	var wantGroup Group
	wantGroup.Name = "group3"
	wantGroup.GroupName = "group3"
	wantGroup.OrgName = "Test Org, LLC"
	wantGroup.Actors = []string{"tester"}
	wantGroup.Clients = []string{"tester"}
	wantGroup.Groups = []string{"nested-group"}
	if !reflect.DeepEqual(group, wantGroup) {
		t.Errorf("Groups.Get returned %+v, want %+v", group, wantGroup)
	}

	// test Create
	res, err := client.Groups.Create(wantGroup)
	if err != nil {
		t.Errorf("Groups.Create returned error: %s", err.Error())
	}

	createResult := &GroupResult{"http://localhost:4545/groups/group3"}
	if !reflect.DeepEqual(createResult, res) {
		t.Errorf("Groups.Post returned %+v, want %+v", res, createResult)
	}

	// test Update
	groupupdate := GroupUpdate{}
	groupupdate.Name = "group3"
	groupupdate.GroupName = "group3"
	groupupdate.Actors.Clients = []string{}
	groupupdate.Actors.Groups = []string{}
	groupupdate.Actors.Users = []string{"tester2"}
	updateRes, err := client.Groups.Update(groupupdate)
	if err != nil {
		t.Errorf("Groups Update returned error %v", err)
	}

	if !reflect.DeepEqual(updateRes, groupupdate) {
		t.Errorf("Groups Update returned %+v, want %+v", updateRes, groupupdate)
	}

	// test Delete
	err = client.Groups.Delete(group.Name)
	if err != nil {
		t.Errorf("Groups.Delete returned error: %v", err)
	}
}
