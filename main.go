package main

import (
	"fmt"

	"github.com/englandrecoil/go-pokedex-cli/internal/api"
)

func main() {
	pokemon := api.GetPokemon()
	fmt.Println(pokemon)
}
