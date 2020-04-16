package chef

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPrincipalsGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/principals/client_node", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
		"principals": [{
			"name": "client_node",
			"type": "client",
			"authz_id": "afe1234",
			"org_member": true,
			"public_key": "-----BEGIN PUBLIC KEY"
}]}`)
	})

	p, err := client.Principals.Get("client_node")
	if err != nil {
		t.Errorf("GET principal error %+v making request: ", err)
		return
	}

	pWant := Principal{}
	client := Principals{
		Name: "client_node",
		Type: "client",
		PublicKey: "-----BEGIN PUBLIC KEY",
		AuthzId: "afe1234",
		OrgMember: true,
	}
	pWant.Principals[0] = client

 //  type Principal struct {
 //         Principals []struct {
 //                 Name      string `json:"name"`
 //                 Type      string `json:"type"`
 //                 PublicKey string `json:"public_key"`
 //                 AuthzId   string `json:"authz_id"`
 //                 OrgMember bool   `json:"org_member"`
 //         } `json:"principals"`
 // }

	if !reflect.DeepEqual(p, pWant) {
		t.Errorf("Unexpected principal values got: %+v wanted: %+v", p, pWant)
	}
}
