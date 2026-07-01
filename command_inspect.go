package main

import "fmt"

func commandInspect(cfg *config, args ...string) error {
	pokemon, ok := cfg.caughtPokemon[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("ID: %d\n", pokemon.ID)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Printf("Forms: %s\n", pokemon.Forms)
		fmt.Printf("Stats: %v\n", pokemon.Stats)
		fmt.Printf("Types: %v\n", pokemon.Types)
	}
	return nil
}
