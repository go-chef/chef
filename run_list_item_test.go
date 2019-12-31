package chef

import (
	"testing"
)

func TestQualifiedRecipeWithVersion(t *testing.T) {
	runListItem := "recipe[my_recipe@1.0.0]"
	rli, err := NewRunListItem(runListItem)
	if err != nil {
		t.Errorf(`NewRunListItem("%s") did not correctly parse.`, runListItem)
	}
	if rli.Name != "my_recipe" {
		t.Errorf(`NewRunListItem("%s") did not correctly set name`, runListItem)
	}
	if rli.Type != "recipe" {
		t.Errorf(`NewRunListItem("%s") did not correctly set type`, runListItem)
	}
	if rli.Version != "1.0.0" {
		t.Errorf(`NewRunListItem("%s") did not correctly set version`, runListItem)
	}
}

func TestQualifiedRecipe(t *testing.T) {
	runListItem := "recipe[my_recipe]"
	rli, err := NewRunListItem(runListItem)
	if err != nil {
		t.Errorf(`NewRunListItem("%s") did not correctly parse.`, runListItem)
	}
	if rli.Name != "my_recipe" {
		t.Errorf(`NewRunListItem("%s") did not correctly set name`, runListItem)
	}
	if rli.Type != "recipe" {
		t.Errorf(`NewRunListItem("%s") did not correctly set type`, runListItem)
	}
	if rli.Version != "" {
		t.Errorf(`NewRunListItem("%s") did not correctly set version`, runListItem)
	}
}

func TestQualifiedRole(t *testing.T) {
	runListItem := "role[my_role]"
	rli, err := NewRunListItem(runListItem)
	if err != nil {
		t.Errorf(`NewRunListItem("%s") did not correctly parse.`, runListItem)
	}
	if rli.Name != "my_role" {
		t.Errorf(`NewRunListItem("%s") did not correctly set name`, runListItem)
	}
	if rli.Type != "role" {
		t.Errorf(`NewRunListItem("%s") did not correctly set type, %s`, runListItem, rli.Type)
	}
	if rli.Version != "" {
		t.Errorf(`NewRunListItem("%s") did not correctly set version`, runListItem)
	}
}

func TestVersionedUnqualifiedRecipe(t *testing.T) {
	runListItem := "my_recipe@1.0.0"
	rli, err := NewRunListItem(runListItem)
	if err != nil {
		t.Errorf(`NewRunListItem("%s") did not correctly parse.`, runListItem)
	}
	if rli.Name != "my_recipe" {
		t.Errorf(`NewRunListItem("%s") did not correctly set name`, runListItem)
	}
	if rli.Type != "recipe" {
		t.Errorf(`NewRunListItem("%s") did not correctly set type`, runListItem)
	}
	if rli.Version != "1.0.0" {
		t.Errorf(`NewRunListItem("%s") did not correctly set version`, runListItem)
	}
}

func TestRecipeNameAlone(t *testing.T) {
	runListItem := "my_recipe"
	rli, err := NewRunListItem(runListItem)
	if err != nil {
		t.Errorf(`NewRunListItem("%s") did not correctly parse.`, runListItem)
	}
	if rli.Name != "my_recipe" {
		t.Errorf(`NewRunListItem("%s") did not correctly set name`, runListItem)
	}
	if rli.Type != "recipe" {
		t.Errorf(`NewRunListItem("%s") did not correctly set type`, runListItem)
	}
	if rli.Version != "" {
		t.Errorf(`NewRunListItem("%s") did not correctly set version`, runListItem)
	}
}

func TestBadCase(t *testing.T) {
	runListItem := "Recipe[my_recipe]" // Capital R
	_, err := NewRunListItem(runListItem)

	if err == nil {
		t.Errorf(`NewRunListItem("%s") should have returned and error and it didn't.`, runListItem)
	}
}

func TestSpaceThrowsError(t *testing.T) {
	runListItem := "recipe [my_recipe]"
	_, err := NewRunListItem(runListItem)

	if err == nil {
		t.Errorf(`NewRunListItem("%s") should have returned and error and it didn't.`, runListItem)
	}
}

func TestMissingClosingBracketThrowsError(t *testing.T) {
	runListItem := "recipe[my_recipe"
	_, err := NewRunListItem(runListItem)

	if err == nil {
		t.Errorf(`NewRunListItem("%s") should have returned and error and it didn't.`, runListItem)
	}
}

func TestStringConversion(t *testing.T) {
	runListItems := []string{
		"recipe[my_recipe@1.0.0]",
		"recipe[my_recipe]",
		"role[my_recipe]",
	}

	for _, runListItem := range runListItems {
		rli, err := NewRunListItem(runListItem)
		if err != nil || rli.String() != runListItem {
			t.Errorf(`NewRunListItem("%s").String() does not match %s`, runListItem, runListItem)
		}
	}
}

func TestIsRecipe(t *testing.T) {
	recipe := RunListItem{
		Type: "recipe",
	}
	if recipe.IsRecipe() != true {
		t.Error(`IsRecipe() should return true for recipe`)
	}

	role := RunListItem{
		Type: "role",
	}
	if role.IsRecipe() != false {
		t.Error(`IsRecipe() should return false for role`)
	}
}

func TestIsRole(t *testing.T) {
	recipe := RunListItem{
		Type: "role",
	}
	if recipe.IsRole() != true {
		t.Error(`IsRole() should return true for role`)
	}

	role := RunListItem{
		Type: "recipe",
	}
	if role.IsRole() != false {
		t.Error(`IsRole() should return false for recipe`)
	}
}
