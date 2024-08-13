package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BradMyrick/pokedex/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Results []LocationArea `json:"results"`
}

type PokemonMap struct {
	index            int
	locationsPerPage int
	locationAreas    []LocationArea
}

type LocationAreaDetails struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type Pokedex struct {
	caughtPokemon map[string]Pokemon
}

var locationMap *PokemonMap
var cache = pokecache.NewCache(5 * time.Minute)
var pokedex *Pokedex

func newPokemonMap() *PokemonMap {
	return &PokemonMap{
		index:            0,
		locationsPerPage: 20,
		locationAreas:    []LocationArea{},
	}
}

func commandMap(args []string) error {
	if locationMap == nil {
		locationMap = newPokemonMap()
		if err := locationMap.fetchLocationAreas(); err != nil {
			return err
		}

		for location := range locationMap.locationAreas {
			fmt.Println(locationMap.locationAreas[location].Name)
		}
		locationMap.index += locationMap.locationsPerPage
		return nil
	}

	if err := locationMap.fetchLocationAreas(); err != nil {
		return err
	}

	for location := range locationMap.locationAreas {
		fmt.Println(locationMap.locationAreas[location].Name)
	}

	locationMap.index += locationMap.locationsPerPage

	return nil
}

func (pm *PokemonMap) fetchLocationAreas() error {
	url := "https://pokeapi.co/api/v2/location-area/?offset=" + fmt.Sprintf("%d", pm.index) + "&limit=" + fmt.Sprintf("%d", pm.locationsPerPage)

	if data, ok := cache.Get(url); ok {
		var locationAreaResp LocationAreaResponse
		err := json.Unmarshal(data, &locationAreaResp)
		if err != nil {
			return err
		}
		pm.locationAreas = locationAreaResp.Results
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var locationAreaResp LocationAreaResponse
	err = json.NewDecoder(resp.Body).Decode(&locationAreaResp)
	if err != nil {
		return err
	}

	data, err := json.Marshal(locationAreaResp)
	if err != nil {
		return err
	}
	cache.Add(url, data)

	pm.locationAreas = locationAreaResp.Results
	return nil
}

func commandExplore(areaName string) error {
	if areaName == "" {
		return fmt.Errorf("please provide a location area name")
	}

	fmt.Printf("Exploring %s...\n", areaName)

	url := "https://pokeapi.co/api/v2/location-area/" + areaName

	var data []byte
	var ok bool

	if data, ok = cache.Get(url); !ok {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cache.Add(url, data)
	}

	var details LocationAreaDetails
	err := json.Unmarshal(data, &details)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range details.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(pokemonName string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	var data []byte
	var ok bool

	if data, ok = cache.Get(url); !ok {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cache.Add(url, data)
	}

	var pokemonData Pokemon
	err := json.Unmarshal(data, &pokemonData)
	if err != nil {
		return err
	}

	catchChance := 100 - pokemonData.BaseExperience/2
	if catchChance < 1 {
		catchChance = 1
	}

    if rand.Intn(100) < catchChance {
        pokedex.caughtPokemon[pokemonName] = pokemonData
        fmt.Printf("%s was caught!\n", pokemonName)
        fmt.Println("You may now inspect it with the inspect command.")
    } else {
        fmt.Printf("%s escaped!\n", pokemonName)
    }

    return nil
}

func getCommands() map[string]cliCommand {
	commandHelp := func(args []string) error {
		fmt.Println("Available commands:")
		for name, command := range getCommands() {
			fmt.Printf("%s: %s\n", name, command.description)
		}
		return nil
	}

	commandExit := func(args []string) error {
		os.Exit(0)
		return nil
	}

	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "View a map of the pokemon locations",
			callback:    commandMap,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area and list the Pokemon found",
			callback: func(args []string) error {
				if len(args) < 2 {
					return fmt.Errorf("please provide a location area name")
				}
				return commandExplore(args[1])
			},
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon by name",
			callback: func(args []string) error {
				if len(args) < 2 {
					return fmt.Errorf("please provide a Pokemon name")
				}
				return commandCatch(args[1])
			},
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon",
			callback: func(args []string) error {
				if len(args) < 2 {
					return fmt.Errorf("please provide a Pokemon name")
				}
				return commandInspect(args[1])
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show a list of caught Pokemon",
			callback: func(args []string) error {
				return commandPokedex()
			},
		},
	}
}

func main() {
	pokedex = &Pokedex{
		caughtPokemon: make(map[string]Pokemon),
	}

	commands := getCommands()

	fmt.Println("Welcome to the Pokedex!")

	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")

		args := strings.Split(text, " ")
		commandName := args[0]

		if command, ok := commands[commandName]; ok {
			fmt.Printf("Running command: %s\n", command.name)

			if err := command.callback(args); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func commandInspect(pokemonName string) error {
	if pokemonName == "" {
		return fmt.Errorf("please provide a Pokemon name")
	}

	pokemon, ok := pokedex.caughtPokemon[pokemonName]
	if !ok {
		fmt.Printf("You have not caught %s\n", pokemonName)
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typ := range pokemon.Types {
		fmt.Printf("  - %s\n", typ.Type.Name)
	}

	return nil
}

func commandPokedex() error {
	fmt.Println("Your Pokedex:")
	if len(pokedex.caughtPokemon) == 0 {
		fmt.Println("  You haven't caught any Pokemon yet.")
		return nil
	}

	for _, pokemon := range pokedex.caughtPokemon {
		fmt.Printf("  - %s\n", pokemon.Name)
	}
	return nil
}
