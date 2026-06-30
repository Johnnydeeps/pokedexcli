package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type locationAreasResponse struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
	} `json:"results"`
}

func commandMap(cfg *config) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	//.  "Default to page 1. But if cfg already has a next-page URL stored from a previous call,
	// use that instead." The *cfg.nextLocationsURL dereferences the pointer to get the actual string.
	if cfg.nextLocationsURL != nil {
		// this is where the string in memory is deferenced (*)
		url = *cfg.nextLocationsURL
	}

	// ***********************************************************************************************
	//. {
	// "count": 1089,
	//"next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
	//"previous": null,
	//"results": [ ... 20 items ... ]
	//}
	// this is the first/ intitial response from the PokeApi server from the first get request,
	// which contains a json (snapshot) representation of the resource on the server.
	// ie. .../location-area. this is resource would be considered RESTful because the resource
	// corresponds with a URL.
	// 		Additionally, the response json contains pagination information which
	// is contained in the Next and Previous fields of the locationAreasResponse struct. These are
	// accessed through the parse stage/unmarshal on the response.Body below and stored via
	// cfg.nextLocationsURL = locations.Next. cfg is a config struct with pointers which are replaced
	// with the ones from the response.Body. note the *string pointer type match in the base location
	// struct above. this allows the next map function call in the CLI to use the query parts of the
	// returned url in the Next field to bookmark the position of the pages in the cfg struct.
	// ************************************ IMPORTANT SUMMARY*****************************************
	// The URL string is allocated by json.Unmarshal and lives in heap memory.
	// locations.Next and cfg.nextLocationsURL are both *string pointers holding that string's address,
	// not the string itself. *cfg.nextLocationsURL(above at the start oif commandMap) dereferences
	// the address to read the actual text.
	// The string outlives the local locations variable because cfg still points to it and go's garbage
	// collector keeps it alive after the function call ends.
	// Unmarshal allocates the string once; locations.Next and cfg.nextLocationsURL
	// are two pointers to that same address. Dereferencing (*) reads the value back out.
	// ***********************************************************************************************
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bodyContent, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	locations := locationAreasResponse{}
	err = json.Unmarshal(bodyContent, &locations)
	if err != nil {
		return err
	}

	//. write new urls back into cfg *config, to save the "page" that you are on. this is the
	// feedback loop that is stored in the config struct in repl.go, and called in main.go
	// with cfg := &config{}
	cfg.nextLocationsURL = locations.Next
	cfg.prevLocationsURL = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(cfg *config) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("you're on the first page")
	}

	url := *cfg.prevLocationsURL

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	bodyContent, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	locations := locationAreasResponse{}
	err = json.Unmarshal(bodyContent, &locations)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locations.Next
	cfg.prevLocationsURL = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}
