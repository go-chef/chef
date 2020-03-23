//
// Test the go-chef/chef chef server api /group endpoints against a live server
//
package main

import (
	"fmt"
        "chefapi_test/testapi"
        "github.com/go-chef/chef"
	"os"
)


// main Exercise the chef server api
func main() {
	client := testapi.Client()

	// List the current groups
	groupList, err := client.Groups.List()
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue printing the existing groups:", err)
	}
	fmt.Printf("List initial groups %+vEndInitialList\n", groupList)

	// Build a stucture to define a group
        group1 := chef.Group {
		Name: "group1",
		GroupName: "group1",
	}

	// Add a new group
	groupAdd, err := client.Groups.Create(group1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue adding group1:", err)
	}
	fmt.Println("Added group1", groupAdd)

	// Add group again
	groupAdd, err = client.Groups.Create(group1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue recreating group1:", err)
	}
	fmt.Println("Recreated group1", groupAdd)

	// List groups after adding
	groupList, err = client.Groups.List()
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue printing the existing groups:", err)
	}
	fmt.Printf("List groups after adding group1 %+vEndAddList\n", groupList)

	// Get new group
	groupOut, err := client.Groups.Get("group1")
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue getting group1:", err)
	}
	fmt.Printf("Get group1 %+v\n", groupOut)

	// Try to get a missing group 
	groupOutMissing, err := client.Groups.Get("nothere")
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue getting nothere:", err)
	}
	fmt.Println("Get nothere", groupOutMissing)

	// Update a group
	group1.GroupName = "group1-new"
	group1.Users = append(group1.Users, "pivotal")
	groupUpdate, err := client.Groups.Update(group1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue updating group1:", err)
	}
	fmt.Printf("Update group1 %+v\n", groupUpdate)

        // Clean up
	err = client.Groups.Delete("group1-new")
	fmt.Println("Delete group1", err)

	// Final list of groups
	groupList, err = client.Groups.List()
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue listing the final groups:", err)
	}
	fmt.Printf("List groups after cleanup %+vEndFinalList\n", groupList)
}
