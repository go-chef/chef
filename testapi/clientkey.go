// Test the go-chef/chef chef server api /clients/CLIENTNAME/keys endpoints against a live chef server
package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
	"strings"
)

// clientkey exercise the chef server api
func Clientkey() {
	client := Client(nil)

	// Create a new private key when adding the client
	clnt1 := chef.ApiNewClient{
		ClientName: "clnt1",
		CreateKey:  true,
	}
	// CreateKey not specified
	clnt3 := chef.ApiNewClient{
		ClientName: "clnt3",
	}

	_ = createClient(client, clnt1)
	fmt.Printf("Add clnt1\n")
	_ = createClient(client, clnt3)
	fmt.Printf("Add clnt3\n")

	// Client Keys
	clientkeys := listClientKeys(client, "clnt1")
	fmt.Printf("List initial client clnt1 keys %+v\n", clientkeys)
	clientkeys = listClientKeys(client, "clnt3")
	fmt.Printf("List initial client clnt3 keys %+v\n", clientkeys)

	// Add a key to a client
	keyadd := chef.AccessKey{
		Name:           "newkey",
		PublicKey:      "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoYyN0AIhUh7Fw1+gQtR+ \n0/HY3625IUlVheoUeUz3WnsTrUGSSS4fHvxUiCJlNni1sQvcJ0xC9Bw3iMz7YVFO\nWz5SeKmajqKEnNywN8/NByZhhlLdBxBX/UN04/7aHZMoZxrrjXGLcyjvXN3uxyCO\nyPY989pa68LJ9jXWyyfKjCYdztSFcRuwF7tWgqnlsc8pve/UaWamNOTXQnyrQ6Dp\ndn+1jiNbEJIdxiza7DJMH/9/i/mLIDEFCLRPQ3RqW4T8QrSbkyzPO/iwaHl9U196\n06Ajv1RNnfyHnBXIM+I5mxJRyJCyDFo/MACc5AgO6M0a7sJ/sdX+WccgcHEVbPAl\n1wIDAQAB \n-----END PUBLIC KEY-----\n\n",
		ExpirationDate: "infinity",
	}
	keyout, err := addClientKey(client, "clnt1", keyadd)
	fmt.Printf("Add clnt1 key %+v\n", keyout)
	// List the client keys after adding
	clientkeys = listClientKeys(client, "clnt1")
	fmt.Printf("List after add clnt1 keys %+v\n", clientkeys)

	// Add a defaultkey to client clnt3
	keyadd.Name = "default"
	keyout, err = addClientKey(client, "clnt3", keyadd)
	fmt.Printf("Add clnt3 key %+v\n", keyout)
	// List the client keys after adding
	clientkeys = listClientKeys(client, "clnt3")
	fmt.Printf("List after add clnt3 keys %+v\n", clientkeys)

	// Get key detail
	keydetail, err := client.Clients.GetKey("clnt1", "default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error displaying key detail %+v\n", err)
	}
	keyfold := strings.Replace(fmt.Sprintf("%+v", keydetail), "\n", "", -1)
	fmt.Printf("Key detail clnt1 default %+v\n", keyfold)

	// update a key
	keyadd.Name = "default"
	keyupdate, err := client.Clients.UpdateKey("clnt1", "default", keyadd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating clnt1 default key%+v\n", err)
	}
	keyfold = strings.Replace(fmt.Sprintf("%+v", keyupdate), "\n", "", -1)
	fmt.Printf("Key update output clnt1 default %+v\n", keyfold)
	// Get key detail after update
	keydetail, err = client.Clients.GetKey("clnt1", "default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error displaying key detail %+v\n", err)
	}
	keyfold = strings.Replace(fmt.Sprintf("%+v", keydetail), "\n", "", -1)
	fmt.Printf("Updated key detail clnt1 default %+v\n", keyfold)

	// delete the key
	keydel, err := client.Clients.DeleteKey("clnt1", "default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting key %+v\n", err)
	}
	keyfold = strings.Replace(fmt.Sprintf("%+v", keydel), "\n", "", -1)
	fmt.Printf("List delete result clnt1 keys %+v\n", keyfold)
	// list the key after delete - expect 404
	keydetail, err = client.Clients.GetKey("clnt1", "default")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error displaying key detail %+v\n", err)
	}
	fmt.Printf("Deleted key detail clnt1 default %+v\n", keydetail)

	// Delete the clients
	err = deleteClient(client, "clnt1")
	fmt.Printf("Delete clnt1 %+v\n", err)
	err = deleteClient(client, "clnt3")
	fmt.Printf("Delete clnt3 %+v\n", err)

}

// listClientKeys uses the chef server api to show the keys for a client
func listClientKeys(client *chef.Client, name string) (clientkeys []chef.KeyItem) {
	clientkeys, err := client.Clients.ListKeys(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue showing keys for client %s: %+v\n", name, err)
	}
	return clientkeys
}

// addClientKey uses the chef server api to add a key to client
func addClientKey(client *chef.Client, name string, keyadd chef.AccessKey) (clientkey chef.KeyItem, err error) {
	clientkey, err = client.Clients.AddKey(name, keyadd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", err)
	}
	return
}

// createClient uses the chef server api to create a single client
func createClient(client *chef.Client, apiclient chef.ApiNewClient) *chef.ApiClientCreateResult {
	clntResult, err := client.Clients.Create(apiclient)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating client:", err)
	}
	return clntResult
}

// deleteClient uses the chef server api to delete a single client
func deleteClient(client *chef.Client, name string) (err error) {
	err = client.Clients.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", err)
	}
	return
}
