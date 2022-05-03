package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Ingredient struct {
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
}

type Recipe struct {
	Name         string       `json:"name"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions string       `json:"instructions"`
	Empty        bool
}

func getRecipes() ([]Recipe, error) {
	res, err := filepath.Glob("./recipes/*.nam")
	if err != nil {
		return nil, err
	}

	var recipes []Recipe

	for _, recipePath := range res {
		file, err := os.Open(recipePath)

		if err != nil {
			log.Println("Failed to open recipe: ", recipePath, " - ", err)
			continue
		}

		var recipe Recipe

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&recipe)
		if err != nil {
			log.Println("Failed to JSON parse nam file: ", recipePath, " - ", err)
			continue
		}

		recipes = append(recipes, recipe)

		err = file.Close()
		if err != nil {
			log.Println("Failed to close nam file: ", recipePath, " - ", err)
			continue
		}
	}

	return recipes, nil
}
