//
// Test the go-chef/chef chef server api /group endpoints against a live server
//

// TODO: add users and then add them to groups. Seems to fail. pivotal is maybe not a good test. adding fails silently
package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// group exercise the chef server api
func Group() {
	client := Client()

	// List the current groups
	groupList, err := client.Groups.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue printing the existing groups:", err)
	}
	fmt.Printf("List initial groups %+vEndInitialList\n", groupList)

	// Build a stucture to define a group
	group1 := chef.Group{
		Name:      "group1",
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
	cerr, err := chef.ChefError(err)
	if cerr != nil {
		fmt.Fprintln(os.Stderr, "Issue recreating group1:", cerr.StatusCode())
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
	groupupdate := chef.GroupUpdate{}
	groupupdate.Name = "group1"
	groupupdate.GroupName = "group1"
	groupupdate.Actors.Clients = groupOut.Clients
	groupupdate.Actors.Groups = []string{}
	groupupdate.Actors.Users = []string{"pivotal"}
	groupUpOut, err := client.Groups.Update(groupupdate)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating group1:", err)
	}
	fmt.Printf("Update group1 %+v\n", groupUpOut)

	// Get new group after update
	groupOut, err = client.Groups.Get("group1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue getting group1:", err)
	}
	fmt.Printf("Get group1 after update %+v\n", groupOut)

	// Update a group add groups and change the group name
	groupupdate = chef.GroupUpdate{}
	groupupdate.Name = "group1"
	groupupdate.GroupName = "group1-new"
	groupupdate.Actors.Clients = []string{}
	groupupdate.Actors.Groups = []string{"admins"}
	groupupdate.Actors.Users = []string{}
	groupUpOut, err = client.Groups.Update(groupupdate)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating group1:", err)
	}
	fmt.Printf("Update group1 %+v\n", groupUpOut)

	// Get new group after update and rename
	groupOut, err = client.Groups.Get("group1-new")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue getting group1-new:", err)
	}
	fmt.Printf("Get group1-new after update %+v\n", groupOut)

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
