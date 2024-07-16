// Test the go-chef/chef chef server api /users/USERNAME/keys endpoints against a live chef server
package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
	"strings"
)

// userkey exercise the chef server api
func Userkey() {
	client := Client(nil)

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

	// Neither PublicKey nor CreateKey specified
	usr3 := chef.User{UserName: "usr3",
		Email:       "user3@domain.io",
		FirstName:   "user3",
		LastName:    "lastname",
		DisplayName: "User3 Lastname",
		Password:    "Logn12ComplexPwd#",
	}

	_ = createUser_key(client, usr1)
	fmt.Printf("Add usr1\n")
	_ = createUser_key(client, usr2)
	fmt.Printf("Add usr2\n")
	_ = createUser_key(client, usr3)
	fmt.Printf("Add usr3\n")

	// User Keys
	userkeys := listUserKeys(client, "usr1")
	fmt.Printf("List initial user usr1 keys %+v\n", userkeys)

	userkeys = listUserKeys(client, "usr2")
	fmt.Printf("List initial user usr2 keys %+v\n", userkeys)

	userkeys = listUserKeys(client, "usr3")
	fmt.Printf("List initial user usr3 keys %+v\n", userkeys)

	// Add a key to a user
	keyadd := chef.AccessKey{
		Name:           "newkey",
		PublicKey:      "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoYyN0AIhUh7Fw1+gQtR+ \n0/HY3625IUlVheoUeUz3WnsTrUGSSS4fHvxUiCJlNni1sQvcJ0xC9Bw3iMz7YVFO\nWz5SeKmajqKEnNywN8/NByZhhlLdBxBX/UN04/7aHZMoZxrrjXGLcyjvXN3uxyCO\nyPY989pa68LJ9jXWyyfKjCYdztSFcRuwF7tWgqnlsc8pve/UaWamNOTXQnyrQ6Dp\ndn+1jiNbEJIdxiza7DJMH/9/i/mLIDEFCLRPQ3RqW4T8QrSbkyzPO/iwaHl9U196\n06Ajv1RNnfyHnBXIM+I5mxJRyJCyDFo/MACc5AgO6M0a7sJ/sdX+WccgcHEVbPAl\n1wIDAQAB \n-----END PUBLIC KEY-----\n\n",
		ExpirationDate: "infinity",
	}
	keyout, err := addUserKey(client, "usr1", keyadd)
	fmt.Printf("Add usr1 key %+v\n", keyout)
	// List the user keys after adding
	userkeys = listUserKeys(client, "usr1")
	fmt.Printf("List after add usr1 keys %+v\n", userkeys)

	// Add a defaultkey to user usr3
	keyadd.Name = "default"
	keyout, err = addUserKey(client, "usr3", keyadd)
	fmt.Printf("Add usr3 key %+v\n", keyout)
	// List the user keys after adding
	userkeys = listUserKeys(client, "usr3")
	fmt.Printf("List after add usr3 keys %+v\n", userkeys)

	// Get key detail
	keydetail, err := client.Users.GetKey("usr1", "default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error displaying key detail %+v\n", err)
	}
	keyfold := strings.Replace(fmt.Sprintf("%+v", keydetail), "\n", "", -1)
	fmt.Printf("Key detail usr1 default %+v\n", keyfold)

	// update a key
	keyadd.Name = "default"
	keyupdate, err := client.Users.UpdateKey("usr1", "default", keyadd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating usr1 default key%+v\n", err)
	}
	keyfold = strings.Replace(fmt.Sprintf("%+v", keyupdate), "\n", "", -1)
	fmt.Printf("Key update output usr1 default %+v\n", keyfold)
	// Get key detail after update
	keydetail, err = client.Users.GetKey("usr1", "default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error displaying key detail %+v\n", err)
	}
	keyfold = strings.Replace(fmt.Sprintf("%+v", keydetail), "\n", "", -1)
	fmt.Printf("Updated key detail usr1 default %+v\n", keyfold)

	// delete the key
	keydel, err := client.Users.DeleteKey("usr1", "default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting key %+v\n", err)
	}
	keyfold = strings.Replace(fmt.Sprintf("%+v", keydel), "\n", "", -1)
	fmt.Printf("List delete result usr1 keys %+v\n", keyfold)
	// list the key after delete - expect 404
	keydetail, err = client.Users.GetKey("usr1", "default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error displaying key detail %+v\n", err)
	}
	fmt.Printf("Deleted key detail usr1 default %+v\n", keydetail)

	// Delete the users
	err = deleteUser_key(client, "usr1")
	fmt.Printf("Delete usr1 %+v\n", err)
	err = deleteUser_key(client, "usr2")
	fmt.Printf("Delete usr2 %+v\n", err)
	err = deleteUser_key(client, "usr3")
	fmt.Printf("Delete usr3 %+v\n", err)

}

// listUserKeys uses the chef server api to show the keys for a user
func listUserKeys(client *chef.Client, name string) (userkeys []chef.KeyItem) {
	userkeys, err := client.Users.ListKeys(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue showing keys for user %s: %+v\n", name, err)
	}
	return userkeys
}

// addUserKey uses the chef server api to add a key to user
func addUserKey(client *chef.Client, name string, keyadd chef.AccessKey) (userkey chef.KeyItem, err error) {
	userkey, err = client.Users.AddKey(name, keyadd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", err)
	}
	return
}

// createUser_key uses the chef server api to create a single user
func createUser_key(client *chef.Client, user chef.User) chef.UserResult {
	usrResult, err := client.Users.Create(user)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating user:", err)
	}
	return usrResult
}

// deleteUser_key uses the chef server api to delete a single user
func deleteUser_key(client *chef.Client, name string) (err error) {
	err = client.Users.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", err)
	}
	return
}
