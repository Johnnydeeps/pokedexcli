package main

import (
	"time"

	"github.com/Johnnydeeps/pokedexcli/internal/pokecache"
)

func main() {
	config := config{
		cache: pokecache.NewCache(5 * time.Minute),
	}
	cfg := &config
	startRepl(cfg)
}
