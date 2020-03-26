//
// Test the go-chef/chef chef server api /databag endpoints against a live server
//
package main

import (
	"chefapi_test/testapi"
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// main Exercise the chef server api
func main() {
	client := testapi.Client()

	// List the current databags
	BagList, err := client.DataBags.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue printing the existing databags:", err)
	}
	fmt.Printf("List initial databags %+v\n", BagList)

	databag1 := chef.DataBag{
		Name: "databag1",
	}

	// Add a new databag
	databagAdd, err := client.DataBags.Create(&databag1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue adding databag1:", err)
	}
	fmt.Println("Added databag1", databagAdd)

	// Add databag again
	databagAdd, err = client.DataBags.Create(&databag1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue recreating databag1:", err)
	}
	fmt.Println("Recreated databag1", databagAdd)

	// Try to get a missing databag
	databagOutMissing, err := client.DataBags.ListItems("nothere")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue getting nothere:", err)
	}
	fmt.Println("Get nothere", databagOutMissing)

	// List databags after adding
	BagList, err = client.DataBags.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue printing the existing databags:", err)
	}
	fmt.Printf("List databags after adding databag1 %+v\n", BagList)

	// Add item to a data bag
	item1data := map[string]string{
		"id":    "item1",
		"type":  "password",
		"value": "secret",
	}
	item1 := chef.DataBagItem(item1data)
	err = client.DataBags.CreateItem("databag1", item1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating databag1:", err)
	}
	fmt.Printf("Create databag1::item1 %+v\n", err)

	// Update a databag item
	item1data["value"] = "next"
	item1upd := chef.DataBagItem(item1data)
	err = client.DataBags.UpdateItem("databag1", "item1", item1upd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating databag1:", err)
	}
	fmt.Printf("Update databag1::item1 %+v\n", err)

	// list databag items
	databagItems, err := client.DataBags.ListItems("databag1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue getting databag1:", err)
	}
	fmt.Printf("List databag1 items %+v\n", databagItems)

	// Get the contents of a data bag item
	dataItem, err := client.DataBags.GetItem("databag1", "item1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue getting databag1::item1 :", err)
	}
	fmt.Printf("Get databag1::item1 %+v\n", dataItem)

	// Delete a databag item
	err = client.DataBags.DeleteItem("databag1", "item1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting databag1::item1 :", err)
	}
	fmt.Printf("Delete databag1::item1 %+v\n", err)

	// List items
	databagItems, err = client.DataBags.ListItems("databag1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue getting databag1:", err)
	}
	fmt.Printf("List databag1 items after delete %+v\n", databagItems)

	// Clean up
	databag, err := client.DataBags.Delete("databag1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting databag1:", err)
	}
	fmt.Printf("Delete databag1 %+v\n", databag)

	// Final list of databags
	BagList, err = client.DataBags.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing the final databags:", err)
	}
	fmt.Printf("List databags after cleanup %+v\n", BagList)
}
