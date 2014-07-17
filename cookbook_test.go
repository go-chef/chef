package chef

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const cookbookListResponseFile = "test/cookbooks_response.json"

func TestCookbookList(t *testing.T) {
	setup()
	defer teardown()

	responseData, err := ioutil.ReadFile(cookbookListResponseFile)
	if err != nil {
		t.Error(err)
	}

	mux.HandleFunc("/cookbooks", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, string(responseData))
	})

	_, err = client.Cookbooks.List()
	if err != nil {
		t.Error(err)
	}

	_, err = client.Cookbooks.ListVersions("3")
	if err != nil {
		t.Error(err)
	}

	_, err = client.Cookbooks.ListVersions("0")
	if err == nil {
		t.Error("0 value ListVersion should not be allowed", err)
	}
}

func TestCookbookListVersions_0(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/cookbooks", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "BAD FUCKING REQUEST", 503)
	})

	_, err := client.Cookbooks.ListVersions("2")
	if err == nil {
		t.Error("We expected this bad request to error", err)
	}
}
