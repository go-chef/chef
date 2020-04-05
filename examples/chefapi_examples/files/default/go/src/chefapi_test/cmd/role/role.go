//
// Test the go-chef/chef chef server api /role endpoints against a live server
//
package main

import (
	"encoding/json"
	"fmt"
        "chefapi_test/testapi"
        "github.com/go-chef/chef"
	"os"
)


// main Exercise the chef server api
func main() {
	client := testapi.Client()

	# ? Add a cookbook 
	# ? Add two environments
	# Update the role
	# Get the role
	# Show the role environments  -> array of environment names
	# Show a role environment runlist -> some horrid runlist structure
	# Delete the role
	# ? Delete the environments
	# ? Delete the cookbook

	// Build a stucture to define a role
	defIn := []byte(`{
		"git_repo": "ssh://here.git",
		"users": ["root", "moe"]
	}`)
	ovrIn := []byte(`{
		"env": {
		  "mine": "ample",
		  "yours": "full"
		}
	}`)
	var defAtt interface{}
	var ovrAtt interface{}
	json.Unmarshal(defIn, &defAtt)
	json.Unmarshal(ovrIn, &ovrAtt)
        role1 := chef.Role {
		Name: "role1",
		DefaultAttributes: defAtt,
		Description: "Test role",
		EnvRunLIst:  EnvRunList{
			en1: []string{"recipe[foo1]", "recipe[foo2]"}
			en2: []string{"recipe[foo2]"}
			en3: []string{"recipe[foo3]"}
		},
		RunList: []string{"recipe[foo]", "recipe[baz]", "role[banana]"},
		OverrideAttributes: ovrIn,
	}

	// Add a new role
	roleAdd, err := client.Roles.Create(role1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue adding role1:", err)
	}
	fmt.Println("Added role1", roleAdd)

	// Add role again
	roleAdd, err = client.Roles.Create(role1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue recreating role1:", err)
	}
	fmt.Println("Recreated role1", roleAdd)

	// List roles after adding
	roleList, err := client.Roles.List()
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue printing the existing roles:", err)
	}
	fmt.Printf("List roles after adding role1 %+v\n", roleList)

	// Get new role
	roleOut, err := client.Roles.Get("role1")
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue getting role1:", err)
	}
	fmt.Printf("Get role1 %+v\n", roleOut)

	// Try to get a missing role 
	roleOutMissing, err := client.Roles.Get("nothere")
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue getting nothere:", err)
	}
	fmt.Println("Get nothere", roleOutMissing)

	// Update a role
	role1.Description = "Changed Role"
	// TODO: try changing the runlists, attributes, environments
	roleUpdate, err := client.Roles.Update(role1)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue updating role1:", err)
	}
	fmt.Printf("Update role1 %+v\n", roleUpdate)

	envList, err := GetEnvironments(rolename)
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue listing environments for role1:", err)
	}
	fmt.Printf("Environments for role1 %+v\n", envList)

	envRunList := GetEnvironmentRunlist("role1", "en1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing runlist for role1::en1:", err)
	}
	fmt.Printf("Environments for role1 %+v\n", envList)

        // Clean up
	err = client.Roles.Delete("role1")
	fmt.Println("Delete role1", err)

	// Final list of roles
	roleList, err = client.Roles.List()
	if err != nil {
	       fmt.Fprintln(os.Stderr, "Issue listing the final roles:", err)
	}
	fmt.Printf("List roles after cleanup %+vEndFinalList\n", roleList)
}
