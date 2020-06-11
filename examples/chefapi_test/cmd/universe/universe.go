//
// Test the go-chef/chef chef server api /universe endpoints against a live server
//
package testapi

import (
<<<<<<< HEAD:testapi/universe.go
	"fmt"
)

// universe exercise the chef server api
func Universe() {
	// Create a client for access
	client := Client()
=======
	"chefapi_test/testapi"
	"fmt"
)

// main Exercise the chef server api
func main() {
	// Create a client for access
	client := testapi.Client()
>>>>>>> master:examples/chefapi_test/cmd/universe/universe.go

	universe, err := client.Universe.Get()
	if err != nil {
		fmt.Println("Issue getting universe information", err)
	}
	fmt.Printf("List universe: %+v", universe)
}
