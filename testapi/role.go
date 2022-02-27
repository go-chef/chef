//
// Test the go-chef/chef chef server api /role endpoints against a live server
//
package testapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// role exercise the chef server api
func Role() {
	client := Client()

	// The environment need to exist for the GetEnvironmentRunlist function to work
	create_en1(client)

	// Build a stucture to define a role
	role1 := create_role1()

	// Add a new role
	roleAdd, err := client.Roles.Create(&role1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue adding role1:", err)
	}
	fmt.Printf("Added role1 %+v\n", roleAdd)

	// Build a stucture to define a role
	roleNR := create_role_norunlist()
	// Add a new role
	roleAdd, err = client.Roles.Create(&roleNR)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue adding roleNR:", err)
	}
	fmt.Printf("Added roleNR %+v\n", roleAdd)

	// Add role again
	roleAdd, err = client.Roles.Create(&role1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue recreating role1:", err)
	}
	cerr, err := chef.ChefError(err)
	if cerr != nil {
		fmt.Fprintln(os.Stderr, "Issue recreating role1:", cerr.StatusCode())
	}
	fmt.Printf("Recreated role1 %+v\n", roleAdd)

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
	fmt.Printf("Get nothere %+v\n", roleOutMissing)

	// Update a role
	role1.Description = "Changed Role"
	// TODO: try changing the runlists, attributes, environment run list
	roleUpdate, err := client.Roles.Put(&role1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating role1:", err)
	}
	fmt.Printf("Update role1 %+v\n", roleUpdate)

	envList, err := client.Roles.GetEnvironments("role1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing environments for role1:", err)
	}
	fmt.Printf("Environments for role1 %+v\n", envList)

	envRunList, err := client.Roles.GetEnvironmentRunlist("role1", "en1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing runlist for role1::en1:", err)
	}
	fmt.Printf("Environments for role1 %+v\n", envRunList)

	// Clean up
	err = client.Roles.Delete("role1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting role1", err)
	}
	fmt.Printf("Delete role1 %+v\n", err)

	err = client.Roles.Delete("roleNR")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting roleNR", err)
	}
	fmt.Printf("Delete roleNR %+v\n", err)

	_, err = client.Environments.Delete("en1")

	// Final list of roles
	roleList, err = client.Roles.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing the final roles:", err)
	}
	fmt.Printf("List roles after cleanup %+v\n", roleList)
}

func create_en1(client *chef.Client) {
	en1 := chef.Environment{
		Name:        "en1",
		Description: "Test environment",
		CookbookVersions: map[string]string{
			"a": "0.0.0",
		},
	}
	_, err := client.Environments.Create(&en1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue adding en1:", err)
	}
	return
}

func create_role1() chef.Role {
	defIn := []byte(`{
		"git_repo": "here.git",
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
	role1 := chef.Role{
		Name:              "role1",
		DefaultAttributes: defAtt,
		Description:       "Test role",
		EnvRunList: chef.EnvRunList{
			"en1": []string{"recipe[foo1]", "recipe[foo2]"},
			"en2": []string{"recipe[foo2]"},
		},
		RunList:            []string{"recipe[foo]", "recipe[baz]", "role[banana]"},
		OverrideAttributes: ovrAtt,
	}
	return role1
}

func create_role_norunlist() chef.Role {
	role1 := chef.Role{
		Name:              "roleNR",
		Description:       "Test role",
		RunList:            []string{"recipe[foo]", "recipe[baz]", "role[banana]"},
	}
	return role1
}
