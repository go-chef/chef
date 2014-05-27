package chef

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	// "io/ioutil"
	"log"
	"os"
	// "path"
	"testing"
)

var (
	testEnvironmentJSON = "test/environment.json"
	// FML
	testEnvironmentMapStringInterfaceLol, _ = NewEnvironment(&Reader{
		"name": "testenvironment",
		// "run_list":   []string{"recipe[foo]", "recipe[baz]", "role[banana]"},
		"chef_type":  "environment",
		"json_class": "Chef::Environment",
		"default_attributes": map[string]interface{}{
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
		"cookbook_versions": map[string]interface{}{
			"openssh": "= 11.0.0",
			"couchdb": "~> 1.2.0",
		},
	})
)

func TestEnvironmentName(t *testing.T) {
	// BUG(spheromak): Pull these constructors out into a Convey Decorator
	e1 := testEnvironmentMapStringInterfaceLol // (*Environment)
	name := e1.Name

	Convey("Environment name is 'testenvironment'", t, func() {
		So(name, ShouldEqual, "testenvironment")
	})

	swordWithoutASheathe, err := NewEnvironment(&Reader{
		"foobar": "baz",
	})

	name = swordWithoutASheathe.Name
	Convey("Environment without a name", t, func() {
		So(name, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestEnvironmentCookbookVersions(t *testing.T) {
	e := testEnvironmentMapStringInterfaceLol
	cv := e.CookbookVersion

	var
	for key
	Convey("Environment.Cookbook() should be a Cookbook", t, func() {
		So(cv, ShouldHaveSameTypeAs, Cookbook{})
	})

	fmt.Println(cv)
	// Convey("Environment.Cookbook() should be populated", t, func() {
	// 	So(cv, ShouldContain, "openssh")
	// 	So(cv, ShouldContain, "couchdb")
	// })
	//
	// rl = RunList{}
	// Convey("Empty RunList should be valid", t, func() {
	// 	So(rl, ShouldBeEmpty)
	// })

}

// func TestNodeAttribute(t *testing.T) {
// 	n := testEnvironmentMapStringInterfaceLol
// 	attr := n.Normal
// 	// BUG(spheromak): Holy shit this is ugly. Need to do something to make this easier for sure.
// 	ugh := attr["openssh"].(map[string]interface{})["server"].(map[string]string)["permit_root_login"]
// 	Convey("Node.Normal should map", t, func() {
// 		So(ugh, ShouldEqual, "no")
// 	})
// }

//
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

//
// // TestNewNode checks the NewNode Reader chain for Type
// func TestNewNode(t *testing.T) {
// 	var v interface{}
// 	v = testEnvironmentMapStringInterfaceLol
// 	Convey("NewNode should create a Node", t, func() {
// 		So(v, ShouldHaveSameTypeAs, &Node{})
// 	})
//
// 	Convey("NewNode should error if decode fails", t, func() {
//
// 		failNode, err := NewNode(&Reader{
// 			"name": struct{}{},
// 		})
//
// 		So(err, ShouldNotBeNil)
// 		So(failNode, ShouldBeNil)
// 	})
// }
//
// // TestNodeReadIntoFile tests that Read() can be used to read by io.Readers
// // BUG(fujin): re-do with goconvey
// func TestNodeReadIntoFile(t *testing.T) {
// 	e1 := testEnvironmentMapStringInterfaceLol // (*Node)
// 	tf, _ := ioutil.TempFile("test", "node-to-file")
// 	// Copy to tempfile (I use Read() internally)
// 	// BUG(fujin): this is currently doing that weird 32768 bytes read thing again.
// 	io.Copy(tf, e1)
//
// 	// Close and remove tempfile
// 	tf.Close()
// 	os.Remove(path.Clean(tf.Name()))
// }
