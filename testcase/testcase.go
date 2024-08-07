package main

import (
	"fmt"
	"github.com/go-chef/chef/testapi"
	"os"
)

var cases = map[string]func(){
	"acl":                 testapi.ACL,
	"association":         testapi.Association,
	"association_cleanup": testapi.AssociationCleanup,
	"association_setup":   testapi.AssociationSetup,
	"authenticate":        testapi.Authenticate,
	"client":              testapi.ApiClient,
	"clientkey":           testapi.Clientkey,
	"container":           testapi.Container,
	"cookbook":            testapi.Cookbook,
	"databag":             testapi.Databag,
	"environment":         testapi.Environment,
	"group":               testapi.Group,
	"http":                testapi.Http,
	"license":             testapi.License,
	"node":                testapi.Node,
	"organization":        testapi.Organization,
	"policy":              testapi.Policy,
	"policygroup":         testapi.PolicyGroup,
	"principals":          testapi.Principals,
	"principals_add":      testapi.PrincipalsAdd,
	"principals_del":      testapi.PrincipalsDel,
	"role":                testapi.Role,
	"required_recipe":     testapi.RequiredRecipe,
	// TODO: fix sandbox
	"sandbox":           testapi.Sandbox,
	"search":            testapi.Search,
	"search_pagination": testapi.SearchPagination,
	"stats":             testapi.Stats,
	"status":            testapi.Status,
	"universe":          testapi.Universe,
	"user":              testapi.User,
	"userkey":           testapi.Userkey,
}

// Invoke the requested testapi test function
func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: specify the test case name arg length %v  values %+v\n", len(os.Args), os.Args)
	} else {
		testcase := os.Args[1]
		fn, ok := cases[testcase]
		if ok {
			fn()
		} else {
			fmt.Fprintf(os.Stderr, "Requested case %+s was not found\n", testcase)
		}
	}
}
