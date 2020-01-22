package chef

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestUniverseGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/universe", func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var request Universe
		dec.Decode(&request)
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `{
			"ffmpeg": {
				"0.1.0": {
					"location_path": "http://supermarket.chef.io/api/v1/cookbooks/ffmpeg/0.1.0/download",
					"location_type": "supermarket",
					"dependencies": {
						"git": ">= 0.0.0",
						"build-essential": ">= 0.0.0",
						"libvpx": "~> 0.1.1",
						"x264": "~> 0.1.1"
					}
				},
				"0.1.1": {
					"location_path": "http://supermarket.chef.io/api/v1/cookbooks/ffmpeg/0.1.1/download",
					"location_type": "supermarket",
					"dependencies": {
						"git": ">= 0.0.0",
						"build-essential": ">= 0.0.0",
						"libvpx": "~> 0.1.1",
						"x264": "~> 0.1.1"
					}
				}
			},
			"pssh": {
				"0.1.0": {
					"location_path": "http://supermarket.chef.io/api/v1/cookbooks/pssh.1.0/download",
					"location_type": "supermarket",
					"dependencies": {}
				}
			}
		        }`)
		}
	})

	wantU := Universe{}
	wantU.Books = make(map[string]UniverseBook)
	ffmpeg := UniverseBook{}
	ffmpeg.Versions = make(map[string]UniverseVersion)
	ffmpeg.Versions["0.1.0"] = UniverseVersion{
		LocationPath: "http://supermarket.chef.io/api/v1/cookbooks/ffmpeg/0.1.0/download",
		LocationType: "supermarket",
		Dependencies: map[string]string{
			"git":             ">= 0.0.0",
			"build-essential": ">= 0.0.0",
			"libvpx":          "~> 0.1.1",
			"x264":            "~> 0.1.1",
		},
	}
	ffmpeg.Versions["0.1.1"] = UniverseVersion{
		LocationPath: "http://supermarket.chef.io/api/v1/cookbooks/ffmpeg/0.1.1/download",
		LocationType: "supermarket",
		Dependencies: map[string]string{
			"git":             ">= 0.0.0",
			"build-essential": ">= 0.0.0",
			"libvpx":          "~> 0.1.1",
			"x264":            "~> 0.1.1",
		},
	}
	pssh := UniverseBook{}
	pssh.Versions = make(map[string]UniverseVersion)
	pssh.Versions["0.1.0"] = UniverseVersion{
		LocationPath: "http://supermarket.chef.io/api/v1/cookbooks/pssh.1.0/download",
		LocationType: "supermarket",
		Dependencies: map[string]string{},
	}
	wantU.Books["ffmpeg"] = ffmpeg
	wantU.Books["pssh"] = pssh

	universe, err := client.Universe.Get()
	if err != nil {
		t.Errorf("Universe.Get returned error: %s", err.Error())
	}

	if !reflect.DeepEqual(universe, wantU) {
		t.Errorf("Universe.Get returned %+v, want %+v", universe, wantU)
	}
}
