package main

import (
	"encoding/json"
	"errors"
	"fmt"
	g "github.com/AllenDang/giu"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

/*
	This contains everything needed for the recipe add state and rendering the
	recipe adder.
*/

// AdderState contains all the state that the add recipe view needs, and also
// implements methods for resetting the state and rendering it.
type AdderState struct {
	Name          string
	Ingredients   []Ingredient
	CurIngredient Ingredient
	Instructions  string
}

func (i Ingredient) Render() g.Widget {
	return g.Row(
		g.Label(i.Name),
		g.Label(fmt.Sprintf("%v", i.Amount)),
	)
}

func RenderIngredients(list []Ingredient) g.Widget {
	var ingredientWidgets []g.Widget
	for _, i := range list {
		ingredientWidgets = append(ingredientWidgets, i.Render())
	}

	return g.Column(ingredientWidgets...)
}

// Reset makes state empty and ready for another add.
func (as *AdderState) Reset() {
	as.Name = ""
	as.Ingredients = make([]Ingredient, 0)
	as.Instructions = ""
}

// GetRecipe forms adder state into a recipe struct.
func (as AdderState) GetRecipe() Recipe {
	return Recipe{Name: as.Name, Instructions: as.Instructions, Ingredients: as.Ingredients}
}

func (as *AdderState) AddRecipe(state *State) {
	// Don't add the recipe if name is empty.
	if len(strings.TrimSpace(as.Name)) == 0 {
		g.Msgbox("Error!", "Name can't be empty.")
		return
	}

	recipeToAdd := as.GetRecipe()
	recipePath := "./recipes/" + recipeToAdd.Name + ".nam"
	state.Recipes = append(state.Recipes, recipeToAdd)

	// Save the recipe to disk as .nam file (is just JSON)
	if _, err := os.Stat(recipePath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(recipePath)
		if err != nil {
			g.Msgbox("Error!", fmt.Sprintf("Error saving your recipe: %v", err))
			return
		}

		jsondata, err := json.MarshalIndent(recipeToAdd, "", "   ")
		if err != nil {
			g.Msgbox("Error!", fmt.Sprintf("Error saving your recipe: %v", err))
			return
		}

		_, err = file.Write(jsondata)
		if err != nil {
			g.Msgbox("Error!", fmt.Sprintf("Error saving your recipe: %v", err))
			return
		}
	} else {
		g.Msgbox("Error!", "There already is a recipe with the same filename. Try another one.")
		return
	}

	as.Reset()
	state.CurView = ""

	// HACK: So the package Giu does not export the selectedIndex field. This is a problem
	// as we need to change it. This is why we do some really hacky and unsafe stuff here.
	// I hope you can forgive me.
	var listState g.ListBoxState
	listStateRef := reflect.ValueOf(&listState).Elem()
	// This can't be changed.
	indexPtr := listStateRef.Field(0)
	// Now it can, we use unsafe black magic to get a pointer that doesn't complain about it.
	indexPtr = reflect.NewAt(indexPtr.Type(), unsafe.Pointer(indexPtr.UnsafeAddr())).Elem()

	// Now that we can change things, we set the selected value from the list to the newly added
	// recipe and do the same for our program state.
	indexPtr.Set(reflect.ValueOf(len(state.Recipes) - 1))
	state.SelectedRecipe = len(state.Recipes) - 1
	g.Context.SetState("RecipeBox", &listState)
}

// RenderAdder returns a view for adding a recipe.
func (as *AdderState) RenderAdder(state *State) g.Widget {
	return g.Column(
		g.Style().SetFontSize(40).To(g.Label("Add recipe").Wrapped(true)),
		g.InputText(&as.Name).Hint("Recipe name"),
		RenderIngredients(as.Ingredients),
		g.Row(
			g.InputText(&as.CurIngredient.Name).Size(320).Hint("Ingredient name"),
			g.InputFloat(&as.CurIngredient.Amount).Size(100),
			g.Button("Add").Size(84, 20).OnClick(func() {
				as.Ingredients = append(as.Ingredients, as.CurIngredient)
				as.CurIngredient = Ingredient{}
			}),
		),
		g.InputTextMultiline(&as.Instructions),
		g.Button("Submit").OnClick(func() { as.AddRecipe(state) }),
	)
}
