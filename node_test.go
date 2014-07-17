package chef

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
	// . "github.com/smartystreets/goconvey/convey"
)

var (
	testNodeJSON = "test/node.json"
)

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

func TestNodesService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"node1":"https://chef/nodes/node1", "node2":"https://chef/nodes/node2"}`)
	})

	nodes, err := client.Nodes.List()
	if err != nil {
		t.Errorf("Nodes.List returned error: %v", err)
	}

	want := map[string]string{"node1": "https://chef/nodes/node1", "node2": "https://chef/nodes/node2"}

	if !reflect.DeepEqual(nodes, want) {
		t.Errorf("Nodes.List returned %+v, want %+v", nodes, want)
	}
}

func TestNodesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/nodes/node1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
	    "name": "node1",
	    "json_class": "Chef::Node",
	    "chef_type": "node",
	    "chef_environment": "development"
		}`)
	})

	nodes, err := client.Nodes.Get("node1")
	if err != nil {
		t.Errorf("Nodes.Get returned error: %v", err)
	}

	want := Node{
		Name:        "node1",
		JsonClass:   "Chef::Node",
		ChefType:    "node",
		Environment: "development",
	}

	if !reflect.DeepEqual(nodes, want) {
		t.Errorf("Nodes.Get returned %+v, want %+v", nodes, want)
	}
}
