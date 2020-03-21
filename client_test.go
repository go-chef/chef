package chef

import (
	"encoding/json"
	"fmt"
	"github.com/r3labs/diff"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
)

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

func TestClientsService_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"uri": "http://localhost/clients/client", 
		"chef_key": {
		  "name": "default",
		  "expiration_date": "infinity",
		  "uri": "http://localhost/clients/client/keys/default",
		  "public_key": "-----BEGIN PUBLIC KEY-----",
		  "private_key": "-----BEGIN PRIVATE KEY-----"
	        }
	  }`)
	})

	newclient := ApiNewClient{Name: "client"}
	response, err := client.Clients.Create(newclient)
	if err != nil {
		t.Errorf("Clients.Create returned error: %v", err)
	}

	want := &ApiClientCreateResult{
		Uri: "http://localhost/clients/client",
		ChefKey: ChefKey{
			Name:           "default",
			ExpirationDate: "infinity",
			Uri:            "http://localhost/clients/client/keys/default",
			PrivateKey:     "-----BEGIN PRIVATE KEY-----",
			PublicKey:      "-----BEGIN PUBLIC KEY-----",
		},
	}
	if !reflect.DeepEqual(response, want) {
		t.Errorf("Clients.Create returned %+v, want %+v", response, want)
	}
}

func TestClientsService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{}`)
	})

	err := client.Clients.Delete("client1")
	if err != nil {
		t.Errorf("Clients.Delete returned error: %v", err)
	}
}

func TestClientsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
      "name": "node_name",
      "clientname": "client1",
      "validator": false,
      "orgname": "org_name",
      "json_class": "Chef::ApiClient",
      "chef_type": "client"
    }`)
	})

	_, err := client.Clients.Get("client1")
	if err != nil {
		t.Errorf("Clients.Get returned error: %v", err)
	}
}

func TestClientsService_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{
      "name": "client2",
      "clientname": "client2",
      "validator": false,
      "json_class": "Chef::ApiClient",
      "chef_type": "client"
    }`)
	})

	apinewclient := ApiNewClient{
		Name:      "client2",
		Validator: false,
	}
	updateresult, err := client.Clients.Update("client1", apinewclient)
	if err != nil {
		t.Errorf("Clients.Update returned error: %v", err)
	}
	want := ApiClient{
		Name:       "client2",
		ClientName: "client2",
		Validator:  false,
		JsonClass:  "Chef::ApiClient",
		ChefType:   "client",
	}
	if !reflect.DeepEqual(updateresult, &want) {
		t.Errorf("Clients.Update returned %+v, want %+v", updateresult, want)
	}
}

func TestCListKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1/keys", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `[
			       {
				       "name": "default",
                                	"uri": "https://chefserver/clients/client1/keys/default",
                                	"expired": false
                         	}
		 	]`)
		}
	})

	keyresult, err := client.Clients.ListKeys("client1")
	if err != nil {
		t.Errorf("Clients.ListKeys returned error: %v", err)
	}
	defaultItem := KeyItem{
		Name:    "default",
		Uri:     "https://chefserver/clients/client1/keys/default",
		Expired: false,
	}
	Want := []KeyItem{defaultItem}
	if !reflect.DeepEqual(keyresult, Want) {
		t.Errorf("Clients.ListKeys returned %+v, want %+v", keyresult, Want)
	}
}

func TestCAddKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1/keys", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST":
			fmt.Fprintf(w, `{
             			        "name": "newkey",
                                	"uri": "https://chefserver/clients/client1/keys/newkey",
                                	"expired": false
                         	}`)
		}
	})

	keyadd := AccessKey{
		Name:           "newkey",
		PublicKey:      "RSA KEY",
		ExpirationDate: "infinity",
	}
	keyresult, err := client.Clients.AddKey("client1", keyadd)
	if err != nil {
		t.Errorf("Clients.AddKey returned error: %v", err)
	}
	Want := KeyItem{
		Name:    "newkey",
		Uri:     "https://chefserver/clients/client1/keys/newkey",
		Expired: false,
	}
	if !reflect.DeepEqual(keyresult, Want) {
		t.Errorf("Clients.AddKey returned %+v, want %+v", keyresult, Want)
	}
}

func TestCDeleteKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1/keys/newkey", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "DELETE":
			fmt.Fprintf(w, `{
             			        "name": "newkey",
                                	"public_key": "RSA KEY",
                                	"expiration_date": "infinity"
                         	}`)
		}
	})

	keyresult, err := client.Clients.DeleteKey("client1", "newkey")
	if err != nil {
		t.Errorf("Clients.DeleteKey returned error: %v", err)
	}
	Want := AccessKey{
		Name:           "newkey",
		PublicKey:      "RSA KEY",
		ExpirationDate: "infinity",
	}
	if !reflect.DeepEqual(keyresult, Want) {
		diff, _ := diff.Diff(keyresult, Want)
		t.Errorf("Clients.DeleteKey returned %+v, want %+v, differences %+v", keyresult, Want, diff)
	}
}

func TestCGetKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1/keys/newkey", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{
             			        "name": "newkey",
                                	"public_key": "RSA KEY",
                                	"expiration_date": "infinity"
                         	}`)
		}
	})

	keyresult, err := client.Clients.GetKey("client1", "newkey")
	if err != nil {
		t.Errorf("Clients.GetKey returned error: %v", err)
	}
	Want := AccessKey{
		Name:           "newkey",
		PublicKey:      "RSA KEY",
		ExpirationDate: "infinity",
	}
	if !reflect.DeepEqual(keyresult, Want) {
		diff, _ := diff.Diff(keyresult, Want)
		t.Errorf("Clients.GetKey returned %+v, want %+v, differences %+v", keyresult, Want, diff)
	}
}

func TestCUpdateKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/clients/client1/keys/newkey", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "PUT":
			fmt.Fprintf(w, `{
             			        "name": "newkey",
                                	"public_key": "RSA NEW KEY",
                                	"expiration_date": "infinity"
                         	}`)
		}
	})

	updkey := AccessKey{
		Name:           "newkey",
		PublicKey:      "RSA NEW KEY",
		ExpirationDate: "infinity",
	}
	keyresult, err := client.Clients.UpdateKey("client1", "newkey", updkey)
	if err != nil {
		t.Errorf("Clients.UpdateKey returned error: %v", err)
	}
	Want := AccessKey{
		Name:           "newkey",
		PublicKey:      "RSA NEW KEY",
		ExpirationDate: "infinity",
	}
	if !reflect.DeepEqual(keyresult, Want) {
		diff, _ := diff.Diff(keyresult, Want)
		t.Errorf("Clients.UpdateKey returned %+v, want %+v, differences %+v", keyresult, Want, diff)
	}
}
