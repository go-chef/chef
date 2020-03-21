package main

import (
	"fmt"
	"os"

	"chefapi_test/testapi"
	"github.com/go-chef/chef"
)

func main() {
	// Use the default test org
	client := testapi.Client()

	// List initial roles
	roleList, err := client.Roles.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list roles: ", err)
	}
	fmt.Println("List initial roles", roleList)

	// Define a Role object
	role1 := chef.Role("role1")
	// TODO- add something to the role
	role1.RunList = []string{"pwn"}
	fmt.Println("Define role1", role1)

	// Delete role1 ignoring errors :)
	err = client.Roles.Delete(role1.Name)

	// Create
	roleResult, err := client.Roles.Post(role1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create role role1. ", err)
	}
	fmt.Println("Added role1", roleResult)

	// List roles
	roleList, err = client.Roles.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list roles: ", err)
	}
	fmt.Println("List roles after adding role1", roleList)

	// Create a second time
	roleResult, err = client.Roles.Post(role1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't recreate role role1. ", err)
	}
	fmt.Println("Added role1", roleResult)

	// Read role1 information
	serverRole, err := client.Roles.Get("role1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get role: ", err)
	}
	fmt.Printf("Get role1 %+v\n", serverRole)

	// update role
	// Update the role with something
	updateRole, err := client.Roles.Put(role1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't update role: ", err)
	}
	fmt.Println("Update role1", updateRole)

	// Info after update
	serverRole, err = client.Roles.Get("role1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get role: ", err)
	}
	fmt.Printf("Get role1 after update %+v\n", serverRole)

	// Delete role ignoring errors :)
	err = client.Roles.Delete(role1.Name)
	fmt.Println("Delete role1", err)

	// List roles
	roleList, err = client.Roles.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list roles: ", err)
	}
	fmt.Println("List roles after cleanup", roleList)
}
