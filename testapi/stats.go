//
// Test the go-chef/chef chef server api /stats endpoint against a live server
//
package testapi

import (
	"fmt"
	"os"
)

// stats exercise the chef server api
func Stats() {
	// Create a client for access
	client := Client()
	password  := os.Args[6]
	fmt.Println("password", password)

	stats, err := client.Stats.Get("json", "statsuser", password)
	if err != nil {
		fmt.Println("Issue getting stats information", err)
	}
	fmt.Printf("List stats json format: %+v", stats)
}
