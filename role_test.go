package chef_test

import (
	"encoding/json"
	"github.com/go-chef/chef"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"log"
	"os"
	"testing"
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
