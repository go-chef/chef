
package main

import (
	"fmt"
	"os"
	"testapi"
)

var cases = map[string]func(){
	"association": testapi.Association,
	"association_cleanup": testapi.AssociationCleanup,
	"association_setup": testapi.AssociationSetup,
	"authenticate": testapi.Authenticate,
	"client": testapi.ApiClient,
	"clientkey": testapi.Clientkey,
	"cookbook": testapi.Cookbook,
	"databag": testapi.Databag,
	// TODO: fix environment and sandbox
	"environment": testapi.Environment,
	"group": testapi.Group,
	"license": testapi.License,
	"node": testapi.Node,
	"organization": testapi.Organization,
	"role": testapi.Role,
	"sandbox": testapi.Sandbox,
	"search": testapi.Search,
	"search_pagination": testapi.SearchPagination,
	"status": testapi.Status,
	"universe": testapi.Universe,
	"user": testapi.User,
	"userkey": testapi.Userkey,
}

// Invoke the requested testapi test function
func main() {
	testcase := os.Args[1]
	fn, ok := cases[testcase]
	if ok {
		fn()
        } else {
		fmt.Fprintf(os.Stderr, "Requested case %+s was not found\n", testcase)
	}
}

