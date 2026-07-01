package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type locationAreaEncounters struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("too many locations/arguements...pick one area to explore only")
	}
	//. adding cli argument to URL api endpoint request, argument should be an area listed from the map
	// function, which will then give us all the pokemon that you could encounter in that specified area
	url := "https://pokeapi.co/api/v2/location-area/" + args[0]

	//. check cache to if the constructed url above is present as a key value in the cache struct, if it is,
	// if check triggers and unpacks previously cached raw json response without having to make a new get()
	// request.
	cachedJson, ok := cfg.cache.Get(url)
	if ok {
		encounters := locationAreaEncounters{}
		err := json.Unmarshal(cachedJson, &encounters)
		if err != nil {
			return err
		}
		fmt.Printf("Exploring %s...\n", args[0])
		fmt.Println("Found Pokemon:")
		for _, encounter := range encounters.PokemonEncounters {
			fmt.Printf(" - %s\n", encounter.Pokemon.Name)
		}
		return nil
	}
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bodyContent, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	//add url to cache for caching functionality to take place, stores raw json response as a slice []byte
	cfg.cache.Add(url, bodyContent)

	encounters := locationAreaEncounters{}
	err = json.Unmarshal(bodyContent, &encounters)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", args[0])
	fmt.Println("Found Pokemon:")
	for _, encounter := range encounters.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}
