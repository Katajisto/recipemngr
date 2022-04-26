package main

import (
	"log"
	"strings"

	g "github.com/AllenDang/giu"
)

func (state *State) onRecipeSelect(selectedIndex int) {
	state.SelectedRecipe = selectedIndex
	state.CurView = "recipe"
}

func (state *State) curRecipe() *Recipe {
	return &state.Recipes[state.SelectedRecipe]
}

func (state *State) getMainView() g.Widget {
	switch {
	case state.CurView == "recipe" || state.CurView == "":
		return g.Column(
			g.Style().SetFontSize(50).To(g.Label(state.Recipes[state.SelectedRecipe].Name).Wrapped(true)),
			g.Label(state.curRecipe().Instructions).Wrapped(true),
		)
	case state.CurView == "add":
		return state.addRecipeState.RenderAdder(state)
	}

	return g.Label("")
}

func getLoop(state *State) func() {
	return func() {
		var recipeNames []string
		for _, r := range state.Recipes {
			// Only show recipes that contain the filter.
			if len(state.CurFilter) != 0 && !strings.Contains(strings.ToLower(r.Name), strings.ToLower(state.CurFilter)) {
				continue
			}

			recipeNames = append(recipeNames, r.Name)
		}

		g.SingleWindow().Layout(

			g.Row(
				g.Column(
					g.InputText(&state.CurFilter).Hint("What do you seek?").Size(200),
					g.ListBox("RecipeBox", recipeNames).OnChange(state.onRecipeSelect).Size(200, 420),
					g.Button("+ Add recipe").OnClick(func() {
						state.addRecipeState.Reset()
						state.CurView = "add"
					}),
				),
				state.getMainView(),
			),
			g.PrepareMsgbox(), // Just a slot/prep for displaying message boxes later in the app.
		)
	}
}

type State struct {
	addRecipeState AdderState
	Recipes        []Recipe
	SelectedRecipe int
	CurFilter      string
	CurView        string
}

func main() {
	var AppState State
	recipes, err := getRecipes()
	log.Println(recipes, err)
	AppState.Recipes = recipes

	wnd := g.NewMasterWindow("Recipe Manager - 0.2", 800, 500, g.MasterWindowFlagsNotResizable)
	wnd.Run(getLoop(&AppState))
}
