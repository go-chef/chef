package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
	"strings"
)

// principle test the chef api
func Principals() {
	// Use the default test org
	client := Client()

	// Create a client
	client1 := chef.ApiNewClient{
		Name:      "client1",
		CreateKey: true,
		Validator: false,
	}
	_, err := client.Clients.Create(client1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't create client client1. ", err)
	}

	// Get principals - public key information
	clientInfo, err := client.Principals.Get("client1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get client1 keys: ", err)
	}
	clientinfof := strings.Replace(fmt.Sprintf("%+v", clientInfo), "\n", "", -1)
	fmt.Printf("Client principal %+v\n", clientinfof)

	userInfo, err := client.Principals.Get("usr1")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't get usr1 keys: ", err)
	}
	userinfof := strings.Replace(fmt.Sprintf("%+v", userInfo), "\n", "", -1)
	fmt.Printf("User principal %+v\n", userinfof)

	// Delete clients
	err = client.Clients.Delete("client1")
	fmt.Printf("Delete client1 %+v\n", err)
}
