package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const policyListResponseFile = "test/policies_response.json"

func TestPolicyList(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(policyListResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/policies", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.Policies.List()
	if err != nil {
		t.Error(err)
	}

	if data == nil {
		t.Fatal("We should have some data")
	}
	fmt.Println(data)

	if len(data) != 2 {
		t.Error("Mismatch in expected policies count. Expected 2, Got: ", len(data))
	}

	if _, ok := data["aar"]; !ok {
		t.Error("aar policy should be listed")
	}

	if _, ok := data["jenkins"]; !ok {
		t.Error("jenkins policy should be listed")
	}

}
