package chef

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
)

const clientKeyListResponseFile = "test/client_keys_response.json"
const clientKeyTestFile = "test/client_key.json"

var (
	testClientJSON = "test/client.json"
)

func TestClientFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testClientJSON); err != nil {
		t.Error("unexpected error", err, "during os.Open on", testClientJSON)
	} else {
		dec := json.NewDecoder(file)
		var n Client
		if err := dec.Decode(&n); err == io.EOF {
			log.Println(n)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func TestClientsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"client1": "http://localhost/clients/client1", "client2": "http://localhost/clients/client2"}`)
	})
	response, err := client.Clients.List()
	if err != nil {
		t.Errorf("Clients.List returned error: %v", err)
	}

	// The order printed by the String function is not defined
	want := "client1 => http://localhost/clients/client1\nclient2 => http://localhost/clients/client2\n"
	want2 := "client2 => http://localhost/clients/client2\nclient1 => http://localhost/clients/client1\n"
	rstr := response.String()
	if rstr != want && rstr != want2 {
		t.Errorf("Clients.List returned:\n%+v\nwant:\n%+v\n", rstr, want)
	}
}

func TestClientsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
      "clientname": "client1",
      "orgname": "org_name",
      "validator": false,
      "certificate": "-----BEGIN CERTIFICATE-----",
      "name": "node_name"
    }`)
	})

	_, err := client.Clients.Get("client1")
	if err != nil {
		t.Errorf("Clients.Get returned error: %v", err)
	}
}

func TestClientsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"uri": "http://localhost/clients/client", "private_key": "-----BEGIN PRIVATE KEY-----"}`)
	})

	response, err := client.Clients.Create("client", false)
	if err != nil {
		t.Errorf("Clients.Create returned error: %v", err)
	}

	want := &ApiClientCreateResult{Uri: "http://localhost/clients/client", PrivateKey: "-----BEGIN PRIVATE KEY-----"}
	if !reflect.DeepEqual(response, want) {
		t.Errorf("Clients.Create returned %+v, want %+v", response, want)
	}
}

func TestClientsService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"name": "client1", "json_class": "Chef::Client", "chef_type": "client"}`)
	})

	err := client.Clients.Delete("client1")
	if err != nil {
		t.Errorf("Clients.Delete returned error: %v", err)
	}
}

func TestClientsService_ListKeys(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(clientKeyListResponseFile)
	if err != nil {
		t.Fatal(err)
	}

	mux.HandleFunc("/clients/client1/keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	keys, err := client.Clients.ListKeys("client1")
	if err != nil {
		t.Fatal(err)
	}

	if len(*keys) != 2 {
		t.Error("expected len(keys) to be 2")
	}
}

func TestClientsService_GetKey(t *testing.T) {
	setup()
	defer teardown()

	file, err := ioutil.ReadFile(clientKeyTestFile)
	if err != nil {
		t.Fatal(err)
	}

	mux.HandleFunc("/clients/client1/keys/default", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(file))
	})

	key, err := client.Clients.GetKey("client1", "default")
	if err != nil {
		t.Fatal(err)
	}

	if key.PublicKey != "-----BEGIN PUBLIC KEY-----\n...\n-----END PUBLIC KEY-----" {
		t.Error("expected key.PublicKey to match fixture")
	}
}
