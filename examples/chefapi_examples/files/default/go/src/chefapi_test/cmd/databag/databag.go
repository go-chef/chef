//
// Test the go-chef/chef chef server api /databag endpoints against a live server
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

	// List the current databags
	BagList, err := client.DataBag.List()
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue printing the existing databags:", err)
	}
	fmt.Printf("List initial databags %+v\n", BagList)

	databag1 := chef.Databag{
		Name: "databag1",
	}

	// Add a new databag
	databagAdd, err := client.DataBags.Create(databag1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue adding databag1:", err)
	}
	fmt.Println("Added databag1", databagAdd)

	// Add databag again
	databagAdd, err = client.DataBags.Create(databag1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue recreating databag1:", err)
	}
	fmt.Println("Recreated databag1", databagAdd)

	// Try to get a missing databag 
	databagOutMissing, err := client.DataBags.Get("nothere")
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue getting nothere:", err)
	}
	fmt.Println("Get nothere", databagOutMissing)

	// List databags after adding
	databagList, err = client.DataBags.List()
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue printing the existing databags:", err)
	}
	fmt.Printf("List databags after adding databag1 %+v\n", databagList)

	// Add items to a data bag

	// Update a databag item
	databag1.BagName = "databag1-new"
	databag1.Users = append(databag1.Users, "pivotal")
	databagUpdate, err := client.DataBags.Update(databag1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue updating databag1:", err)
	}
	fmt.Printf("Update databag1 %+v\n", databagUpdate)

	// list databag items
	databagOut, err := client.DataBags.Get("databag1")
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue getting databag1:", err)
	}
	fmt.Printf("Get databag1 %+v\n", databagOut)


	// Get the contents of a data bag item

	// Delete a databag item

	// List items

        // Clean up
	err = client.DataBags.Delete("databag1-new")
	fmt.Println("Delete databag1", err)

	// Final list of databags
	databagList, err = client.DataBags.List()
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue listing the final databags:", err)
	}
	fmt.Printf("List databags after cleanup %+v\n", databagList)
}
