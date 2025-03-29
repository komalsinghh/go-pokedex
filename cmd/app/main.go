package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	httppokedex "github.com/komalsinghh/go-pokedex/internal/pokecache"
)

type Config struct {
	Next     string
	Previous string
	Cache    *httppokedex.Cache
	Pokedex  map[string]httppokedex.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string)
}

var commandsMap map[string]cliCommand

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := &Config{
		Next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		Previous: "",
		Cache:    httppokedex.NewCache(10 * time.Second),
		Pokedex:  make(map[string]httppokedex.Pokemon),
	}

	for {
		fmt.Print("Pokedex> ")
		if scanner.Scan() {
			text := scanner.Text()
			if len(text) == 0 {
				continue
			}
			part := cleanInput(text)
			commandName := part[0]
			args := []string{}
			if len(part) > 1 {
				args = part[1:]
			}
			command, exists := getCommand()[commandName]
			if exists {
				command.callback(config, args...)
				continue
			} else {
				fmt.Println("Unknown command. Type 'help' for available commands.")
			}
		}
	}
}

func getCommand() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"clear": {
			name:        "clear",
			description: "Clear the Pokedex terminal",
			callback:    commandClearScreen,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map":     {"map", "Displays the next 20 locations", mapLocation},
		"mapb":    {"mapb", "Displays the previous 20 locations", mapPreviousLocation},
		"explore": {"explore", "Displays list of all the Pokémon located there", explorePokemonLocation},
		"catch":   {"catch", "Try to catch a Pokémon", catchPokemon},
		"pokedex": {"pokedex", "Displays your caught Pokémon", showPokedex},
	}
}

func cleanInput(input string) []string {
	input = strings.ToLower(input)
	inputArr := strings.Fields(input)
	return inputArr
}

func mapLocation(config *Config, args ...string) {
	fetchAndDisplayLocations(config.Next, config)
}

func mapPreviousLocation(config *Config, args ...string) {
	fetchAndDisplayLocations(config.Previous, config)
}

func fetchAndDisplayLocations(url string, config *Config) {
	if url == "" {
		fmt.Println("No more locations available.")
		return
	}

	locationResponse, err := httppokedex.GetLocation(url, config.Cache)
	if err != nil {
		fmt.Println("Error fetching locations:", err)
		return
	}

	fmt.Println("Location Areas:")
	for _, loc := range locationResponse.Results {
		fmt.Println("-", loc.Name)
	}

	config.Next = locationResponse.Next
	config.Previous = locationResponse.Previous
}

func explorePokemonLocation(config *Config, args ...string) {
	if len(args) != 1 {
		fmt.Println("you must provide a location name")
		return
	}

	name := args[0]
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", name)
	fmt.Println("url----->", url)
	pokemonEncounterResult, err := httppokedex.GetPokemonLocation(url, config.Cache)
	if err != nil {
		fmt.Println("Error fetching locations of Pokemon:", err)
		return
	}

	fmt.Println("Found Pokemon...")
	for _, res := range pokemonEncounterResult.PokemonEncounter {
		fmt.Println("-", res.Pokemon.Name)
	}
}
func catchPokemon(config *Config, args ...string) {
	if len(args) != 1 {
		fmt.Println("you must provide a pokemon name")
		return
	}

	name := args[0]

	if _, exists := config.Pokedex[name]; exists {
		fmt.Printf("%v already caught", name)
		return
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	fmt.Println(url)
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	catchPokemon, err := httppokedex.GetCatchPokemon(url)
	if err != nil {
		fmt.Println("An error occured ", err)
		return
	}

	res := rand.Intn(catchPokemon.BaseExperience)
	if res > 40 {
		fmt.Printf("%s escaped!\n", name)
	} else {
		fmt.Printf("%s was caught!\n", name)
		config.Pokedex[name] = httppokedex.Pokemon{Name: name}
	}
}

func showPokedex(config *Config, args ...string) {
	fmt.Println("Your Pokedex:")
	for _, value := range config.Pokedex {
		fmt.Printf(" - %v\n", value.Name)
	}
}

func commandExit(config *Config, args ...string) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func commandClearScreen(config *Config, args ...string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func commandHelp(config *Config, args ...string) {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for key, value := range commandsMap {
		fmt.Println(key, " ", value.description)
	}
}
