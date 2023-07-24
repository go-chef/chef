// Test the go-chef/chef chef server api /policies endpoints against a live server
package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// policy exercise the chef server api
func Policy() {
	client := Client()

	// List policies
	policyList, err := client.Policies.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue printing the existing policies:", err)
	}
	fmt.Printf("List policies %+v\n", policyList)

	policyName, policy := firstPolicy(policyList)
	revisionID := firstRevision(policy)

	// Get policy
	policyOut, err := client.Policies.Get(policyName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue getting %+v err %+v\n", policyName, err)
	}
	fmt.Printf("Get %+v %+v\n", policyName, policyOut)

	// Get policy revision
	policyRevOut, err := client.Policies.GetRevisionDetails(policyName, revisionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue getting %+v err %+v\n", policyName, err)
	}
	fmt.Printf("Get %+v revision %+v\n", policyName, policyRevOut)

	// Delete a revision from a policy
	revOutDel, err := client.Policies.DeleteRevision(policyName, revisionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue deleting revision from %+v: %+v\n", policyName, err)
	}
	fmt.Printf("Delete revision %v from %+v %+v\n", revisionID, policyName, revOutDel)

	// Try to get a missing policy
	policyOutMissing, err := client.Policies.Get("nothere")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue getting nothere: %+v\n", err)
	}
	fmt.Printf("Get nothere %+v\n", policyOutMissing)

	// Delete a policy
	policyOutDel, err := client.Policies.Delete("testsamp2")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue deleting testsamp2: %+v\n", err)
	}
	fmt.Printf("Delete testsamp2 %+v\n", policyOutDel)
}

func firstPolicy(policyList chef.PoliciesGetResponse) (string, chef.Policy) {
	for key, val := range policyList {
		return key, val
	}
	return "", chef.Policy{}
}

func firstRevision(policy chef.Policy) string {
	for key, _ := range policy.Revisions {
		return key
	}
	return ""
}
