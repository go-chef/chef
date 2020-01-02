//
// Test the go-chef/chef chef authenticate_user api endpoint against a live server
//
package main

import (
	"fmt"
	"os"

	"github.com/go-chef/chef"
	"chefapi_test/testapi"
)


// main Exercise the chef server api
func main() {
        // Create a client for access
	client := testapi.Client()

	// Create a user
        var usr1 chef.User
        usr1 = chef.User{ UserName: "usr1",
                           Email: "user1@domain.io",
                           FirstName: "user1",
                           LastName: "fullname",
                           DisplayName: "User1 Fullname",
                           Password: "Logn12ComplexPwd#",
                   }
        userResult := createUser(client, usr1)

	var ar Authenticate
	// Authenticate with a valid password
	ar.UserName =  "usr1"
	ar.Password =  "Logn12ComplexPwd#"
	err := client.AuthenticateUser.Authenticate(ar)
	fmt.Printf("Authenticate with a valid password %+v", err)

	// Authenticate with an invalid password
	ar.Password =  "Logn12ComplexPwd#"
	err = client.AuthenticateUser.Authenticate(ar)
	fmt.Printf("Authenticate with an invalid password %+v", err)
}

// createUser uses the chef server api to create a single user
func createUser(client *chef.Client, user chef.User) chef.UserResult {
        usrResult, err := client.Users.Create(user)
        if err != nil {
                fmt.Fprintln(os.Stderr, "Issue creating user:", err)
        }
        return usrResult
}
