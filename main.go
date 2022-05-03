package main

import (
	"strings"

	g "github.com/AllenDang/giu"
)

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
			g.TabBar().TabItems(g.TabItem("Recipes").Layout(g.Row(
				g.Column(
					g.InputText(&state.CurFilter).Hint("What do you seek?").Size(200),
					g.ListBox("RecipeBox", recipeNames).OnChange(state.onRecipeSelect).Size(200, 400),
					g.Button("+ Add recipe").OnClick(func() {
						state.addRecipeState.Reset()
						state.CurView = "add"
					}),
				),
				state.getMainView(),
			)), g.TabItem("Plan").Layout(
				state.MealPlannerState.Render(),
			)),
			g.PrepareMsgbox(), // Just a slot/prep for displaying message boxes later in the app.
		)
	}
}

// GLOBAL_STATE_PTR is meant to reduce dumb state pointer passing in functions.
var GLOBAL_STATE_PTR *State

func main() {
	var AppState State
	GLOBAL_STATE_PTR = &AppState
	recipes, err := getRecipes()

	if err != nil {
		panic(err)
	}

	AppState.Recipes = recipes

	wnd := g.NewMasterWindow("Recipe Manager - 0.2", 800, 500, g.MasterWindowFlagsNotResizable)
	wnd.Run(getLoop(&AppState))
}
