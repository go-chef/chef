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
	testNodeJSON = "test/node.json"
	// FML
	testNodeMapStringInterfaceLol, _ = NewNode(&Reader{
		"name":       "test",
		"run_list":   []string{"recipe[foo]", "recipe[baz]", "role[banana]"},
		"chef_type":  "node",
		"json_class": "Chef::Node",
		"normal": map[string]interface{}{
			"tags": map[string]interface{}{},
			"openssh": map[string]interface{}{
				"server": map[string]string{
					"permit_root_login": "no",
					"max_auth_tries":    "3",
				},
			},
		},
		"override": map[string]interface{}{
			"openssh": map[string]interface{}{
				"server": map[string]string{
					"permit_root_login": "yes",
					"max_auth_tries":    "1",
				},
			},
		},
	})
)

func TestNodeName(t *testing.T) {
	// BUG(spheromak): Pull these constructors out into a Convey Decorator
	n1 := testNodeMapStringInterfaceLol // (*Node)
	name := n1.Name

	Convey("Node name is 'test'", t, func() {
		So(name, ShouldEqual, "test")
	})

	swordWithoutASheathe, err := NewNode(&Reader{
		"foobar": "baz",
	})

	name = swordWithoutASheathe.Name
	Convey("Node without a name", t, func() {
		So(name, ShouldBeEmpty)
		So(err, ShouldBeNil)
	})
}

func TestNodeRunList(t *testing.T) {
	n := testNodeMapStringInterfaceLol
	rl := n.RunList

	Convey("Node.RunList() should be a RunList", t, func() {
		So(rl, ShouldHaveSameTypeAs, RunList{})
	})

	Convey("Node.RunList() should be populated", t, func() {
		So(rl, ShouldContain, "recipe[foo]")
		So(rl, ShouldContain, "recipe[baz]")
		So(rl, ShouldContain, "role[banana]")
	})

	rl = RunList{}
	Convey("Empty RunList should be valid", t, func() {
		So(rl, ShouldBeEmpty)
	})

}

func TestNodeAttribute(t *testing.T) {
	n := testNodeMapStringInterfaceLol
	attr := n.Normal
	// BUG(spheromak): Holy shit this is ugly. Need to do something to make this easier for sure.
	ugh := attr["openssh"].(map[string]interface{})["server"].(map[string]string)["permit_root_login"]
	Convey("Node.Normal should map", t, func() {
		So(ugh, ShouldEqual, "no")
	})
}

// BUG(fujin): re-do with goconvey
func TestNodeFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testNodeJSON); err != nil {
		t.Error("unexpected error", err, "during os.Open on", testNodeJSON)
	} else {
		dec := json.NewDecoder(file)
		var n Node
		if err := dec.Decode(&n); err == io.EOF {
			log.Println(n)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

// TestNewNode checks the NewNode Reader chain for Type
func TestNewNode(t *testing.T) {
	var v interface{}
	v = testNodeMapStringInterfaceLol
	Convey("NewNode should create a Node", t, func() {
		So(v, ShouldHaveSameTypeAs, &Node{})
	})

	Convey("NewNode should error if decode fails", t, func() {

		failNode, err := NewNode(&Reader{
			"name": struct{}{},
		})

		So(err, ShouldNotBeNil)
		So(failNode, ShouldBeNil)
	})
}

// TestNodeReadIntoFile tests that Read() can be used to read by io.Readers
// BUG(fujin): re-do with goconvey
func TestNodeReadIntoFile(t *testing.T) {
	n1 := testNodeMapStringInterfaceLol // (*Node)
	tf, _ := ioutil.TempFile("test", "node-to-file")
	// Copy to tempfile (I use Read() internally)
	// BUG(fujin): this is currently doing that weird 32768 bytes read thing again.
	io.Copy(tf, n1)

	// Close and remove tempfile
	tf.Close()
	os.Remove(path.Clean(tf.Name()))
}
