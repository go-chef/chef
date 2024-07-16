// Test the go-chef/chef chef authenticate_user api endpoint against a live server
package testapi

import (
	"fmt"
	"os"

	"github.com/go-chef/chef"
)

// authenticate exercise the chef server api
func Authenticate() {
	// Create a client for access
	client := Client(nil)

	// Create a user
	var usr chef.User
	usr = chef.User{UserName: "usrauth",
		Email:       "usrauth@domain.io",
		FirstName:   "usrauth",
		LastName:    "fullname",
		DisplayName: "Userauth Fullname",
		Password:    "Logn12ComplexPwd#",
	}
	createUser_auth(client, usr)

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
	deleteUser_auth(client, "usrauth")
}

// createUser_auth uses the chef server api to create a single user
func createUser_auth(client *chef.Client, user chef.User) chef.UserResult {
	usrResult, err := client.Users.Create(user)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating user:", err)
	}
	return usrResult
}

// deleteUser_auth uses the chef server api to delete a single user
func deleteUser_auth(client *chef.Client, name string) (err error) {
	err = client.Users.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", err)
	}
	return
}
