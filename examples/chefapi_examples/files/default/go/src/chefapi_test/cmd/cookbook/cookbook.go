//
// Test the go-chef/chef chef server api /organizations/*/cookbooks endpoints against a live chef server
//
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	chef "github.com/go-chef/chef"
	//chef "github.com/MarkGibbons/chefapi"
)


// main Exercise the chef server api
func main() {
	#  values passed in and client are common code, use struct and shared
	user := os.Args[1]
	keyfile := os.Args[2]
	chefurl := os.Args[3]
	ssl

        // Create a client for user access
	client := buildClient(user, keyfile, chefurl)

	// Cookbooks
        fmt.Println("")
	fmt.Println("Starting cookbooks")
	cookbookList := listCookbooks(client) 
	fmt.Println(cookbookList))

	// cookbooks GET
	// List cookbooks
	cookList, err := client.Cookbooks.List()
	if err != nil {
                fmt.Fprintln(os.STDERR, "Issue listing cookbooks:", err)
        }
	// cookbooks/_laters GET
	// cookbooks/_recipes GET
	// cookbooks/NAME GET
	// cookbooks/NAME/VERSION DELETE, GET, PUT

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

// listCookbooks uses the chef server api to list all cookbooks
func listCookbooks(client *chef.Client, filters ...string) map[string]string {
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
