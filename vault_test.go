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
		fmt.Fprintf(w, `{"secrets":"http://localhost/data/secrets","more":"http://localhost/data/more","justabag":"http://localhost/data/justabag"}`)
	})

	mux.HandleFunc("/data/secrets", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"item1":"http://localhost/data/secrets/item1", "item1_keys":"http://localhost/data/secrets/item1_keys"}`)
	})

	mux.HandleFunc("/data/more", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"item1":"http://localhost/data/more/item1", "item1_keys":"http://localhost/data/more/item1_keys"}`)
	})

	mux.HandleFunc("/data/justabag", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"item1":"http://localhost/data/justabag/item1"}`)
	})

	vaults, err := client.Vaults.List()
	if err != nil {
		t.Errorf("Vaults.List returned error: %v", err)
	}

	want := &VaultListResult{"secrets": "http://localhost/data/secrets", "more": "http://localhost/data/more"}
	if !reflect.DeepEqual(vaults, want) {
		t.Errorf("Vaults.List returned %+v, want %+v", vaults, want)
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

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"uri": "http://localhost/data/vaults"}`)
	})

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

func TestVaultsService_CreateItem_Existing(t *testing.T) {
	setup()
	defer teardown()

	// Return 409 when the attempt is made to create an existing vault
	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	})

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

func TestVaultsService_ListItems(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/data/vault1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"item1":"http://localhost/data/vault1/item1", "item1_keys":"http://localhost/data/vault1/item1_keys"}`)
	})

	vaultItems, err := client.Vaults.ListItems("vault1")
	if err != nil {
		t.Errorf("Vaults.ListItem returned error: %v", err)
	}

	want := []string{"item1"}
	if !reflect.DeepEqual(vaultItems, want) {
		t.Fatalf("Vault items returned did not match desired: %v != %v", vaultItems, want)
	}
}

func TestVaultsService_UpdateItem(t *testing.T) {
	setup()
	defer teardown()
	var secretsData string

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"uri": "http://localhost/data/vaults"}`)
	})

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
