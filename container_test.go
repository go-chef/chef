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
	testContainerJSON = "test/container.json"
)

func TestContainerFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testContainerJSON); err != nil {
		t.Error("unexpected error", err, "during os.Open on", testContainerJSON)
	} else {
		dec := json.NewDecoder(file)
		var n Container
		if err := dec.Decode(&n); err == io.EOF {
			log.Println(n)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func TestContainersService_Methods(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{"container1":"https://chef/containers/container1", "container2":"https://chef/containers/container2"}`)
		case r.Method == "POST":
			fmt.Fprintf(w, `{ "uri": "http://localhost:4545/containers/container1" }`)
		}
	})

	mux.HandleFunc("/containers/container1", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" || r.Method == "PUT":
			fmt.Fprintf(w, `{
	    "containername": "container1",
	    "containerpath": "container1"
		}`)
		case r.Method == "DELETE":
		}
	})

	// Test list
	containers, err := client.Containers.List()
	if err != nil {
		t.Errorf("Containers.List returned error: %v", err)
	}

	listWant := ContainerListResult{"container1": "https://chef/containers/container1", "container2": "https://chef/containers/container2"}

	if !reflect.DeepEqual(containers, listWant) {
		t.Errorf("Containers.List returned %+v, want %+v", containers, listWant)
	}

	// test Get
	container, err := client.Containers.Get("container1")
	if err != nil {
		t.Errorf("Containers.Get returned error: %v", err)
	}

	wantContainer := Container{
		ContainerName: "container1",
		ContainerPath: "container1",
	}
	if !reflect.DeepEqual(container, wantContainer) {
		t.Errorf("Containers.Get returned %+v, want %+v", container, wantContainer)
	}

	// test Post
	res, err := client.Containers.Create(wantContainer)
	if err != nil {
		t.Errorf("Containers.Post returned error: %s", err.Error())
	}

	postResult := &ContainerCreateResult{Uri: "http://localhost:4545/containers/container1"}
	if !reflect.DeepEqual(postResult, res) {
		t.Errorf("Containers.Post returned %+v, want %+v", res, postResult)
	}

	// test Delete
	err = client.Containers.Delete(container.ContainerName)
	if err != nil {
		t.Errorf("Containers.Delete returned error %+v", err)
	}
}
