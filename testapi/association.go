//
// Test the go-chef/chef chef server api /organization/:org/user and /organization/:org/association_requests
// endpoints against a live server
//
package testapi

import (
	"fmt"
	"github.com/go-chef/chef"
	"os"
)

// association exercise the chef server api
func Association() {
	client := Client()

	// Build stuctures to invite users
	invite := chef.Request{
		User: "usrinvite",
	}
	invite2 := chef.Request{
		User: "usr2invite",
	}
	invitemissing := chef.Request{
		User: "nouser",
	}
	add1 := chef.AddNow{
		Username: "usradd",
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

	// Find a pending invitation by user name
	id, err := client.Associations.InviteId("usr2invite")
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
	for _, in := range outl {
		outd, err := client.Associations.DeleteInvite(in.Id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Issue deleting an invitation for %s  %+v\n", in.UserName, err)
		}
		fmt.Printf("Deleted invitation %s for %s %+v\n", in.Id, in.UserName, outd)
	}

	// Add a user to the test organization
	err = client.Associations.Add(add1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue adding user usradd: %+v\n", err)
	}
	fmt.Printf("User added: %+v\n", add1)
	// List the users
	ulist, err := client.Associations.List()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue listing the users: %+v\n", err)
	}
	fmt.Printf("Users list: %+v\n", ulist)
	// Get the user details
	uget, err := client.Associations.Get("usradd")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue getting user details: %+v\n", err)
	}
	fmt.Printf("User details: %+v\n", uget)
	// Delete a user from the organization
	udel, err := client.Associations.Get("usradd")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Issue deleting usradd: %+v\n", err)
	}
	fmt.Printf("User deleted: %+v\n", udel)
}
