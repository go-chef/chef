// Test the go-chef/chef chef server api /policy_groups endpoints against a live server
package testapi

import (
	"fmt"
	"os"
)

// policy exercise the chef server api
func PolicyGroup() {
	client := Client(nil)

	// List policy_groups
	policygroupList, err := client.PolicyGroups.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue printing the existing policy_groups:", err)
	}
	fmt.Printf("List policy_groups %+v\n", policygroupList)

	// Get specific policy_group
	policygroupOut, err := client.PolicyGroups.Get("testgroup")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue getting policy_group testgroup:", err)
	}
	fmt.Printf("Get testgroup %+v\n", policygroupOut)

	// Get policy from policy group
	policyOut, err := client.PolicyGroups.GetPolicy("testgroup", "testsamp")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue getting policy testsamp:", err)
	}
	fmt.Printf("Get testgroup::testsamp %+v\n", policyOut)

	// Delete policy from policy group
	policyDel, err := client.PolicyGroups.DeletePolicy("testgroup", "testsamp")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting policy testsamp:", err)
	}
	fmt.Printf("Delete testgroup::testsamp %+v\n", policyDel)

	// Delete specific policy_group
	policygroupDel, err := client.PolicyGroups.Delete("testgroup")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting policy_group testgroup:", err)
	}
	fmt.Printf("Delete testgroup %+v\n", policygroupDel)

}
