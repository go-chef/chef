//
// Test the go-chef/chef chef server api /users endpoints against a live chef server
//
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	// "github.com/go-chef/chef"
	chef "github.com/MarkGibbons/chefapi"
)


// main Exercise the chef server api
func main() {
	// Pass in the database and chef-server api credentials.
	usr1 := "usr1"
	usr2 := "usr2"
	user := os.Args[1]
	keyfile := os.Args[2]
	chefurl := os.Args[3]

        // Create a client for user access
	client := buildClient(user, keyfile, chefurl)

	// Users
        fmt.Println("")
	fmt.Println("Starting users")
	userList := listUsers(client) // map[string]string  a1px:https://chefp01.nordstrom.net/users/a1px
	fmt.Println(userList)

        fmt.Println("")
	fmt.Println("Added usr1")
	userResult := createUser(client, chef.User{ UserName: usr1, DisplayName: "Show Me", Email: "user1@domain.io", FirstName: "user1", FullName: "All the Name", LastName: "fullname", MiddleName: "J", Password: "Logn12ComplexPwd#" })
	fmt.Println(userResult)

        fmt.Println("")
	fmt.Println("Filter users email=user1")
	userList = listUsers(client, "email=user1@domain.io") // map[string]string  a1px:https://chefp01.nordstrom.net/users/a1px
	fmt.Println(userList)

        fmt.Println("")
	fmt.Println("Verbose users")
	userVerboseOut := listUsersVerbose(client) // []UserVerbose
	fmt.Printf("Verbose Out %v\n", userVerboseOut)

        err := client.AuthenticateUser.Authenticate(chef.AuthenticateUser{ UserName: usr1, Password: "Logn12ComplexPwd#" })
        fmt.Println("")
        fmt.Println("")
        fmt.Println("Error returned from authenticate: ", err)

        fmt.Println("")
	fmt.Println("Added usr1 again")
	userResult = createUser(client, chef.User{ UserName: usr1, Email: "user1@domain.io", FirstName: "user1", LastName: "fullname", DisplayName: "User1 Fullname", Password: "mary" })
	fmt.Println(userResult)

        fmt.Println("")
	fmt.Println("Added usr2")
	userResult = createUser(client, chef.User{ UserName: usr2, Email: "user2@domain.io", FirstName: "User2", LastName: "Fullname", DisplayName: "User2 Fullname", ExternalAuthenticationUid: "mary" })
	fmt.Println(userResult) 

        fmt.Println("")
	fmt.Println("Filter users external_authentication_uid=mary")
	userList = listUsers(client, "external_authentication_uid=mary")
	fmt.Println(userList)

        fmt.Println("")
	fmt.Println("Original usr1")
	userout := getUser(client, usr1)
	fmt.Println(userout)

        fmt.Println("")
	fmt.Println("Original pivotal")
	userout = getUser(client, "pivotal")
	fmt.Println(userout)

        fmt.Println("")
	fmt.Println("After adding usr1 and usr2")
	userList = listUsers(client)
	fmt.Println(userList)

        fmt.Println("")
	fmt.Println("After updating usr1")
        userbody := chef.User{ FullName: "usr1new" }
        fmt.Println("User request", userbody)
	userresult := updateUser(client, "usr1", userbody)
	fmt.Println(userresult)

        fmt.Println("")
	fmt.Println("Get usr1")
	userout = getUser(client, usr1)
	fmt.Println(userout)

        fmt.Println("")
	fmt.Println("Delete usr2")
	userd, userErr := deleteUser(client, usr2)
	fmt.Println(userErr)
	fmt.Println("delete data", userd)

        fmt.Println("")
	fmt.Println("Delete usr1")
	userd, userErr = deleteUser(client, usr1)
	fmt.Println(userErr)
	fmt.Println("delete data", userd)

        fmt.Println("")
	fmt.Println("list after deleting users")
	userList = listUsers(client)
	fmt.Println(userList)

}

// buildClient creates a connection to a chef server using the chef api.
func buildClient(user string, keyfile string, baseurl string) *chef.Client {
	key := clientKey(keyfile)
	client, err := chef.NewClient(&chef.Config{
		Name:    user,
		Key:     string(key),
		BaseURL: baseurl,
		// goiardi is on port 4545 by default, chef-zero is 8889, chef-server is on 443
	})
	if err != nil {
		fmt.Println("Issue setting up client:", err)
		os.Exit(1)
	}
	return client
}

// clientKey reads the pem file containing the credentials needed to use the chef client.
func clientKey(filepath string) string {
	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println("Couldn't read key.pem:", err)
		os.Exit(1)
	}
	return string(key)
}

// createOrganization uses the chef server api to create a single organization
func createOrganization(client *chef.Client, org chef.Organization) chef.OrganizationResult {
	orgResult, err := client.Organizations.Create(org)
	if err != nil {
		fmt.Println("Issue creating org:", err)
	}
	return orgResult
}

// deleteOrganization uses the chef server api to delete a single organization
func deleteOrganization(client *chef.Client, name string) error {
        err := client.Organizations.Delete(name)
        if err != nil {
                fmt.Println("Issue deleting org:", err)
        }
        return err
}

// createUser uses the chef server api to create a single organization
func createUser(client *chef.Client, user chef.User) chef.UserResult {
	usrResult, err := client.Users.Create(user)
	if err != nil {
		fmt.Println("Issue creating user:", err)
	}
	return usrResult
}

// deleteUser uses the chef server api to delete a single organization
func deleteUser(client *chef.Client, name string) (data chef.UserResult, err error) {
	err = client.Users.Delete(name)
	if err != nil {
		fmt.Println("Issue deleting org:", err)
	}
	return
}

// getUser uses the chef server api to get information for a single user
func getUser(client *chef.Client, name string) chef.User {
	userList, err := client.Users.Get(name)
	if err != nil {
		fmt.Println("Issue listing user", err)
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
		fmt.Println("Issue listing users:", err)
	}
	return userList
}

// listUsersVerbose uses the chef server api to list all users and return verbose output
func listUsersVerbose(client *chef.Client) map[string]chef.UsersVerboseItem {
	userList, err := client.Users.ListVerbose()
        fmt.Println("VERBOSE LIST", userList)
	if err != nil {
		fmt.Println("Issue listing verbose users:", err)
	}
	return userList
}

// updateUser uses the chef server api to update information for a single user
func updateUser(client *chef.Client, username string, user chef.User) chef.User {
	user_update, err := client.Users.Update(username, user)
	if err != nil {
		fmt.Println("Issue updating user:", err)
	}
	return user_update
}
