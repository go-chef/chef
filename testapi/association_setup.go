// Test the go-chef/chef chef server api /organization/:org/user and /organization/:org/association_requests
// endpoints against a live server
package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// association_setup exercise the chef server api
func AssociationSetup() {
	client := Client(nil)

	// Create a user
	var usr chef.User
	usr = chef.User{UserName: "usrinvite",
		Email:       "usrauth@domain.io",
		FirstName:   "usr",
		LastName:    "invite",
		DisplayName: "Userauth Fullname",
		Password:    "Logn12ComplexPwd#",
	}
	createUser(client, usr)

	usr = chef.User{UserName: "usr2invite",
		Email:       "usr22auth@domain.io",
		FirstName:   "usr22",
		LastName:    "invite",
		DisplayName: "User22auth Fullname",
		Password:    "Logn12ComplexPwd#",
	}
	createUser(client, usr)

	usr = chef.User{UserName: "usradd",
		Email:       "usradd@domain.io",
		FirstName:   "usr",
		LastName:    "add",
		DisplayName: "UserAdd Fullname",
		Password:    "Logn12ComplexPwd#",
	}
	createUser(client, usr)

}

// createUser uses the chef server api to create a single user
func createUser(client *chef.Client, user chef.User) chef.UserResult {
	usrResult, err := client.Users.Create(user)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating user:", err)
	}
	return usrResult
}
