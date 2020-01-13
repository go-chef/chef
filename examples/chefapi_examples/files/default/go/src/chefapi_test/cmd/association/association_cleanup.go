//
// Test the go-chef/chef chef server api /organization/:org/user and /organization/:org/association_requests
// endpoints against a live server
//
package main

import (
	"fmt"
        "chefapi_test/testapi"
        "github.com/go-chef/chef"
	"os"
)


// main Exercise the chef server api
func main() {
	client := testapi.Client()
        deleteUser(client, "usrinvite")
        deleteUser(client, "usr2invite")
        deleteUser(client, "usradd")

}

 // deleteUser uses the chef server api to delete a single user
 func deleteUser(client *chef.Client, name string) (err error) {
         err = client.Users.Delete(name)
         if err != nil {
                 fmt.Fprintln(os.Stderr, "Issue deleting user:", err)
         }
         return
 }
