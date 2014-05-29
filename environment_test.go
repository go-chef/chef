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
	testEnvironmentJSON = "test/environment.json"
	// FML
	testEnvironmentMapStringInterfaceLol, _ = NewEnvironment(&Reader{
		"name":       "testenvironment",
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
			"openssh": map[string]interface{}{
				"name":    "openssh",
				"version": "= 1.0.0",
			},
			"couchdb": map[string]interface{}{
				"name":    "couchdb",
				"version": "~> 1.2.0",
			},
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
	envCooks := testEnvironmentMapStringInterfaceLol
	Convey("CookbookVersion should have cooks", t, func() {
		So(envCooks.CookbookVersions["openssh"], ShouldHaveSameTypeAs, nativeCookbook{})
	})

}

func TestEnvironmentAttribute(t *testing.T) {
	n := testEnvironmentMapStringInterfaceLol
	attr := n.Default
	// BUG(spheromak): Holy shit this is ugly. Need to do something to make this easier for sure.
	ugh := attr["openssh"].(map[string]interface{})["server"].(map[string]string)["permit_root_login"]
	Convey("Environment.Default should map", t, func() {
		So(ugh, ShouldEqual, "no")
	})
}

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

// TestNewEnvironment checks the NewEnvironment Reader chain for Type
func TestNewEnvironment(t *testing.T) {
	var v interface{}
	v = testEnvironmentMapStringInterfaceLol
	Convey("NewEnvironment should create a Environment", t, func() {
		So(v, ShouldHaveSameTypeAs, &Environment{})
	})

	Convey("NewEnvironment should error if decode fails", t, func() {

		failEnvironment, err := NewEnvironment(&Reader{
			"name": struct{}{},
		})

		So(err, ShouldNotBeNil)
		So(failEnvironment, ShouldBeNil)
	})
}

// TestEnvironmentReadIntoFile tests that Read() can be used to read by io.Readers
func TestEnvironmentReadIntoFile(t *testing.T) {
	e1 := testEnvironmentMapStringInterfaceLol // (*Environment)
	tf, _ := ioutil.TempFile("test", "environment-to-file")
	Convey("Environment To File", t, func() {
		_, err := io.Copy(tf, e1)
		So(err, ShouldBeNil)
	})
	// Close and remove tempfile
	tf.Close()
	os.Remove(path.Clean(tf.Name()))
}
