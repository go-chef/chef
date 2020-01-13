package chef

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestVaultsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"secrets":"http://localhost/data/secrets","secrets_keys":"http://localhost/data/secret_keys","bag1":"http://localhost/data/bag1"}`)
	})

	databags, err := client.Vaults.List()
	if err != nil {
		t.Errorf("Vaults.List returned error: %v", err)
	}

	want := &VaultListResult{"secrets": "http://localhost/data/secrets"}
	if !reflect.DeepEqual(databags, want) {
		t.Errorf("Vaults.List returned %+v, want %+v", databags, want)
	}
}

func TestVaultsService_GetItem(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/data/vaults/secrets", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"id":"secrets"}`)
	})

	mux.HandleFunc("/data/vaults/secrets_keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"id":"secrets_keys"}`)
	})
	_, err := client.Vaults.GetItem("vaults", "secrets")
	if err != nil {
		t.Errorf("Vaults.GetItem returned error: %v", err)
	}
}

func TestVaultsService_CreateItem(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/data/vaults", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ``)
	})

	mux.HandleFunc("/data/vaults/secrets_keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ``)
	})

	_, err := client.Vaults.CreateItem("vaults", "secrets")
	if err != nil {
		t.Errorf("Vaults.CreateItem returned error: %v", err)
	}
}

func TestVaultsService_DeleteItem(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/data/vaults/secrets", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ``)
	})

	mux.HandleFunc("/data/vaults/secret_keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ``)
	})

	err := client.Vaults.DeleteItem("vaults", "secrets")
	if err != nil {
		t.Errorf("Vaults.DeleteItem returned error: %v", err)
	}
}

func TestVaultsService_UpdateItem(t *testing.T) {
	setup()
	defer teardown()
	var secretsData string

	mux.HandleFunc("/data/vaults/secrets_keys", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ``)
	})

	mux.HandleFunc("/data/vaults/secrets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			fmt.Fprint(w, secretsData)
		case "PUT":
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			secretsData = buf.String()
			fmt.Fprintf(w, ``)
		default:
			fmt.Fprintf(w, ``)
		}
	})

	mux.HandleFunc("/data/vaults", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, ``)
	})

	data := map[string]interface{}{
		"id":  "secrets",
		"foo": "bar",
	}

	item, err := client.Vaults.CreateItem("vaults", "secrets")
	if err != nil {
		t.Fatalf("Vaults.CreateItem returned an error: %v", err)
	}
	if item == nil {
		t.Fatalf("Vaults.CreateItem returned nothing: %q", err)
	}

	err = client.Vaults.UpdateItem(item, data)
	if err != nil {
		t.Fatalf("Vaults.UpdateItem returned an error: %v", err)
	}

	if secretsData == "" {
		t.Fatalf("Vaults.UpdateItem did not update the data bag: %v", err)
	}

	updatedData, err := item.Decrypt()
	if err != nil {
		t.Fatalf("Vaults.Decrypt returned an error: %v", err)
	}

	if !reflect.DeepEqual(*updatedData, data) {
		t.Fatalf("Updated data did not match: %v != %v", updatedData, data)
	}
}
