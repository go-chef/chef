package chef

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequiredRecipeGet(t *testing.T) {
	setup()
	defer teardown()

	recipeText := "file 'test'"
	mux.HandleFunc("/required_recipe", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, recipeText)
	})

	wantRecipe := RequiredRecipe(recipeText)

	recipe, err := client.RequiredRecipe.Get()
	if err != nil {
		t.Errorf("RequiredRecipe.Get returned error: %s", err.Error())
	}

	if recipe != wantRecipe {
		t.Errorf("RequiredRecipe.Get returned %+v, want %+v, error %+v", recipe, wantRecipe, err)
	}

}
