package chef

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

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

	want := &RoleListResult{"foo": "http://localhost:4000/roles/foo", "webserver": "http://localhost:4000/roles/webserver"}

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

	want := &Role{
		Name:               "webserver",
		ChefType:           "role",
		JsonClass:          "Chef::Role",
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

	role := &Role{
		Name:               "webserver",
		ChefType:           "role",
		JsonClass:          "Chef::Role",
		DefaultAttributes:  "",
		Description:        "A webserver",
		RunList:            []string{"recipe[unicorn]", "recipe[apache2]"},
		OverrideAttributes: "",
	}

	uri, err := client.Roles.Create(role)
	if err != nil {
		t.Errorf("Roles.Create returned error: %v", err)
	}

	want := &RoleCreateResult{"uri": "http://localhost:4000/roles/webserver"}

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

	role := &Role{
		Name:        "webserver",
		ChefType:    "role",
		JsonClass:   "Chef::Role",
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
	r := &RoleListResult{"foo": "http://localhost:4000/roles/foo"}
	want := "foo => http://localhost:4000/roles/foo\n"
	if r.String() != want {
		t.Errorf("RoleListResult.String returned %+v, want %+v", r.String(), want)
	}
}

func TestRolesService_RoleCreateResultString(t *testing.T) {
	r := &RoleCreateResult{"uri": "http://localhost:4000/roles/webserver"}
	want := "uri => http://localhost:4000/roles/webserver\n"
	if r.String() != want {
		t.Errorf("RoleCreateResult.String returned %+v, want %+v", r.String(), want)
	}
}
