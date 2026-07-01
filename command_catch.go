package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

type pokemonGeneric struct {
	Name           string `json:"name"`
	ID             int    `json:"id"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Forms          []struct {
		Name string `json:"name"`
	} `json:"forms"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("too many arguments...pick one Pokemon to catch")
	}
	//. adding cli argument to URL api endpoint request, argument should be a pokemon listed from the area
	// returned from the explore command function, which will then give us all the pokemon
	// that you could encounter in that specified area this function will then be able to catch.
	url := "https://pokeapi.co/api/v2/pokemon/" + args[0]

	//. check cache to if the constructed url above is present as a key value in the cache struct, if it is,
	// if check triggers and unpacks previously cached raw json response without having to make a new get()
	// request.
	cachedJson, ok := cfg.cache.Get(url)
	if ok {
		pokemon := pokemonGeneric{}
		err := json.Unmarshal(cachedJson, &pokemon)
		if err != nil {
			return err
		}
		fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
		roll := rand.Intn(pokemon.BaseExperience)
		if roll > 75 {
			fmt.Printf("%s has escaped!...\n", args[0])
		} else {
			cfg.caughtPokemon[args[0]] = pokemon
			fmt.Printf("%s has been caught!\n", args[0])

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

	pokemon := pokemonGeneric{}
	err = json.Unmarshal(bodyContent, &pokemon)
	if err != nil {
		return err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[0])
	roll := rand.Intn(pokemon.BaseExperience)
	if roll > 75 {
		fmt.Printf("%s has escaped!...\n", args[0])
	} else {
		cfg.caughtPokemon[args[0]] = pokemon
		fmt.Printf("%s has been caught!\n", args[0])

	}
	return nil
}
