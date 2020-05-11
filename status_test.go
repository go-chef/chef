package chef

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
)

var (
	testStatusJSON = "test/status.json"
)

func TestStatusFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testStatusJSON); err != nil {
		t.Error("Unexpected error '", err, "' during os.Open on", testStatusJSON)
	} else {
		dec := json.NewDecoder(file)
		var g Status
		if err := dec.Decode(&g); err == io.EOF {
			log.Fatal(g)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func TestStatusGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/_status", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var request Status
		dec.Decode(&request)
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{
                        	"status": "pong",
	                        "upstreams": {
		                        "chef_elasticsearch": "pong",
		                        "chef_sql": "pong",
		                        "chef_index": "pong",
		                        "oc_chef_authz": "pong",
		                        "data_collector": "pong"
	                        },
	                        "keygen": {
		                        "keys": 10,
		                        "max": 10,
		                        "max_workers": 2,
		                        "cur_max_workers": 2,
		                        "inflight": 0,
		                        "avail_workers": 2,
		                        "start_size": 0
	                        }
                        }`)
		}

	})

	wantStatus := Status{
		Status:    "pong",
		Upstreams: map[string]string{"chef_elasticsearch": "pong", "chef_index": "pong", "chef_sql": "pong", "data_collector": "pong", "oc_chef_authz": "pong"},
		Keygen:    map[string]int{"avail_workers": 2, "cur_max_workers": 2, "inflight": 0, "keys": 10, "max": 10, "max_workers": 2, "start_size": 0},
	}

	status, err := client.Status.Get()
	if err != nil {
		t.Errorf("Status.Get returned error: %s", err.Error())
	}

	if !reflect.DeepEqual(status, wantStatus) {
		t.Errorf("Status.Get returned %+v, want %+v, error %+v", status, wantStatus, err)
	}

}
