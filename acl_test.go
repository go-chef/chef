package chef

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestACLService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/nodes/hostname/_acl", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
      "create": {
        "actors": [
          "hostname",
          "pivotal"
        ],
        "groups": [
          "clients",
          "users",
          "admins"
        ],
        "users": [
          "pivotal"
        ],
        "clients": [
          "hostname"
        ]
      },
      "read": {
        "actors": [
          "hostname",
          "pivotal"
        ],
        "groups": [
          "clients",
          "users",
          "admins"
        ],
        "users": [
          "pivotal"
        ],
        "clients": [
          "hostname"
        ]
      },
      "update": {
        "actors": [
          "hostname",
          "pivotal"
        ],
        "groups": [
          "users",
          "admins"
        ],
        "users": [
          "pivotal"
        ],
        "clients": [
          "hostname"
        ]
      },
      "delete": {
        "actors": [
          "hostname",
          "pivotal"
        ],
        "groups": [
          "users",
          "admins"
        ],
        "users": [
          "pivotal"
        ],
        "clients": [
          "hostname"
        ]
      },
      "grant": {
        "actors": [
          "hostname",
          "pivotal"
        ],
        "groups": [
          "admins"
        ],
        "users": [
          "pivotal"
        ],
        "clients": [
          "hostname"
        ]
      }
    }
    `)
	})

	acl, err := client.ACLs.Get("nodes", "hostname")
	assert.Nil(t, err, "Get returned error")

	want := ACL{
		"create": ACLitems{Groups: []string{"clients", "users", "admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
		"read":   ACLitems{Groups: []string{"clients", "users", "admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
		"update": ACLitems{Groups: []string{"users", "admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
		"delete": ACLitems{Groups: []string{"users", "admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
		"grant":  ACLitems{Groups: []string{"admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
	}

	assert.Equal(t, want, acl, "Get Return")

}

func TestNewACL(t *testing.T) {

}

func TestACLAdminAccess(t *testing.T) {
	acl := NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{"pivotal", "user"}, []string{"client"})
	err := ACLAdminAccess(acl)
	assert.Nil(t, err, fmt.Sprintf("Pivotal missing %+v\n", acl))

	acl = NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{"pivotal", "folks", "other"}, []string{})
	err = ACLAdminAccess(acl)
	assert.Nil(t, err, fmt.Sprintf("Pivotal first %+v\n", acl))

	acl = NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{"other", "pivotal", "folks"}, []string{})
	err = ACLAdminAccess(acl)
	assert.Nil(t, err, fmt.Sprintf("Pivotal first %+v\n", acl))

	acl = NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{"other", "folks", "pivotal"}, []string{})
	err = ACLAdminAccess(acl)
	assert.Nil(t, err, fmt.Sprintf("Pivotal last %+v\n", acl))

	acl = NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{}, []string{})
	err = ACLAdminAccess(acl)
	assert.NotNil(t, err, fmt.Sprintf("Pivotal not there %+v\n", acl))

	acl = NewACL("create", []string{"pivotal"}, []string{"admins"}, nil, []string{})
	err = ACLAdminAccess(acl)
	assert.NotNil(t, err, fmt.Sprintf("Nil user array %+v\n", acl))

	myacl := *NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{"pivotal"}, []string{})
	myacl["read"] = *NewACLItems([]string{"pivotal"}, []string{"admins"}, []string{"pivotal"}, []string{})
	myacl["destroy"] = *NewACLItems([]string{"pivotal"}, []string{"admins"}, []string{"pivotal"}, []string{})
	err = ACLAdminAccess(&myacl)
	assert.Nil(t, err, fmt.Sprintf("mutliple types ok %+v\n", myacl))

	myacl["read"] = *NewACLItems([]string{"pivotal"}, []string{"admins"}, []string{}, []string{})
	err = ACLAdminAccess(&myacl)
	assert.NotNil(t, err, fmt.Sprintf("mutliple types missing pivotal %+v\n", myacl))

}

func TestACLService_Put(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/nodes/hostname/_acl/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ``)
	})

	acl := NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{"pivotal"}, []string{})
	err := client.ACLs.Put("nodes", "hostname", "create", acl)
	assert.Nil(t, err, "Put returned error")

	acl = NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{}, []string{})
	err = client.ACLs.Put("nodes", "hostname", "create", acl)
	assert.NotNil(t, err, "Put should return error, pivotal not in users list")
}
