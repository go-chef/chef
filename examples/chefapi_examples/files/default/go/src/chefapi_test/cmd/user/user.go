//
// Test the go-chef/chef chef server api /users endpoints against a live chef server
//
package main

import (
	"fmt"
	"github.com/go-chef/chef"
	"chefapi_test/testapi"
	"os"
)


// main Exercise the chef server api
func main() {
        client := testapi.Client()

	var usr1 chef.User
	usr1 = chef.User{ UserName: "usr1",
	                   Email: "user1@domain.io",
			   FirstName: "user1",
			   LastName: "fullname",
			   DisplayName: "User1 Fullname",
			   Password: "Logn12ComplexPwd#",
		   }

	// Users
	userList := listUsers(client)
	fmt.Printf("List initial users %+v EndInitialList\n", userList)

	userResult := createUser(client, usr1)
	fmt.Println("Add usr1", userResult)

	userList = listUsers(client, "email=user1@domain.io")
	fmt.Printf("Filter users %+v\n", userList)

	// userVerboseOut := listUsersVerbose(client)
	// fmt.Printf("Verbose out %v\n", userVerboseOut)

        // err := client.AuthenticateUser.Authenticate(chef.AuthenticateUser{ UserName: usr1, Password: "Logn12ComplexPwd#" })
        // fmt.Println("Authenticate usr1 ", err)

	userResult = createUser(client, usr1)
	fmt.Printf("Add user1 again %+v\n", userResult)

	userout := getUser(client, "usr1")
	fmt.Printf("Get usr1 %+v\n", userout)

	userout = getUser(client, "pivotal")
	fmt.Printf("Pivotal user %+v\n", userout)

	userList = listUsers(client)
	fmt.Printf("List after adding %+v EndAddList\n", userList)

        // userbody := chef.User{ FullName: "usr1new" }
	// userresult := updateUser(client, "usr1", userbody)
	// fmt.Printf("Update user1 %+v", userresult)

	// userout = getUser(client, "usr1")
	// fmt.Println("Get usr1 after update %+v\n", userout)

	err := deleteUser(client, "usr1")
	fmt.Println("Delete usr1", err)

	userList = listUsers(client)
	fmt.Printf("List after cleanup %+v EndCleanupList\n", userList)
}

// createUser uses the chef server api to create a single organization
func createUser(client *chef.Client, user chef.User) chef.UserResult {
	usrResult, err := client.Users.Create(user)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating user:", err)
	}
	return usrResult
}

// deleteUser uses the chef server api to delete a single organization
func deleteUser(client *chef.Client, name string) (err error) {
	err = client.Users.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", err)
	}
	return
}

// getUser uses the chef server api to get information for a single user
func getUser(client *chef.Client, name string) chef.User {
	userList, err := client.Users.Get(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing user", err)
	}
	return userList
}

// listUsers uses the chef server api to list all users
func listUsers(client *chef.Client, filters ...string) map[string]string {
        var filter string
        if len(filters) > 0 {
        	filter = filters[0]
        }
	userList, err := client.Users.List(filter)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing users:", err)
	}
	return userList
}

// listUsersVerbose uses the chef server api to list all users and return verbose output
//func listUsersVerbose(client *chef.Client) map[string]chef.UsersVerboseItem {
//	userList, err := client.Users.ListVerbose()
 //       fmt.Println("VERBOSE LIST", userList)
//	if err != nil {
//		fmt.Println("Issue listing verbose users:", err)
//	}
//	return userList
//}

// updateUser uses the chef server api to update information for a single user
// func updateUser(client *chef.Client, username string, user chef.User) chef.User {
//	user_update, err := client.Users.Update(username, user)
	//if err != nil {
	//i	fmt.Println("Issue updating user:", err)
	//
//}
//	return user_update
//}
