package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// principle test the chef api
func PrincipalsDel() {
	// Use the default test org
	client := Client(nil)

	_ = deleteUser_p(client, "client1")
	_ = deleteUser_p(client, "usr1")
}

// deleteUser_p uses the chef server api to delete a single user
func deleteUser_p(client *chef.Client, name string) (err error) {
	err = client.Users.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting user:", err)
	}
	return
}
