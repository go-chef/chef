//
// Test the go-chef/chef chef server api /organizations/*/cookbooks endpoints against a live chef server
//
package main

import (
	"fmt"
	"os"

	// chef "github.com/go-chef/chef"
	"chefapi_test/testapi"
)


// main Exercise the chef server api
func main() {
        // Create a client for user access
	client := testapi.Client()

	// Prep by adding a couple versions of some cookbooks before running this code
	// testbook version 0.1.0 and 0.2.0
	// sampbook version 0.1.0 and 0.2.0

	fmt.Println("Starting cookbooks")
	// Cookbooks
	cookList, err := client.Cookbooks.List()
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue listing cookbooks:", err)
        }
	fmt.Printf("List initial cookbooks %+v\nEndInitialList\n", cookList)

	// cookbook GET info
	testbook, err := client.Cookbooks.Get("testbook")
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue getting cookbook testbook:", err)
        }
	fmt.Printf("Get cookbook testbook %+v\n", testbook)

	// GET missing cookbook
	nothere, err := client.Cookbooks.Get("nothere")
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue getting cookbook nothere:", err)
        }
	fmt.Printf("Get cookbook nothere %+v\n", nothere)

	// list available versions of a cookbook
	testbookversions, err := client.Cookbooks.GetAvailableVersions("testbook", "0")
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue getting cookbook versions for testbook:", err)
        }
	fmt.Printf("Get cookbook versions testbook %+v\n", testbookversions)

	// list available versions of a cookbook
	sampbookversions, err := client.Cookbooks.GetAvailableVersions("sampbook", "0")
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue getting cookbook versions for sampbook:", err)
        }
	fmt.Printf("Get cookbook versions sampbook %+v\n", sampbookversions)

	// get specific versions of a cookbook
	testbookversions1, err := client.Cookbooks.GetVersion("testbook", "0.1.0")
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue getting specific cookbook versions for testbook:", err)
        }
	fmt.Printf("Get specific cookbook version testbook %+v\n", testbookversions1)

	// list all recipes
	allRecipes, err := client.Cookbooks.ListAllRecipes()
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue getting all recipes:", err)
        }
	fmt.Printf("Get all recipes %+v\n", allRecipes)

	// delete version
	err = client.Cookbooks.Delete("testbook", "0.1.0")
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue deleting testbook 0.1.0:", err)
        }
	fmt.Printf("Delete testbook 0.1.0 %+v\n", err)

	// delete version
	err = client.Cookbooks.Delete("testbook", "0.2.0")
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue deleting testbook 0.2.0:", err)
        }
	fmt.Printf("Delete testbook 0.2.0 %+v\n", err)

	// List cookbooks
	cookList, err = client.Cookbooks.List()
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue listing cookbooks:", err)
        }
	fmt.Printf("Final cookbook list %+v\n", cookList)

	// list available versions of a cookbook
	sampbookversions, err = client.Cookbooks.GetAvailableVersions("sampbook", "0")
	if err != nil {
                fmt.Fprintln(os.Stderr, "Issue getting cookbook versions for sampbook:", err)
        }
	fmt.Printf("Final cookbook versions sampbook %+v\n", sampbookversions)
}
