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
		fmt.Fprintf(w, `{"bag1":"http://localhost/data/bag1", "bag2":"http://localhost/data/bag2"}`)
	})

	databags, err := client.Vaults.List()
	if err != nil {
		t.Errorf("Vaults.List returned error: %v", err)
	}

	want := &VaultListResult{"bag1": "http://localhost/data/bag1", "bag2": "http://localhost/data/bag2"}
	if !reflect.DeepEqual(databags, want) {
		t.Errorf("Vaults.List returned %+v, want %+v", databags, want)
	}
}

// func TestVaultsService_Create(t *testing.T) {
//     setup()
//     defer teardown()

//     mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, `{"uri": "http://localhost/data/newdatabag"}`)
//     })

//     databag := &Vault{Name: "newdatabag"}
//     response, err := client.Vaults.Create(databag)
//     if err != nil {
//         t.Errorf("Vaults.Create returned error: %v", err)
//     }

//     want := &VaultCreateResult{URI: "http://localhost/data/newdatabag"}
//     if !reflect.DeepEqual(response, want) {
//         t.Errorf("Vaults.Create returned %+v, want %+v", response, want)
//     }
// }

// func TestVaultsService_Delete(t *testing.T) {
//     setup()
//     defer teardown()

//     mux.HandleFunc("/data/databag", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, `{"name": "databag", "json_class": "Chef::Vault", "chef_type": "data_bag"}`)
//     })

//     response, err := client.Vaults.Delete("databag")
//     if err != nil {
//         t.Errorf("Vaults.Delete returned error: %v", err)
//     }

//     want := &Vault{
//         Name:      "databag",
//         JsonClass: "Chef::Vault",
//         ChefType:  "data_bag",
//     }

//     if !reflect.DeepEqual(response, want) {
//         t.Errorf("Vaults.Delete returned %+v, want %+v", response, want)
//     }
// }

// func TestVaultsService_ListItems(t *testing.T) {
//     setup()
//     defer teardown()

//     mux.HandleFunc("/data/bag1", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, `{"item1":"http://localhost/data/bag1/item1", "item2":"http://localhost/data/bag1/item2"}`)
//     })

//     databags, err := client.Vaults.ListItems("bag1")
//     if err != nil {
//         t.Errorf("Vaults.ListItems returned error: %v", err)
//     }

//     want := &VaultListResult{"item1": "http://localhost/data/bag1/item1", "item2": "http://localhost/data/bag1/item2"}
//     if !reflect.DeepEqual(databags, want) {
//         t.Errorf("Vaults.ListItems returned %+v, want %+v", databags, want)
//     }
// }

// func TestVaultsService_GetItem(t *testing.T) {
//     setup()
//     defer teardown()

//     mux.HandleFunc("/data/bag1/item1", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, `{"id":"item1", "stuff":"things"}`)
//     })

//     _, err := client.Vaults.GetItem("bag1", "item1")
//     if err != nil {
//         t.Errorf("Vaults.GetItem returned error: %v", err)
//     }
// }

// func TestVaultsService_CreateItem(t *testing.T) {
//     setup()
//     defer teardown()

//     mux.HandleFunc("/data/bag1", func(w http.ResponseWriter, r *http.Request) {
//         fmt.Fprintf(w, ``)
//     })

//     dbi := map[string]string{
//         "id":  "item1",
//         "foo": "test123",
//     }

//     err := client.Vaults.CreateItem("bag1", dbi)
//     if err != nil {
//         t.Errorf("Vaults.CreateItem returned error: %v", err)
//     }
// }

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
