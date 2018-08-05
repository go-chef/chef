package chef

import (
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

// func TestVaultsService_DeleteItem(t *testing.T) {
//     setup()
//     defer teardown()

//     mux.HandleFunc("/data/bag1/item1", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, ``)
//     })

//     err := client.Vaults.DeleteItem("bag1", "item1")
//     if err != nil {
//         t.Errorf("Vaults.DeleteItem returned error: %v", err)
//     }
// }

// func TestVaultsService_UpdateItem(t *testing.T) {
//     setup()
//     defer teardown()

//     mux.HandleFunc("/data/bag1/item1", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, ``)
//     })

//     dbi := map[string]string{
//         "id":  "item1",
//         "foo": "test123",
//     }

//     err := client.Vaults.UpdateItem("bag1", "item1", dbi)
//     if err != nil {
//         t.Errorf("Vaults.UpdateItem returned error: %v", err)
//     }
// }

// func TestVaultsService_VaultListResultString(t *testing.T) {
//     e := &VaultListResult{"bag1": "http://localhost/data/bag1", "bag2": "http://localhost/data/bag2"}
//     want := "bag1 => http://localhost/data/bag1\nbag2 => http://localhost/data/bag2\n"
//     if e.String() != want {
//         t.Errorf("VaultListResult.String returned:\n%+v\nwant:\n%+v\n", e.String(), want)
//     }
// }
