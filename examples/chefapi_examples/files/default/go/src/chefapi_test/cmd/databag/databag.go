//
// Test the go-chef/chef chef databag api against a live server
//
package main

import (
	"fmt"
        "chefapi_test/testapi"
        "github.com/go-chef/chef"
	"os"
)


// main Exercise the chef server databag api
func main() {
	client := testapi.Client()

	// List data bags before
	bags, err := client.DataBags.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing data bags %+v\n",err)
	}
	fmt.Printf("List data bags before %+v\n", bags)


	// Create a data bag
	databag := &chef.DataBag{Name: "testbag"}
	response, err := client.DataBags.Create(databag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue creating data bag testbag %+v\n",err)
	}
	fmt.Printf("Data bag created %+v\n", response)

	// List data bags after
	bags, err = client.DataBags.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing data bags after add %+v\n",err)
	}
	fmt.Printf("List data bags after create %+v\n", bags)

	// Create a data bag item
	dbi := map[string]string{
		"id": "item1",
		"foo": "test123",
	}
	err = client.DataBags.CreateItem("testbag", dbi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue creating item %+v\n",err)
	}
	fmt.Printf("Created item %+v\n", dbi)

	// List data bag items
	bagItems, err := client.DataBags.ListItems("testbag")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing items %+v\n",err)
	}
	fmt.Printf("List bag items %+v\n", bagItems)

	// Get data bag items
	itemOut, err := client.DataBags.GetItem("testbag", "item1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue getting initial item %+v\n",err)
	}
	fmt.Printf("Initial item %+v\n", itemOut)

	// Update a data bag item
	dbi = map[string]string{
		"id": "item1",
		"foo": "update123",
	}
	err = client.DataBags.UpdateItem("testbag", "item1", dbi)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue updating item %+v\n",err)
	}
	fmt.Printf("Update item %+v\n", dbi)

	// Get data bag items
	itemOut, err = client.DataBags.GetItem("testbag", "item1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue getting updated item %+v\n",err)
	}
	fmt.Printf("Updated item %+v\n", itemOut)

	// Delete a data bag item
	err = client.DataBags.DeleteItem("testbag", "item1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue deleting item %+v\n",err)
	}
	fmt.Println("Deleted item")

	// List data bag items
	bagItems, err = client.DataBags.ListItems("testbag")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing items after delete %+v\n",err)
	}
	fmt.Printf("List bag items after delete  %+v\n", bagItems)

	// Delete a data bag
	outBag, err := client.DataBags.Delete("testbag")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue deleting data bag testbag %+v\n",err)
	}
	fmt.Printf("Data bag deleted %+v\n", outBag)

	// List data bags after delete
	bags, err = client.DataBags.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing data bags after delete %+v\n",err)
	}
	fmt.Printf("List data bags after delete %+v\n", bags)
}
