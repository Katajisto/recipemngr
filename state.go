package main

import (
	g "github.com/AllenDang/giu"
	"math/rand"
)

type Day struct {
	Breakfast    Recipe
	Lunch        Recipe
	LunchLocked  bool
	Dinner       Recipe
	DinnerLocked bool
}

func (d *Day) Randomize() {
	if !d.LunchLocked {
		d.Lunch = GLOBAL_STATE_PTR.getRandomRecipe()
	}
	if !d.DinnerLocked {
		d.Dinner = GLOBAL_STATE_PTR.getRandomRecipe()
	}
}

func (d *Day) Render(day string) g.Widget {
	return g.Column(
		LabelWithSize(day, 40),
		g.Row(
			LabelWithSize(d.Lunch.Name, 30),
			g.Button(CTF(d.LunchLocked, "Unlock", "Lock")).OnClick(func() {
				d.LunchLocked = !d.LunchLocked
			}),
			g.Button("Clear").OnClick(func() {
				d.Lunch = Recipe{Name: "[EMPTY]", Empty: true}
			}),
		),
		g.Row(
			LabelWithSize(d.Dinner.Name, 30),
			g.Button(CTF(d.DinnerLocked, "Unlock", "Lock")).OnClick(func() {
				d.DinnerLocked = !d.DinnerLocked
			}),
			g.Button("Clear").OnClick(func() {
				d.Dinner = Recipe{Name: "[EMPTY]", Empty: true}
			}),
		),
	)
}

type MealPlannerState struct {
	Monday    Day
	Tuesday   Day
	Wednesday Day
	Thursday  Day
	Friday    Day
	Saturday  Day
	Sunday    Day
}

func (mps *MealPlannerState) Render() g.Widget {
	return g.Column(
		g.Row(
			LabelWithSize("Meal Planner", 50),
			ButtonWithSize(
				"Randomize",
				40,
				func() {
					mps.Randomize()
				},
			),
		),
		mps.Monday.Render("Monday"),
		mps.Tuesday.Render("Tuesday"),
		mps.Wednesday.Render("Wednesday"),
		mps.Thursday.Render("Thursday"),
		mps.Friday.Render("Friday"),
		mps.Saturday.Render("Saturday"),
		mps.Sunday.Render("Sunday"),
	)
}

func (mps *MealPlannerState) Randomize() {
	mps.Monday.Randomize()
	mps.Tuesday.Randomize()
	mps.Wednesday.Randomize()
	mps.Thursday.Randomize()
	mps.Friday.Randomize()
	mps.Saturday.Randomize()
	mps.Sunday.Randomize()
}

type State struct {
	MealPlannerState MealPlannerState
	addRecipeState   AdderState
	Recipes          []Recipe
	SelectedRecipe   int
	CurFilter        string
	CurView          string
}

func (state *State) onRecipeSelect(selectedIndex int) {
	state.SelectedRecipe = selectedIndex
	state.CurView = "recipe"
}

func (state *State) curRecipe() *Recipe {
	return &state.Recipes[state.SelectedRecipe]
}

func (state *State) getRandomRecipe() Recipe {
	recipeCount := len(state.Recipes)
	randIndex := rand.Intn(recipeCount)
	return state.Recipes[randIndex]
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
