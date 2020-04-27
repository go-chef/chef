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
		Name:      "client_node",
		Type:      "client",
		PublicKey: "-----BEGIN PUBLIC KEY",
		AuthzId:   "afe1234",
		OrgMember: true,
	}
	pWant.Principals = append(pWant.Principals, client)

	if !reflect.DeepEqual(p, pWant) {
		t.Errorf("Unexpected principal values got: %+v wanted: %+v", p, pWant)
	}
}
