package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const policyGroupResponseFile = "test/policy_group_response.json"

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
