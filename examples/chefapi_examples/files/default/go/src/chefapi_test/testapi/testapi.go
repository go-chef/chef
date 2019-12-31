//
// Test the go-chef/chef chef server api /organizations endpoints against a live server
//
package testapi

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	chef "github.com/go-chef/chef"
)

// main Exercise the chef server api
func Client() *chef.Client {
	// Pass in the database and chef-server api credentials.
	user := os.Args[1]
	keyfile := os.Args[2]
	chefurl := os.Args[3]
	skipssl, err := strconv.ParseBool(os.Args[4])
	if err != nil {
	   skipssl = true
        }

        // Create a client for access
	return buildClient(user, keyfile, chefurl, skipssl)
}

// buildClient creates a connection to a chef server using the chef api.
func buildClient(user string, keyfile string, baseurl string, skipssl bool) *chef.Client {
	key := clientKey(keyfile)
	client, err := chef.NewClient(&chef.Config{
		Name:    user,
		Key:     string(key),
		BaseURL: baseurl,
		SkipSSL: skipssl,
		// goiardi uses port 4545 by default, chef-zero uses 8889, chef-server uses 443
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue setting up client:", err)
		os.Exit(1)
	}
	return client
}

// clientKey reads the pem file containing the credentials needed to use the chef client.
func clientKey(filepath string) string {
	key, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't read key.pem:", err)
		os.Exit(1)
	}
	return string(key)
}
