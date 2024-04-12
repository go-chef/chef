package chef

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

const policyListResponseFile = "test/policies_response.json"
const policyRevisionResponseFile = "test/policy_revision_response.json"

func TestListPolicies(t *testing.T) {
	setup()
	defer teardown()

	file, err := os.ReadFile(policyListResponseFile)
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

func TestGetPolicy(t *testing.T) {
	setup()
	defer teardown()

	policyGetJSON := `{
						"revisions": {
		  					"8228b0e381fe1de3ee39bf51e93029dbbdcecc61fb5abea4ca8c82591c0b529b": {

		  						}
							}
	  				 }`
	mux.HandleFunc("/policies/base", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, policyGetJSON)
	})
	mux.HandleFunc("/policies/bad", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", 404)
	})

	data, err := client.Policies.Get("base")
	if err != nil {
		t.Error(err)
	}

	if _, ok := data["revisions"]["8228b0e381fe1de3ee39bf51e93029dbbdcecc61fb5abea4ca8c82591c0b529b"]; !ok {
		t.Error("Missing expected revision for this policy")
	}

	_, err = client.Policies.Get("bad")
	if err == nil {
		t.Error("We expected this bad request to error", err)
	}
}

func TestGetPolicyRevision(t *testing.T) {
	setup()
	defer teardown()

	const policyName = "base"
	const policyRevision = "8228b0e381fe1de3ee39bf51e93029dbbdcecc61fb5abea4ca8c82591c0b529b"

	file, err := os.ReadFile(policyRevisionResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc(fmt.Sprintf("/policies/%s/revisions/%s", policyName, policyRevision), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.Policies.GetRevisionDetails(policyName, policyRevision)
	if err != nil {
		t.Error(err)
	}

	if data.Name != policyName {
		t.Errorf("Unexpected policy name: %+v", data.Name)
	}

	if data.RevisionID != policyRevision {
		t.Errorf("Unexpected policy revision ID: %+v", data.RevisionID)
	}

	if data.RunList[0] != "recipe[base::default]" {
		t.Errorf("Unexpected policy run list: %+v", data.RevisionID)
	}

	if val, ok := data.NamedRunList["os"]; !ok {
		t.Error("Expected os NamedRunList policy to be present in the policy information")
	} else if val[0] != "recipe[hardening::default]" {
		t.Error("Expected named run list for the policy, got: ", val[0])
	}

	if data.IncludedPolicyLocks[0].Name != "other" {
		t.Error("Expected included policy name to be present in the policy information")
	} else if data.IncludedPolicyLocks[0].RevisionID != "7b40995ad1150ec56950c757872d6732aa00e76382dfcd2fddeb3a971e57ba9c" {
		t.Error("Expected included policy revision ID to be present in the policy information")
	}

	if val, ok := data.CookbookLocks["hardening"]; !ok {
		t.Error("Expected hardening policy to be present in the policy information")
	} else if val.Version != "0.1.0" {
		t.Error("Expected hardening policy version to be 0.1.0, got: ", val.Version)
	}

}

func TestDeletePolicyRevision(t *testing.T) {
	setup()
	defer teardown()

	const policyName = "base"
	const policyRevision = "8228b0e381fe1de3ee39bf51e93029dbbdcecc61fb5abea4ca8c82591c0b529b"

	file, err := os.ReadFile(policyRevisionResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc(fmt.Sprintf("/policies/%s/revisions/%s", policyName, policyRevision), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	data, err := client.Policies.DeleteRevision(policyName, policyRevision)
	if err != nil {
		t.Error(err)
	}

	if data.Name != policyName {
		t.Errorf("Unexpected policy name: %+v", data.Name)
	}

	if data.RevisionID != policyRevision {
		t.Errorf("Unexpected policy revision ID: %+v", data.RevisionID)
	}
}
