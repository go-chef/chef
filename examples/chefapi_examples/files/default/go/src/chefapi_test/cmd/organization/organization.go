//
// Test the go-chef/chef chef server api /organizations endpoints against a live server
//
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-chef/chef"
	"chefapi_test/testapi"
)


// main Exercise the chef server api
func main() {
        // Create a client for access
	client := testapi.Client()

	// Organization tests
	org1 := "org1"
	org2 := "org2"

	orgList := listOrganizations(client)
	fmt.Println("List initial organizations", orgList)

	orgResult := createOrganization(client, chef.Organization{Name: "org1", FullName: "organization1"})
	fmt.Println("Added org1", orgResult)

	orgResult = createOrganization(client, chef.Organization{Name: "org1", FullName: "organization1"})
	fmt.Println("Added org1 again", orgResult)

	orgResult = createOrganization(client, chef.Organization{Name: "org2", FullName: "organization2"})
	fmt.Println("Added org2", orgResult)

	orgout := getOrganization(client, org1)
	fmt.Println("Get org1", orgout)

	orgList = listOrganizations(client)
	fmt.Println("List organizations After adding org1 and org2", orgList)

	orgresult := updateOrganization(client, chef.Organization{Name: "org1", FullName: "new_organization1"})
	fmt.Println("Update org1", orgresult)

	orgout = getOrganization(client, org1)
	fmt.Println("Get org1 after update", orgout)

	orgErr := deleteOrganization(client, org2)
	fmt.Println("Delete org2", orgErr)

	orgList = listOrganizations(client)
	fmt.Println("List organizations after deleting org2", orgList)

        // Clean up
	orgErr = deleteOrganization(client, org1)
	fmt.Println("Result from deleting org1", orgErr)

	orgList = listOrganizations(client)
	fmt.Println("List organizations after cleanup", orgList)

}

// buildClient creates a connection to a chef server using the chef api.
func buildClient(user string, keyfile string, baseurl string, skipssl bool) *chef.Client {
	key := clientKey(keyfile)
	client, err := chef.NewClient(&chef.Config{
		Name:    user,
		Key:     string(key),
		BaseURL: baseurl,
		SkipSSL: skipssl,
		// goiardi is on port 4545 by default, chef-zero is 8889, chef-server is on 443
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

// createOrganization uses the chef server api to create a single organization
func createOrganization(client *chef.Client, org chef.Organization) chef.OrganizationResult {
	orgResult, err := client.Organizations.Create(org)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue creating org:", org, err)
	}
	return orgResult
}

// deleteOrganization uses the chef server api to delete a single organization
func deleteOrganization(client *chef.Client, name string) error {
	err := client.Organizations.Delete(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue deleting org:", name, err)
	}
	return err
}

// getOrganization uses the chef server api to get information for a single organization
func getOrganization(client *chef.Client, name string) chef.Organization {
        // todo: everything
	orgList, err := client.Organizations.Get(name)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing org:", name, err)
	}
	return orgList
}

// listOrganizations uses the chef server api to list all organizations
func listOrganizations(client *chef.Client) map[string]string {
	orgList, err := client.Organizations.List()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue listing orgs:", err)
	}
	return orgList
}

// updateOrganization uses the chef server api to update information for a single organization
func updateOrganization(client *chef.Client, org chef.Organization) chef.Organization {
	org_update, err := client.Organizations.Update(org)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Issue updating org:", org, err)
	}
	return org_update
}
