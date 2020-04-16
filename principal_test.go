package chef

import (
	"fmt"
	"net/http"
	_ "reflect"
	"testing"
)

func TestPrincipalsGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/principals/client_node", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
		"name": "client_node",
		"type": "client",
		"public_key": "-----BEGIN PUBLIC KEY-----No, not really-----END PUBLIC KEY-----"
}`)
	})

	p, err := client.Principals.Get("client_node")
	if err != nil {
		t.Errorf("GET principal error making request: %v", err)
		return
	}

	if p.Name != "client_node" {
		t.Errorf("Unexpected principal name: %v", p.Name)
	}
	if p.Type != "client" {
		t.Errorf("Unexpected principal type: %v", p.Type)
	}
	if p.PublicKey != "-----BEGIN PUBLIC KEY-----No, not really-----END PUBLIC KEY-----" {
		t.Errorf("Unexpected principal public key: %v", p.PublicKey)
	}
}
