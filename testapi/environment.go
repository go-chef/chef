package testapi

import (
	"fmt"
	"os"

	"github.com/go-chef/chef"
)

// environment exercise the chef api
func Environment() {
	// Use the default test org
	client := Client(nil)

	// List initial environments
	environmentList, err := client.Environments.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list environments: ", err)
	}
	fmt.Println("List initial environments", environmentList)

	// Define a Environment object
	// TODO set some more things
	environment1 := chef.Environment{Name: "environment1"}
	fmt.Println("Define environment1", environment1)

	// Delete environment1 ignoring errors :)
	_, err = client.Environments.Delete(environment1.Name)

	// Create
	environmentResult, err := client.Environments.Create(&environment1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create environment environment1. ", err)
	}
	fmt.Println("Added environment1", environmentResult)

	// List environments
	environmentList, err = client.Environments.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list environments: ", err)
	}
	fmt.Println("List environments after adding environment1", environmentList)

	// Create a second time
	environmentResult, err = client.Environments.Create(&environment1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't recreate environment environment1. ", err)
	}
	fmt.Println("Added environment1", environmentResult)

	// Read environment1 information
	serverEnvironment, err := client.Environments.Get("environment1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get environment: ", err)
	}
	fmt.Printf("Get environment1 %+v\n", serverEnvironment)

	// update environment
	// TODO update the environment
	updateEnvironment, err := client.Environments.Put(&environment1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't update environment: ", err)
	}
	fmt.Println("Update environment1", updateEnvironment)

	// Info after update
	serverEnvironment, err = client.Environments.Get("environment1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get environment: ", err)
	}
	fmt.Printf("Get environment1 after update %+v\n", serverEnvironment)

	// Delete environment ignoring errors :)
	_, err = client.Environments.Delete(environment1.Name)
	fmt.Println("Delete environment1", err)

	// List environments
	environmentList, err = client.Environments.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list environments: ", err)
	}
	fmt.Println("List environments after cleanup", environmentList)
}
