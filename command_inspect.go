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

		fmt.Println("Forms:")
		for _, form := range pokemon.Forms {
			fmt.Printf("  -%s\n", form.Name)
		}
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  -%s\n", t.Type.Name)
		}
	}
	return nil
}
