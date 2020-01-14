chefapi_test//
// Test the go-chef/chef chef vault support against a live server
//
package main

import (
	"fmt"
        "chefapi_test/testapi"
        "github.com/go-chef/chef"
	"os"
)


// main Exercise the chef server vault support
func main() {
	client := testapi.Client()

	// Add users to the test organizations
        addv := chef.AddNow { Username: "usrv", }
        addv2 := chef.AddNow { Username: "usrv2", }
        client.Associations.Add(addv)
        client.Associations.Add(addv2)
	// TODO: Add a node

	// List vaults before an item is created
	vaultList, err := client.Vaults.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing vaults %+v\n",err)
	}
	fmt.Printf("List vaults before creation %+v\n", vaultList)

        // Create a data bag to hold the vault items
        databag := &chef.DataBag{Name: "testv"}
        response, err := client.DataBags.Create(databag)
        if err != nil {
                fmt.Fprintf(os.Stderr, "Issue creating data bag testv %+v\n",err)
        }
        fmt.Printf("Data bag created %+v\n", response)

	// Create a vault item
	item, err := client.Vaults.CreateItem("testv", "secrets")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue creating testv secrets vault item %+v\n", err)
	} 
	fmt.Printf("Created testv secrets vault item %+v\n", item)

	// Add content to the vault item
	// The vault item has pointers and must not be nil
	data := map[string]interface{}{
                "id":  "secrets",
                "foo": "bar",
        }
	err = client.Vaults.UpdateItem(item, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue updating testv secrets vault item %+v\n", err)
	}
	fmt.Println("Updated testv secrets vault item")

	// TODO:  Listt the items in a vault
	//  Make sure the item is encrypted

	// List vaults after an item is created
	vaultList, err = client.Vaults.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing vaults %+v\n",err)
	}
	fmt.Printf("List vaults after creation %+v\n", vaultList)


	// Get vault item
	vaultItem, err := client.Vaults.GetItem("testv", "secrets")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue getting vault item testv secrets %+v\n",err)
	}
	fmt.Printf("Get testv secrets vault item%+v\n", vaultItem)
	fmt.Printf("Show testv databag item %+v\n", *vaultItem.DataBagItem)
	fmt.Printf("Show testv keys %+v\n", vaultItem.Keys)
	fmt.Printf("Show testv keys %+v\n", vaultItem.Keys)
	fmt.Printf("Show testv keys databagitem %+v\n", *vaultItem.Keys.DataBag.Item)
	fmt.Printf("Show testv admins %+v\n", *vaultItem.Keys.DataBag.Item.clients)

	// Show contents of the vault item
	// Must get the item before decrypting
	contents, err := vaultItem.Decrypt()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue decrypting vault item testv secrets  %+v\n",err)
	}
	fmt.Printf("List initial vault item values %+v\n", contents)

	// Add content to the vault item after we get it
	// The vault item has pointers and must not be nil
	data := map[string]interface{}{
                "id":  "secrets",
                "foo": "bar",
                "jellico": "bats",
        }
	err = client.Vaults.UpdateItem(vaultItem, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue updating based on get of testv secrets vault item %+v\n", err)
	}
	fmt.Println("Updated based on get of testv secrets vault item")

	// TODO:  List the items in a vault after update

	// Delete vault contents
	err = client.Vaults.DeleteItem("testv", "secrets")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue deleting vault item testv secrets %+v\n",err)
	}
	fmt.Println("Delete testv secrets vault item")

	// List vaults after all items are deleted
	vaultList, err = client.Vaults.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing vaults %+v\n",err)
	}
	fmt.Printf("List vaults after creation %+v\n", vaultList)

        // Delete the data bag
        outBag, err := client.DataBags.Delete("testv")
        if err != nil {
                fmt.Fprintf(os.Stderr, "Issue deleting data bag testv %+v\n",err)
        }
        fmt.Printf("Data bag deleted %+v\n", outBag)

}

// add user and node to the admin and client list
// Vaults.GetItem(vaultName, itemName)  (*VaultItem, error)
// Add user2
// Change a value
// Do things using usrv id - admin
// Do things using usrv2 id - client only
// Show key value
// Show encrypted values
// Show admins
// Show clients
// Show search
