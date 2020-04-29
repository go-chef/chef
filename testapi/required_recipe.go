//
// Test the go-chef/chef chef server api /required_recipe endpoint against a live server
//
package testapi

import (
	"fmt"
)

// required_recipe exercise the chef server api
func RequiredRecipe() {
	// Create a client for access
	client := Client()

	required_recipe, err := client.RequiredRecipe.Get()
	if err != nil {
		fmt.Println("Issue getting required_recipe information", err)
	}
	fmt.Printf("List required_recipe: %+v", required_recipe)
}
