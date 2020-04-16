package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// principle test the chef api
func PrincipalsAdd() {
	// Use the default test org
	client := Client()

	// Create a User
        usr1 := chef.User{UserName: "usr1",
                Email:       "user1@domain.io",
                FirstName:   "user1",
                LastName:    "fullname",
                DisplayName: "User1 Fullname",
                Password:    "Logn12ComplexPwd#",
                CreateKey:   true,
        }
        _ = createUser_u(client, usr1)

	// Create a User with the same name as a client
	client1 := chef.User{UserName: "client1",
                Email:       "client@domain.io",
                FirstName:   "user1",
                LastName:    "fullname",
                DisplayName: "User1 Fullname",
                Password:    "Logn12ComplexPwd#",
                CreateKey:   true,
        }
	_ = createUser_u(client, client1)
}

// createUser_p uses the chef server api to create a single user
func createUser_p(client *chef.Client, user chef.User) chef.UserResult {
        usrResult, err := client.Users.Create(user)
        if err != nil {
                fmt.Fprintln(os.Stderr, "Issue creating user:", err)
        }
        return usrResult
}
