//
// Test the go-chef/chef chef server api /license endpoints against a live server
//
package testapi

import (
<<<<<<< HEAD:testapi/license.go
	"fmt"
)

// license exercise the chef server api
func License() {
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
>>>>>>> master:examples/chefapi_test/cmd/license/license.go

	license, err := client.License.Get()
	if err != nil {
		fmt.Println("Issue getting license information", err)
	}
	fmt.Printf("List license: %+v", license)
}
