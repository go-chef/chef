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

	"github.com/stretchr/testify/assert"
)

var (
	testNodeJSON = "test/node.json"
)

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

func TestNode_NewNode(t *testing.T) {
	n := NewNode("testnode")
	expect := Node{
		Name:        "testnode",
		Environment: "_default",
		ChefType:    "node",
		JsonClass:   "Chef::Node",
	}

	if !reflect.DeepEqual(n, expect) {
		t.Errorf("NewNode returned %+v, want %+v", n, expect)
	}
}

func TestNodesService_Methods(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/nodes", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{"node1":"https://chef/nodes/node1", "node2":"https://chef/nodes/node2"}`)
		case r.Method == "POST":
			fmt.Fprintf(w, `{ "uri": "http://localhost:4545/nodes/node1" }`)
		}
	})

	mux.HandleFunc("/nodes/node1", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" || r.Method == "PUT":
			fmt.Fprintf(w, `{
	    "name": "node1",
	    "json_class": "Chef::Node",
	    "chef_type": "node",
	    "chef_environment": "development"
		}`)
		case r.Method == "HEAD":
		case r.Method == "DELETE":
		}
	})

	// Test list
	nodes, err := client.Nodes.List()
	if err != nil {
		t.Errorf("Nodes.List returned error: %v", err)
	}

	listWant := map[string]string{"node1": "https://chef/nodes/node1", "node2": "https://chef/nodes/node2"}

	if !reflect.DeepEqual(nodes, listWant) {
		t.Errorf("Nodes.List returned %+v, want %+v", nodes, listWant)
	}

	// test Get
	node, err := client.Nodes.Get("node1")
	if err != nil {
		t.Errorf("Nodes.Get returned error: %v", err)
	}

	wantNode := NewNode("node1")
	wantNode.Environment = "development"
	if !reflect.DeepEqual(node, wantNode) {
		t.Errorf("Nodes.Get returned %+v, want %+v", node, wantNode)
	}

	// test HEAD
	err = client.Nodes.Head("node1")
	if err != nil {
		t.Errorf("Nodes.Head returned error: %v", err)
	}

	// test Post
	res, err := client.Nodes.Post(wantNode)
	if err != nil {
		t.Errorf("Nodes.Post returned error: %s", err.Error())
	}

	postResult := &NodeResult{"http://localhost:4545/nodes/node1"}
	if !reflect.DeepEqual(postResult, res) {
		t.Errorf("Nodes.Post returned %+v, want %+v", res, postResult)
	}

	// test Put
	putRes, err := client.Nodes.Put(node)
	if err != nil {
		t.Errorf("Nodes.Put returned error %+v", err)
	}

	if !reflect.DeepEqual(putRes, node) {
		t.Errorf("Nodes.Post returned %+v, want %+v", putRes, node)
	}

	// test Delete
	err = client.Nodes.Delete(node.Name)
	if err != nil {
		t.Errorf("Nodes.Delete returned error %+v", err)
	}
}

func TestGetAttribute(t *testing.T) {
	tests := []struct {
		name   string
		node   Node
		paths  []string
		result interface{}
		err    error
	}{
		{
			"no path",
			Node{
				NormalAttributes: map[string]interface{}{
					"foo": "bar",
				},
			},
			[]string{},
			nil,
			ErrNoPathProvided,
		},
		{
			"normal attribute",
			Node{
				NormalAttributes: map[string]interface{}{
					"foo": "bar",
				},
			},
			[]string{"foo"},
			"bar",
			nil,
		},
		{
			"nested normal attribute",
			Node{
				NormalAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "foobar",
					},
				},
			},
			[]string{"foo", "bar"},
			"foobar",
			nil,
		},
		{
			"missing normal attribute",
			Node{
				NormalAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"buzz": "foobar",
					},
				},
			},
			[]string{"foo", "bar"},
			nil,
			ErrPathNotFound,
		},
		{
			"automatic first",
			Node{
				AutomaticAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 1,
					},
				},
				OverrideAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 2,
					},
				},
			},
			[]string{"foo", "bar"},
			1,
			nil,
		},
		{
			"override first",
			Node{
				OverrideAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 1,
					},
				},
				NormalAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 2,
					},
				},
			},
			[]string{"foo", "bar"},
			1,
			nil,
		},
		{
			"normal first",
			Node{
				NormalAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 1,
					},
				},
				DefaultAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 2,
					},
				},
			},
			[]string{"foo", "bar"},
			1,
			nil,
		},
		{
			"correct order",
			Node{
				AutomaticAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 4,
					},
				},
				OverrideAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 3,
					},
				},
				NormalAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 2,
					},
				},
				DefaultAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": 1,
					},
				},
			},
			[]string{"foo", "bar"},
			4,
			nil,
		},
		{
			"first full path",
			Node{
				OverrideAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"buzz": "foobuzz",
					},
				},
				NormalAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": "foobar",
					},
				},
			},
			[]string{"foo", "bar"},
			"foobar",
			nil,
		},
		{
			"incomplete path",
			Node{
				NormalAttributes: map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": map[string]interface{}{
							"foobar": "buzz",
						},
					},
				},
			},
			[]string{"foo", "bar"},
			map[string]interface{}{"foobar": "buzz"},
			nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.node.GetAttribute(tt.paths...)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.result, result)
		})
	}
}
