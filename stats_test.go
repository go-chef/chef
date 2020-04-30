package chef

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestStatsGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/_stats", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET":
			fmt.Fprintf(w, `[{"stat": "value"}]`)
		}
	})

	wantStats := Stats{
		map[string]interface{}{
			"stat": "value",
		},
	}

	stats, err := client.Stats.Get("json", "statsuser", "password")
	if err != nil {
		t.Errorf("Stat.Get returned error: %s", err.Error())
	}

	if !reflect.DeepEqual(stats, wantStats) {
		t.Errorf("Stat.Get returned %+v, want %+v, error %+v", stats, wantStats, err)
	}

}
