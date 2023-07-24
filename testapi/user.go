// Test the go-chef/chef chef server api /users endpoints against a live chef server
package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// User exercise the chef server api
func User() {
	client := Client()

	// Create a new private key when adding the user
	usr1 := chef.User{UserName: "usr1",
		Email:       "user1@domain.io",
		FirstName:   "user1",
		LastName:    "fullname",
		DisplayName: "User1 Fullname",
		Password:    "Logn12ComplexPwd#",
		CreateKey:   true,
	}

	// Supply a public key
	usr2 := chef.User{UserName: "usr2",
		Email:       "user2@domain.io",
		FirstName:   "user2",
		LastName:    "lastname",
		DisplayName: "User2 Lastname",
		Password:    "Logn12ComplexPwd#",
		PublicKey:   "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoYyN0AIhUh7Fw1+gQtR+ \n0/HY3625IUlVheoUeUz3WnsTrUGSSS4fHvxUiCJlNni1sQvcJ0xC9Bw3iMz7YVFO\nWz5SeKmajqKEnNywN8/NByZhhlLdBxBX/UN04/7aHZMoZxrrjXGLcyjvXN3uxyCO\nyPY989pa68LJ9jXWyyfKjCYdztSFcRuwF7tWgqnlsc8pve/UaWamNOTXQnyrQ6Dp\ndn+1jiNbEJIdxiza7DJMH/9/i/mLIDEFCLRPQ3RqW4T8QrSbkyzPO/iwaHl9U196\n06Ajv1RNnfyHnBXIM+I5mxJRyJCyDFo/MACc5AgO6M0a7sJ/sdX+WccgcHEVbPAl\n1wIDAQAB \n-----END PUBLIC KEY-----\n\n",
	}

	err := deleteUser_u(client, "usr2")
	fmt.Println("Delete usr2", err)

	// Neither PublicKey nor CreateKey specified
	usr3 := chef.User{UserName: "usr3",
		Email:       "user3@domain.io",
		FirstName:   "user3",
		LastName:    "lastname",
		DisplayName: "User3 Lastname",
		Password:    "Logn12ComplexPwd#",
	}

	// Users
	userList := listUsers(client)
	fmt.Printf("List initial users %+v EndInitialList\n", userList)

	userout := getUser(client, "pivotal")
	fmt.Printf("Pivotal user %+v\n", userout)

	userResult := createUser_u(client, usr1)
	fmt.Printf("Add usr1 %+v\n", userResult)

	userResult = createUser_u(client, usr2)
	fmt.Printf("Add usr2 %+v\n", userResult)

	err = deleteUser_u(client, "usr2")
	fmt.Println("Delete usr2", err)

	userResult = createUser_u(client, usr3)
	fmt.Printf("Add usr3 %+v\n", userResult)

	err = deleteUser_u(client, "usr3")
	fmt.Println("Delete usr3", err)

	userList = listUsers(client, "email=user1@domain.io")
	fmt.Printf("Filter users %+v\n", userList)

	userVerboseOut := listUsersVerbose(client)
	fmt.Printf("Verbose out %v\n", userVerboseOut)

	userResult = createUser_u(client, usr1)
	fmt.Printf("Add usr1 again %+v\n", userResult)

	userout = getUser(client, "usr1")
	fmt.Printf("Get usr1 %+v\n", userout)

	userList = listUsers(client)
	fmt.Printf("List after adding %+v EndAddList\n", userList)

	userbody := chef.User{UserName: "usr1", DisplayName: "usr1", Email: "myuser@samp.com"}
	userresult := updateUser(client, "usr1", userbody)
	fmt.Printf("Update usr1 partial update %+v\n", userresult)

	userout = getUser(client, "usr1")
	fmt.Printf("Get usr1 after partial update %+v\n", userout)

	userbody = chef.User{UserName: "usr1", DisplayName: "usr1", FirstName: "user", MiddleName: "mid", LastName: "name", Email: "myuser@samp.com"}
	userresult = updateUser(client, "usr1", userbody)
	fmt.Printf("Update usr1 full update %+v\n", userresult)

	userout = getUser(client, "usr1")
	fmt.Printf("Get usr1 after full update %+v\n", userout)

	userbody = chef.User{UserName: "usr1new", DisplayName: "usr1", FirstName: "user", MiddleName: "mid", LastName: "name", Email: "myuser@samp.com"}
	userresult = updateUser(client, "usr1", userbody)
	fmt.Printf("Update usr1 rename %+v\n", userresult)

	userout = getUser(client, "usr1new")
	fmt.Printf("Get usr1 after rename %+v\n", userout)

	userbody = chef.User{UserName: "usr1new", DisplayName: "usr1", FirstName: "user", MiddleName: "mid", LastName: "name", Email: "myuser@samp.com", CreateKey: true}
	userresult = updateUser(client, "usr1new", userbody)
	fmt.Printf("Update usr1new create key  %+v\n", userresult)

	userout = getUser(client, "usr1new")
	fmt.Printf("Get usr1new after create key %+v\n", userout)

	err = deleteUser_u(client, "usr1")
	fmt.Println("Delete usr1", err)

	err = deleteUser_u(client, "usr1new")
	fmt.Println("Delete usr1new", err)

	userList = listUsers(client)
	fmt.Printf("List after cleanup %+v EndCleanupList\n", userList)
}

// createUser_u uses the chef server api to create a single user
func createUser_u(client *chef.Client, user chef.User) chef.UserResult {
	usrResult, err := client.Users.Create(user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue creating user %+v\n", err)
		if cerr, ok := err.(*chef.ErrorResponse); ok {
			fmt.Fprintf(os.Stderr, "Issue creating user err: %+v\n", cerr.Error())
			fmt.Fprintf(os.Stderr, "Issue creating user code: %+v\n", cerr.StatusCode())
			fmt.Fprintf(os.Stderr, "Issue creating user method: %+v\n", cerr.StatusMethod())
			fmt.Fprintf(os.Stderr, "Issue creating user url: %+v\n", cerr.StatusURL().String())
		}
	}
	return usrResult
}

// deleteUser_u uses the chef server api to delete a single user
func deleteUser_u(client *chef.Client, name string) (err error) {
	err = client.Users.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting user:", err)
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
func listUsersVerbose(client *chef.Client) map[string]chef.UserVerboseResult {
	userList, err := client.Users.VerboseList()
	fmt.Printf("VERBOSE LIST %+v\n", userList)
	if err != nil {
		fmt.Println("Issue listing verbose users:", err)
	}
	return userList
}

// updateUser uses the chef server api to update information for a single user
func updateUser(client *chef.Client, username string, user chef.User) chef.UserResult {
	user_update, err := client.Users.Update(username, user)
	if err != nil {
		fmt.Println("Issue updating user:", err)
	}
	return user_update
}
