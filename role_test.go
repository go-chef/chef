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
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testRoleJSON = "test/role.json"
	// FML
	testRole = &chef.Role{
		Name:               "test",
		ChefType:           "role",
		Description:        "Test Role",
		RunList:            []string{"recipe[foo]", "recipe[baz]", "role[banana]"},
		JSONClass:          "Chef::Role",
		DefaultAttributes:  struct{}{},
		OverrideAttributes: struct{}{},
	}
)

func TestRoleName(t *testing.T) {
	// BUG(spheromak): Pull these constructors out into a Convey Decorator
	n1 := testRole
	name := n1.Name

	Convey("Role name is 'test'", t, func() {
		So(name, ShouldEqual, "test")
	})
}

// BUG(fujin): re-do with goconvey
func TestRoleFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testRoleJSON); err != nil {
		t.Error("unexpected error", err, "during os.Open on", testRoleJSON)
	} else {
		dec := json.NewDecoder(file)
		var n chef.Role
		if err := dec.Decode(&n); err == io.EOF {
			log.Println(n)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func TestRolesService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"foo":"http://localhost:4000/roles/foo", "webserver":"http://localhost:4000/roles/webserver"}`)
	})

	roles, err := client.Roles.List()
	if err != nil {
		t.Errorf("Roles.List returned error: %v", err)
	}

	want := &chef.RoleListResult{"foo": "http://localhost:4000/roles/foo", "webserver": "http://localhost:4000/roles/webserver"}

	if !reflect.DeepEqual(roles, want) {
		t.Errorf("Roles.List returned %+v, want %+v", roles, want)
	}
}

func TestRolesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/roles/webserver", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
		  "name": "webserver",
		  "chef_type": "role",
		  "json_class": "Chef::Role",
		  "default_attributes": "",
		  "description": "A webserver",
		  "run_list": [
		    "recipe[unicorn]",
		    "recipe[apache2]"
		  ],
		  "override_attributes": ""
		}
		`)
	})

	role, err := client.Roles.Get("webserver")
	if err != nil {
		t.Errorf("Roles.Get returned error: %v", err)
	}

	want := &chef.Role{
		Name:               "webserver",
		ChefType:           "role",
		JSONClass:          "Chef::Role",
		DefaultAttributes:  "",
		Description:        "A webserver",
		RunList:            []string{"recipe[unicorn]", "recipe[apache2]"},
		OverrideAttributes: "",
	}

	if !reflect.DeepEqual(role, want) {
		t.Errorf("Roles.Get returned %+v, want %+v", role, want)
	}
}

func TestRolesService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "uri": "http://localhost:4000/roles/webserver" }`)
	})

	role := &chef.Role{
		Name:               "webserver",
		ChefType:           "role",
		JSONClass:          "Chef::Role",
		DefaultAttributes:  "",
		Description:        "A webserver",
		RunList:            []string{"recipe[unicorn]", "recipe[apache2]"},
		OverrideAttributes: "",
	}

	uri, err := client.Roles.Create(role)
	if err != nil {
		t.Errorf("Roles.Create returned error: %v", err)
	}

	want := &chef.RoleCreateResult{"uri": "http://localhost:4000/roles/webserver"}

	if !reflect.DeepEqual(uri, want) {
		t.Errorf("Roles.Create returned %+v, want %+v", uri, want)
	}
}

func TestRolesService_Put(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/roles/webserver", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
		  "name": "webserver",
		  "chef_type": "role",
		  "json_class": "Chef::Role",
		  "description": "A webserver",
		  "run_list": [
		    "recipe[apache2]"
		  ]
		}`)
	})

	role := &chef.Role{
		Name:        "webserver",
		ChefType:    "role",
		JSONClass:   "Chef::Role",
		Description: "A webserver",
		RunList:     []string{"recipe[apache2]"},
	}

	updatedRole, err := client.Roles.Put(role)
	if err != nil {
		t.Errorf("Roles.Put returned error: %v", err)
	}

	if !reflect.DeepEqual(updatedRole, role) {
		t.Errorf("Roles.Put returned %+v, want %+v", updatedRole, role)
	}
}

func TestRolesService_RoleListResultString(t *testing.T) {
	r := &chef.RoleListResult{"foo": "http://localhost:4000/roles/foo"}
	want := "foo => http://localhost:4000/roles/foo\n"
	if r.String() != want {
		t.Errorf("RoleListResult.String returned %+v, want %+v", r.String(), want)
	}
}

func TestRolesService_RoleCreateResultString(t *testing.T) {
	r := &chef.RoleCreateResult{"uri": "http://localhost:4000/roles/webserver"}
	want := "uri => http://localhost:4000/roles/webserver\n"
	if r.String() != want {
		t.Errorf("RoleCreateResult.String returned %+v, want %+v", r.String(), want)
	}
}
