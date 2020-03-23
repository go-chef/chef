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
	testLicenseJSON = "test/license.json"
)

func TestLicenseFromJSONDecoder(t *testing.T) {
	if file, err := os.Open(testLicenseJSON); err != nil {
		t.Error("Unexpected error '", err, "' during os.Open on", testLicenseJSON)
	} else {
		dec := json.NewDecoder(file)
		var g License
		if err := dec.Decode(&g); err == io.EOF {
			log.Fatal(g)
		} else if err != nil {
			log.Fatal(err)
		}
	}
}

func TestLicenseGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/license", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var request License
		dec.Decode(&request)
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{
                           "limit_exceeded": false,
                           "node_license": 25,
                           "node_count": 12,
                           "upgrade_url": "http://www.chef.io/contact/on-premises-simple"
		        }`)
		}
	})

	wantLicense := License{
		LimitExceeded: false,
		NodeLicense:   25,
		NodeCount:     12,
		UpgradeUrl:    "http://www.chef.io/contact/on-premises-simple",
	}

	license, err := client.License.Get()
	if err != nil {
		t.Errorf("License.Get returned error: %s", err.Error())
	}
	if !reflect.DeepEqual(license, wantLicense) {
		t.Errorf("License.Get returned %+v, want %+v", license, wantLicense)
	}

}
