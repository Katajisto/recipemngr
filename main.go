package main

import (
	"log"

	g "github.com/AllenDang/giu"
)

func (s *State) onRecipeSelect(selectedIndex int) {
	s.SelectedRecipe = selectedIndex
}

func (s *State) curRecipe() *Recipe {
	return &s.Recipes[s.SelectedRecipe]
}

func getLoop(state *State) func() {
	return func() {
		var recipeNames []string
		var firstRecipe string
		for i, r := range state.Recipes {
			if i == 0 {
				firstRecipe = r.Name
			}
			recipeNames = append(recipeNames, r.Name)
		}

		g.SingleWindow().Layout(

			g.Row(
				g.ListBox(firstRecipe, recipeNames).OnChange(state.onRecipeSelect).Size(200, g.Auto),
				g.Column(
					g.Style().SetFontSize(50).To(g.Label(state.Recipes[state.SelectedRecipe].Name)),
					g.Label(state.curRecipe().Instructions).Wrapped(true),
				),
			),
		)
	}
}

type State struct {
	Recipes        []Recipe
	SelectedRecipe int
}

func main() {
	var AppState State
	recipes, err := getRecipes()
	log.Println(recipes, err)
	AppState.Recipes = recipes

	wnd := g.NewMasterWindow("Recipe Manager - 0.1", 800, 500, g.MasterWindowFlagsNotResizable)
	wnd.Run(getLoop(&AppState))
}
