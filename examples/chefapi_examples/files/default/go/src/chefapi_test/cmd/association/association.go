//
// Test the go-chef/chef chef server api /organization/:org/user and /organization/:org/association_requests
// endpoints against a live server
//
package main

import (
	"fmt"
        "chefapi_test/testapi"
        "github.com/go-chef/chef"
	"os"
)


// main Exercise the chef server api
func main() {
	client := testapi.Client()

	// Build stuctures to invite users
        invite := chef.Request {
		User: "usrinvite",
	}
        invite2 := chef.Request {
		User: "usr2invite",
	}
        invitemissing := chef.Request {
		User: "nouser",
	}

	// Invite the user to the test org
	out, err := client.Associations.Invite(invite)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue inviting a user %+v %+v\n", invite, err)
	}
	fmt.Printf("Invited user %+v %+v\n", invite, out)

	// Invite a second user
	out, err = client.Associations.Invite(invite2)
	fmt.Printf("Invited user %+v %+v\n", invite2, out)

	// fail at inviting a missing user.  Should get a 404
	out, err = client.Associations.Invite(invitemissing)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue inviting a user %+v %+v\n", invitemissing, err)
	}
	fmt.Printf("Invited user %+v %+v\n", invitemissing, out)

	// Find the pending invitation by user name
	id, err :=  client.Associations.InviteId("usr2invite")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue finding an invitation for usr2invite %+v\n", err)
	}
	fmt.Printf("Invitation id for usr2invite %+v\n", id)

	// Accept the invite for invite2
	// outa, err := client.Associations.AcceptInvite(id)
	// if err != nil {
		// fmt.Fprintf(os.Stderr, "Issue accepting the invitation %+v\n", err)
	// }
	// fmt.Printf("Accept invitation %+v\n", outa)

	// List the invites
	outl, err := client.Associations.ListInvites()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing the invitations %+v\n", err)
	}
	fmt.Printf("Invitation list %+v\n", outl)

	// Delete the invitations by id
	for  _, in := range outl {
                outd, err := client.Associations.DeleteInvite(in.Id)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "Issue deleting an invitation for %s  %+v\n", in.UserName, err)
                }
                fmt.Printf("Deleted invitation %s for %s %+v\n", in.Id, in.UserName, outd)
        }
}
