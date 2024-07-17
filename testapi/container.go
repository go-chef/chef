package testapi

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-chef/chef"
)

// container exercise the chef api
func Container() {
	// Use the default test org
	client := Client(nil)

	// List initial containers
	containerList, err := client.Containers.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list containers: ", err)
	}
	fold := strings.Replace(fmt.Sprintf("%+v", containerList), "\n", "", -1)
	fmt.Println("List initial containers", fold)

	// Define a Container object
	container1 := chef.Container{}
	container1.ContainerName = "container1"
	container1.ContainerPath = "container1"

	// Create
	containerResult, err := client.Containers.Create(container1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating container container1. ", err)
	}
	fmt.Println("Added container1", containerResult)

	// List containers
	containerList, err = client.Containers.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing containers after add: ", err)
	}
	fold = strings.Replace(fmt.Sprintf("%+v", containerList), "\n", "", -1)
	fmt.Println("List containers after adding container1", fold)

	// Read container1 information
	containerContents, err := client.Containers.Get("container1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get container1: ", err)
	}
	fmt.Printf("Get container1 %+v\n", containerContents)

	// Read the environment container information
	containerContents, err = client.Containers.Get("environments")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get container environments: ", err)
	}
	fmt.Printf("Get environment %+v\n", containerContents)

	// Delete container
	err = client.Containers.Delete(container1.ContainerName)
	fmt.Println("Delete container1", err)
}
