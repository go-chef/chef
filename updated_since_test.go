package chef

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestUpdatedSinceGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/updated_since", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			w.WriteHeader(404)
		}
	})

	_, err := client.UpdatedSince.Get(1)
	if !strings.Contains(fmt.Sprint(err), "404") {
		t.Errorf("Non 404 return code: %+v", err)
	}
}
