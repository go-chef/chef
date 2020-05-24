package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
	"strings"
)

// client test the chef api
func ApiClient() {
	// Use the default test org
	client := Client()

	// List initial clients
	clientList, err := client.Clients.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list clients: ", err)
	}
	clientfold := strings.Replace(fmt.Sprintf("%+v", clientList), "\n", "", -1)
	fmt.Printf("List initial clients %+v\n", clientfold)

	// Define a Client object
	client1 := chef.ApiNewClient{
		Name:      "client1",
		CreateKey: true,
		Validator: false,
	}
	fmt.Printf("Define client1 %+v\n", client1)

	// Define another Client object
	client2 := chef.ApiNewClient{
		ClientName: "client2",
		CreateKey:  false,
		Validator:  true,
	}
	fmt.Printf("Define client2 %+v\n", client2)

	// Create
	clientResult, err := client.Clients.Create(client1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create client client1. ", err)
	}
	clientresf := strings.Replace(fmt.Sprintf("%+v", clientResult), "\n", "", -1)
	fmt.Printf("Added client1 %+v\n", clientresf)

	clientResult, err = client.Clients.Create(client2)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create client client2. ", err)
	}
	clientresf = strings.Replace(fmt.Sprintf("%+v", clientResult), "\n", "", -1)
	fmt.Printf("Added client2 %+v\n", clientResult)

	// List clients
	clientList, err = client.Clients.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list clients: ", err)
	}
	clientfold = strings.Replace(fmt.Sprintf("%+v", clientList), "\n", "", -1)
	fmt.Printf("List clients after adding %+v\n", clientfold)

	// Create a second time expect 409
	clientResult, err = client.Clients.Create(client1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't recreate client client1. ", err)
	}
	cerr, err := chef.ChefError(err)
	if cerr != nil {
		fmt.Fprintln(os.Stderr, "Couldn't recreate client client1. Code", cerr.StatusCode())
		fmt.Fprintln(os.Stderr, "Couldn't recreate client client1. Msg", cerr.StatusMsg())
		fmt.Fprintln(os.Stderr, "Couldn't recreate client client1. Text", string(cerr.StatusText()))
	}
	fmt.Printf("Added client1 %+v\n", clientResult)

	// Read client1 information
	serverClient, err := client.Clients.Get("client1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get client: ", err)
	}
	fmt.Printf("Get client1 %+v\n", serverClient)

	serverClient, err = client.Clients.Get("client2")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get client2: ", err)
	}
	fmt.Printf("Get client2 %+v\n", serverClient)

	// update client - change the client name
	client1 = chef.ApiNewClient{
		Name:       "clientchanged",
		ClientName: "clientchanged",
		Validator:  true,
	}
	updateClient, err := client.Clients.Update("client1", client1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't update client: ", err)
	}
	fmt.Printf("Update client1 %+v\n", updateClient)

	// update client - change the validator status
	client2 = chef.ApiNewClient{
		Validator: false,
	}
	updateClient, err = client.Clients.Update("client2", client2)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't update client: ", err)
	}
	fmt.Printf("Update client2 %+v\n", updateClient)

	// Info after update
	serverClient, err = client.Clients.Get("clientchanged")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get client: ", err)
	}
	fmt.Printf("Get client1 after update %+v\n", serverClient)

	// Delete clients ignoring errors :)
	err = client.Clients.Delete("clientchanged")
	fmt.Printf("Delete client1 %+v\n", err)
	err = client.Clients.Delete("client2")
	fmt.Printf("Delete client2 %+v\n", err)

	// List clients
	clientList, err = client.Clients.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't list clients: ", err)
	}
	clientfold = strings.Replace(fmt.Sprintf("%+v", clientList), "\n", "", -1)
	fmt.Printf("List clients after cleanup %+v\n", clientfold)
}
