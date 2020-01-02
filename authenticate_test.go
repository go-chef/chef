package chef

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
        "github.com/stretchr/testify/assert"
)

var (
	testAuthenticateJSON = "test/authenticate.json"
)

func TestAuthenticateFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testAuthenticateJSON); err != nil {
		t.Error("Unexpected error '", err, "' during os.Open on", testAuthenticateJSON)
	} else {
		dec := json.NewDecoder(file)
		var g Authenticate
		if err := dec.Decode(&g); err == io.EOF {
			log.Fatal(g)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func TestAuthenticatesCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/authenticate_user", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var request Authenticate
		dec.Decode(&request)
		switch {
		case r.Method == "POST":
			if request.Password == "password" {
				w.WriteHeader(200)
			} else {
				w.WriteHeader(401)
			}
		}
	})
	var request Authenticate
	request.UserName = "user1"
	request.Password = "invalid"
	err := client.AuthenticateUser.Authenticate(request)
	if assert.NotNil(t, err) {
               assert.Contains(t, err.Error(), "401")
        }

	request.UserName = "user1"
	request.Password = "password"
	err = client.AuthenticateUser.Authenticate(request)
	if err != nil {
		t.Errorf("Authenticate returned error: %+v", err)
	}
}
