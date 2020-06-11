//
// Test the go-chef/chef chef server api /_status endpoints against a live server
//
package testapi

import (
<<<<<<< HEAD:testapi/status.go
	"fmt"
)

// status exercise the chef server api
func Status() {
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
>>>>>>> master:examples/chefapi_test/cmd/status/status.go

	status, err := client.Status.Get()
	if err != nil {
		fmt.Println("Issue getting status information", err)
	}
	fmt.Printf("List status: %+v", status)
}
