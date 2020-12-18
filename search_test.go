package chef

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

const partialSearchResponseFile_1 = "test/partial_search_test_1.json"
const partialSearchResponseFile_2 = "test/partial_search_test_2.json"

func TestSearch_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
			"node": "http://localhost:4000/search/node", 
			"role": "http://localhost:4000/search/role", 
			"client": "http://localhost:4000/search/client", 
			"users": "http://localhost:4000/search/users" 
		}`)
	})

	indexes, err := client.Search.Indexes()
	if err != nil {
		t.Errorf("Search.Get returned error: %+v", err)
	}
	wantedIdx := map[string]string{
		"node":   "http://localhost:4000/search/node",
		"role":   "http://localhost:4000/search/role",
		"client": "http://localhost:4000/search/client",
		"users":  "http://localhost:4000/search/users",
	}
	if !reflect.DeepEqual(indexes, wantedIdx) {
		t.Errorf("Search.Get returned %+v, want %+v", indexes, wantedIdx)
	}
}

func TestSearch_ExecDo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/nodes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
	    "total": 1,
	    "start": 0,
	    "rows": [
	       {
	        "overrides": {"hardware_type": "laptop"},
	        "name": "latte",
	        "chef_type": "node",
	        "json_class": "Chef::Node",
	        "attributes": {"hardware_type": "laptop"},
	        "run_list": ["recipe[unicorn]"],
	        "defaults": {}
	       }
				 ]
		}`)
	})

	// test the fail case
	_, err := client.Search.NewQuery("foo", "failsauce")
	if err == nil {
		t.Error("Bad query wasn't caught")
	}

	// test the fail case
	_, err = client.Search.Exec("foo", "failsauce")
	if err == nil {
		t.Error("Bad query wasn't caught")
	}

	// test the positive case
	query, err := client.Search.NewQuery("nodes", "name:latte")
	if err != nil {
		t.Error("failed to create query")
	}

	// for now we aren't testing the result..
	_, err = query.Do(client)
	if err != nil {
		t.Errorf("Search.Exec failed err: %+v", err)
	}

	_, err = client.Search.Exec("nodes", "name:latte")
	if err != nil {
		t.Errorf("Search.Exec failed err: %+v", err)
	}

}

func TestSearch_PartialExec(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/search/node", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
			"total": 1,
			"start": 0,
			"rows": [
			   {
				"overrides": {"hardware_type": "laptop"},
				"name": "latte",
				"chef_type": "node",
				"json_class": "Chef::Node",
				"policy_group": "testing",
				"policy_name": "grafana",
				"policy_revision": "123xyz00009999",
				"attributes": {"hardware_type": "laptop"},
				"run_list": ["recipe[unicorn]"],
				"defaults": {}
			   }
					 ]
			}`)
	})

	query := map[string]interface{}{
		"name":            []string{"name"},
		"policy_group":    []string{"policy_group"},
		"policy_name":     []string{"policy_name"},
		"policy_revision": []string{"policy_revision"},
	}

	pres, err := client.Search.PartialExec("node", "*.*", query)
	if err != nil {
		t.Errorf("Search.PartialExec failed err: %+v", err)
	}

	assert.Len(t, pres.Rows, 1)
	actualNode := Node{}
	assert.NoError(t, json.Unmarshal(pres.Rows[0], &actualNode))
	assert.Equal(t, "grafana", actualNode.PolicyName)

}

func TestSearch_PartialExecMultipleCalls(t *testing.T) {
	setup()
	defer teardown()

	searchResponseOne, err := ioutil.ReadFile(partialSearchResponseFile_1)
	if err != nil {
		t.Error(err)
	}

	searchResponseTwo, err := ioutil.ReadFile(partialSearchResponseFile_2)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/search/node", func(w http.ResponseWriter, r *http.Request) {

		start, ok := r.URL.Query()["start"]

		if !ok || len(start[0]) < 1 {
			fmt.Println("Url Param 'start' is missing")
			return
		}

		if start[0] == "0" {
			fmt.Fprintf(w, string(searchResponseOne))
		} else {
			fmt.Fprintf(w, string(searchResponseTwo))
		}
	})

	query := map[string]interface{}{
		"name":            []string{"name"},
		"policy_group":    []string{"policy_group"},
		"policy_name":     []string{"policy_name"},
		"policy_revision": []string{"policy_revision"},
	}

	pres, err := client.Search.PartialExec("node", "*.*", query)
	if err != nil {
		t.Errorf("Search.PartialExec failed err: %+v", err)
	}

	assert.Len(t, pres.Rows, 1185)

	firstNode := Node{}
	assert.NoError(t, json.Unmarshal(pres.Rows[0], &firstNode))
	assert.Equal(t, "node1", firstNode.Name)

	lastNode := Node{}
	assert.NoError(t, json.Unmarshal(pres.Rows[len(pres.Rows)-1], &lastNode))
	assert.Equal(t, "node1185", lastNode.Name)

}
