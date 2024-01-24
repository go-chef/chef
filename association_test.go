package chef

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestAssociationMethods(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/association_requests", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `[{"id":"1f", "username":"jollystranger"}, {"id":"2b", "username":"fredhamlet"}]`)
		case r.Method == "POST":
			fmt.Fprintf(w, `{ "uri": "http://chef/organizations/test/association_requests/1a" }`)
		}
	})

	mux.HandleFunc("/association_requests/1f", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `[
			{ "user": {"username": "jollystranger"}
		]`)
		case r.Method == "DELETE":
			fmt.Fprintf(w, `{
	    "id": "1f",
	    "orgname": "test",
	    "username": "jollystranger"
		}`)
		}
	})

	// test Invite - Invite a user
	request := Request{
		User: "jollystranger",
	}
	invite, err := client.Associations.Invite(request)
	if err != nil {
		t.Errorf("Associations.Invite returned error: %v", err)
	}
	associationWant := Association{
		Uri: "http://chef/organizations/test/association_requests/1a",
	}
	if !reflect.DeepEqual(invite, associationWant) {
		t.Errorf("Associations.Invite returned %+v, want %+v", invite, associationWant)
	}

	// test ListInvites - return the existing invitations
	invites, err := client.Associations.ListInvites()
	if err != nil {
		t.Errorf("Associations.List returned error: %v", err)
	}
	listWant := []Invite{
		{Id: "1f", UserName: "jollystranger"},
		{Id: "2b", UserName: "fredhamlet"},
	}
	if !reflect.DeepEqual(invites, listWant) {
		t.Errorf("Associations.InviteList returned %+v, want %+v", invites, listWant)
	}

	// test DeleteInvite
	deleted, err := client.Associations.DeleteInvite("1f")
	if err != nil {
		t.Errorf("Associations.Get returned error: %v", err)
	}
	wantInvite := RescindInvite{
		Id:       "1f",
		Orgname:  "test",
		Username: "jollystranger",
	}
	if !reflect.DeepEqual(deleted, wantInvite) {
		t.Errorf("Associations.RescindInvite returned %+v, want %+v", deleted, wantInvite)
	}

	// test InviteId
	id, err := client.Associations.InviteId("fredhamlet")
	if err != nil {
		t.Errorf("Associations.InviteId returned error: %s", err.Error())
	}
	wantId := "2b"
	if !reflect.DeepEqual(id, wantId) {
		t.Errorf("Associations.InviteId returned %+v, want %+v", id, wantId)
	}
}

func TestOrgUserMethods(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			// []OrgUserListEntry
			fmt.Fprintf(w, `[{ "user": {"username": "jollystranger"}}]`)
		case r.Method == "POST":
			// err only
		}
	})

	mux.HandleFunc("/users/jollystranger", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			// OrgUser
			fmt.Fprintf(w, `{ "username": "jollystranger", "email": "jolly.stranger@domain.io", "display_name":"jolly" }`)
		case r.Method == "DELETE":
			// OrgUser
			fmt.Fprintf(w, `{ "username": "jollystranger", "email": "jolly.stranger@domain.io", "display_name":"jolly" }`)
		}
	})

	// test List the users
	users, err := client.Associations.List()
	if err != nil {
		t.Errorf("Associations.List returned error: %v", err)
	}
	var a struct {
		Username string `json:"username,omitempty"`
	}
	a.Username = "jollystranger"
	var wuser OrgUserListEntry
	wuser.User = a
	wantUsers := []OrgUserListEntry{
		wuser,
	}
	if !reflect.DeepEqual(users, wantUsers) {
		t.Errorf("Associations.List returned %+v, want %+v", users, wantUsers)
	}

	// test Add user
	addme := AddNow{
		Username: "jollystranger",
	}
	err = client.Associations.Add(addme)
	if err != nil {
		t.Errorf("Associations.Add returned error: %v", err)
	}

	// test Get user details
	user, err := client.Associations.Get("jollystranger")
	if err != nil {
		t.Errorf("Associations.Get returned error: %v", err)
	}
	wantUser := OrgUser{
		Username: "jollystranger", Email: "jolly.stranger@domain.io", DisplayName: "jolly",
	}
	if !reflect.DeepEqual(user, wantUser) {
		t.Errorf("Associations.Get returned %+v, want %+v", user, wantUser)
	}

	// test Delete user details
	delu, err := client.Associations.Delete("jollystranger")
	if err != nil {
		t.Errorf("Associations.Delete returned error: %v", err)
	}
	delUser := OrgUser{
		Username: "jollystranger", Email: "jolly.stranger@domain.io", DisplayName: "jolly",
	}
	if !reflect.DeepEqual(delu, delUser) {
		t.Errorf("Associations.Delete returned %+v, want %+v", delu, delUser)
	}
}
