package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const policyGroupResponseFile = "test/policy_group_response.json"
const policyGroupFile = "test/policy_group.json"
const revisionDetailsResponseFile = "test/revision_details_response.json"

func TestPolicyGroupList(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(policyGroupResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policy_groups", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.PolicyGroups.List()
	if err != nil {
		t.Error(err)
	}

	if data == nil {
		t.Fatal("We should have some data")
	}

	if len(data) != 1 {
		t.Error("Mismatch in expected policy group count. Expected 1, Got: ", len(data))
	}

	if _, ok := data["demo_policy_group"]; !ok {
		t.Error("demo_policy_group policy should be listed")
	}

}

func TestPolicyGroupGet(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(policyGroupFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policy_groups/testgroup", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.PolicyGroups.Get("testgroup")
	if err != nil {
		t.Error(err)
	}

	if _, ok := data.Policies["testsamp"]; !ok {
		t.Error("testsamp policy should be listed")
	}

}

func TestPolicyGroupDelete(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(policyGroupFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policy_groups/testgroup", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.PolicyGroups.Delete("testgroup")
	if err != nil {
		t.Error(err)
	}

	if _, ok := data.Policies["testsamp"]; !ok {
		t.Error("testsamp policy should be listed")
	}

}

func TestPolicyGroupGetPolicy(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(revisionDetailsResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policy_groups/testgroup/policies/testsamp", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.PolicyGroups.GetPolicy("testgroup", "testsamp")
	if err != nil {
		t.Error(err)
	}

	if data.Name != "testsamp" {
		t.Error("testsamp policy should be retrieved")
	}

}

func TestPolicyGroupDeletePolicy(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(revisionDetailsResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policy_groups/testgroup/policies/testsamp", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.PolicyGroups.DeletePolicy("testgroup", "testsamp")
	if err != nil {
		t.Error(err)
	}

	if data.Name != "testsamp" {
		t.Error("testsamp policy should be retrieved")
	}

}
