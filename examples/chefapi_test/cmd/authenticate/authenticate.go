//
// Test the go-chef/chef chef authenticate_user api endpoint against a live server
//
package main

import (
	"fmt"
	"os"

	"chefapi_test/testapi"
	"github.com/go-chef/chef"
)

// main Exercise the chef server api
func main() {
	// Create a client for access
	client := testapi.Client()

	// Create a user
	var usr chef.User
	usr = chef.User{UserName: "usrauth",
		Email:       "usrauth@domain.io",
		FirstName:   "usrauth",
		LastName:    "fullname",
		DisplayName: "Userauth Fullname",
		Password:    "Logn12ComplexPwd#",
	}
	createUser(client, usr)

	var ar chef.Authenticate
	// Authenticate with a valid password
	ar.UserName = "usrauth"
	ar.Password = "Logn12ComplexPwd#"
	err := client.AuthenticateUser.Authenticate(ar)
	fmt.Printf("Authenticate with a valid password %+vauthenticate\n", err)

	// Authenticate with an invalid password
	ar.Password = "badpassword"
	err = client.AuthenticateUser.Authenticate(ar)
	fmt.Printf("Authenticate with an invalid password %+v\n", err)

	// Cleanup
	deleteUser(client, "usrauth")
}

// createUser uses the chef server api to create a single user
func createUser(client *chef.Client, user chef.User) chef.UserResult {
	usrResult, err := client.Users.Create(user)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating user:", err)
	}
	return usrResult
}

// deleteUser uses the chef server api to delete a single user
func deleteUser(client *chef.Client, name string) (err error) {
	err = client.Users.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", err)
	}
	return
}
