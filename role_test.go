package chef

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
)

var (
	testRoleJSON = "test/role.json"
	// FML
	testRoleMapStringInterfaceLol, _ = NewRole(&Reader{
		"name":       "test",
		"run_list":   []string{"recipe[foo]", "recipe[baz]", "role[banana]"},
		"chef_type":  "role",
		"json_class": "Chef::Role",
		"default_attributes": map[string]interface{}{
			"tags": map[string]interface{}{},
			"openssh": map[string]interface{}{
				"server": map[string]string{
					"permit_root_login": "no",
					"max_auth_tries":    "3",
				},
			},
		},
		"override_attributes": map[string]interface{}{
			"openssh": map[string]interface{}{
				"server": map[string]string{
					"permit_root_login": "yes",
					"max_auth_tries":    "1",
				},
			},
		},
	})
)

func TestRoleName(t *testing.T) {
	// BUG(spheromak): Pull these constructors out into a Convey Decorator
	n1 := testRoleMapStringInterfaceLol // (*Role)
	name := n1.Name

	Convey("Role name is 'test'", t, func() {
		So(name, ShouldEqual, "test")
	})

	swordWithoutASheathe, err := NewRole(&Reader{
		"foobar": "baz",
	})

	name = swordWithoutASheathe.Name
	Convey("Role without a name", t, func() {
		So(name, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestRoleAttribute(t *testing.T) {
	n := testRoleMapStringInterfaceLol
	attr := n.Default
	// BUG(spheromak): Holy shit this is ugly. Need to do something to make this easier for sure.
	ugh := attr["openssh"].(map[string]interface{})["server"].(map[string]string)["permit_root_login"]
	Convey("Role.Default should map", t, func() {
		So(ugh, ShouldEqual, "no")
	})
}

// BUG(fujin): re-do with goconvey
func TestRoleFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testRoleJSON); err != nil {
		t.Error("unexpected error", err, "during os.Open on", testRoleJSON)
	} else {
		dec := json.NewDecoder(file)
		var n Role
		if err := dec.Decode(&n); err == io.EOF {
			log.Println(n)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

// TestNewRole checks the NewRole Reader chain for Type
func TestNewRole(t *testing.T) {
	var v interface{}
	v = testRoleMapStringInterfaceLol
	Convey("NewRole should create a Role", t, func() {
		So(v, ShouldHaveSameTypeAs, &Role{})
	})

	Convey("NewRole should error if decode fails", t, func() {

		failRole, err := NewRole(&Reader{
			"name": struct{}{},
		})

		So(err, ShouldNotBeNil)
		So(failRole, ShouldBeNil)
	})
}

// TestRoleReadIntoFile tests that Read() can be used to read by io.Readers
// BUG(fujin): re-do with goconvey
func TestRoleReadIntoFile(t *testing.T) {
	n1 := testRoleMapStringInterfaceLol // (*Role)
	tf, _ := ioutil.TempFile("test", "role-to-file")
	// Copy to tempfile (I use Read() internally)
	// BUG(fujin): this is currently doing that weird 32768 bytes read thing again.
	io.Copy(tf, n1)

	// Close and remove tempfile
	tf.Close()
	os.Remove(path.Clean(tf.Name()))
}
