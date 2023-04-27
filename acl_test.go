package chef

import (
	"fmt"
	"net/http"
	"reflect"
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
	if err != nil {
		t.Errorf("ACL.Get returned error: %v", err)
	}

	want := ACL{
		"create": ACLitems{Groups: []string{"clients", "users", "admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
		"read":   ACLitems{Groups: []string{"clients", "users", "admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
		"update": ACLitems{Groups: []string{"users", "admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
		"delete": ACLitems{Groups: []string{"users", "admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
		"grant":  ACLitems{Groups: []string{"admins"}, Actors: []string{"hostname", "pivotal"}, Users: []string{"pivotal"}, Clients: []string{"hostname"}},
	}

	if !reflect.DeepEqual(acl, want) {
		t.Errorf("ACL.Get returned %+v, want %+v", acl, want)
	}
}

func TestACLService_Put(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/nodes/hostname/_acl/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ``)
	})

	acl := NewACL("create", []string{"pivotal"}, []string{"admins"}, []string{"pivotal"}, []string{})
	err := client.ACLs.Put("nodes", "hostname", "create", acl)
	if err != nil {
		t.Errorf("ACL.Put returned error: %v", err)
	}
}
