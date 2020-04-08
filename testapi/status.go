//
// Test the go-chef/chef chef server api /_status endpoints against a live server
//
package testapi

import (
	"fmt"
)

// status exercise the chef server api
func Status() {
	// Create a client for access
	client := Client()

	status, err := client.Status.Get()
	if err != nil {
		fmt.Println("Issue getting status information", err)
	}
	fmt.Printf("List status: %+v", status)
}
