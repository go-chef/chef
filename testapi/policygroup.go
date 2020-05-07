//
// Test the go-chef/chef chef server api /policy_groups endpoints against a live server
//
package testapi

import (
	"fmt"
	"os"
)

// policy exercise the chef server api
func PolicyGroup() {
	client := Client()

	// List policy_groups
	policygroupList, err := client.PolicyGroups.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue printing the existing policy_groups:", err)
	}
	fmt.Printf("List policy_groups %+v\n", policygroupList)
}
