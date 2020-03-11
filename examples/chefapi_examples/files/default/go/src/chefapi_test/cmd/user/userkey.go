//
// Test the go-chef/chef chef server api /users/USERNAME/keys endpoints against a live chef server
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

	// Create a new private key when adding the user
	usr1 := chef.User{ UserName: "usr1",
	                   Email: "user1@domain.io",
			   FirstName: "user1",
			   LastName: "fullname",
			   DisplayName: "User1 Fullname",
			   Password: "Logn12ComplexPwd#",
			   CreateKey: true,
		   }

        // Supply a public key
        usr2 := chef.User{ UserName: "usr2",
	                   Email: "user2@domain.io",
			   FirstName: "user2",
			   LastName: "lastname",
			   DisplayName: "User2 Lastname",
			   Password: "Logn12ComplexPwd#",
			   PublicKey: "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoYyN0AIhUh7Fw1+gQtR+ \n0/HY3625IUlVheoUeUz3WnsTrUGSSS4fHvxUiCJlNni1sQvcJ0xC9Bw3iMz7YVFO\nWz5SeKmajqKEnNywN8/NByZhhlLdBxBX/UN04/7aHZMoZxrrjXGLcyjvXN3uxyCO\nyPY989pa68LJ9jXWyyfKjCYdztSFcRuwF7tWgqnlsc8pve/UaWamNOTXQnyrQ6Dp\ndn+1jiNbEJIdxiza7DJMH/9/i/mLIDEFCLRPQ3RqW4T8QrSbkyzPO/iwaHl9U196\n06Ajv1RNnfyHnBXIM+I5mxJRyJCyDFo/MACc5AgO6M0a7sJ/sdX+WccgcHEVbPAl\n1wIDAQAB \n-----END PUBLIC KEY-----\n\n",
		   }

	err := deleteUser(client, "usr2")
	fmt.Println("Delete usr2", err)

        // Neither PublicKey nor CreateKey specified
        usr3 := chef.User{ UserName: "usr3",
	                   Email: "user3@domain.io",
			   FirstName: "user3",
			   LastName: "lastname",
			   DisplayName: "User3 Lastname",
			   Password: "Logn12ComplexPwd#",
		   }

	// User Keys
	userkeys := listUserKeys(client, "usr1")
	fmt.Printf("List initial user usr1 keys %+v EndInitialList\n", userkeys)

	userkeys := listUserKeys(client, "usr2")
	fmt.Printf("List initial user usr2 keys %+v EndInitialList\n", userkeys)

	userkeys := listUserKeys(client, "usr3")
	fmt.Printf("List initial user usr3 keys %+v EndInitialList\n", userkeys)

	// Add a key to a user
	// List the user after adding
	// Get key detail
	// update a key
	// Get key detail after update
	// delete the key 
	// list the user keys after deleting

}

// listUserKeys uses the chef server api to show the keys for a user
func listUserKeys(client *chef.Client, user chef.User) (userkeys []chef.UserKeyItem) {
	usrResult, err := client.Users.ListUserKeys(user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue showing keys for user %s: %+v\n", user, err)
	}
	return userkeys
}

// deleteUser uses the chef server api to delete a single user
func deleteUser(client *chef.Client, name string) (err error) {
	err = client.Users.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", err)
	}
	return
}

// getUserKey uses the chef server api to get information for a single user
func getUserKey(client *chef.Client, name string) chef.User {
	userList, err := client.Users.Get(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing user", err)
	}
	return userList
}
