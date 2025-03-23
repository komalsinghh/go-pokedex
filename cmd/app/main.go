package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	httppokedex "github.com/komalsinghh/go-pokedex/internal/pokecache"
)

type Config struct {
	Next            string
	Previous        string
	Cache           *httppokedex.Cache
	ExploreLocation string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config)
}

var commandsMap map[string]cliCommand

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := &Config{
		Next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		Previous: "",
		Cache:    httppokedex.NewCache(10 * time.Second),
	}
	commandsMap = map[string]cliCommand{
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
	}
	for {
		fmt.Print("Pokedex> ")
		if scanner.Scan() {
			text := scanner.Text()
			if len(text) == 0 {
				continue
			}
			part := strings.Fields(text)
			command := part[0]
			if cmd, found := commandsMap[command]; found {
				if command == "explore" && len(part) > 1 {
					config.ExploreLocation = part[1]
				}
				cmd.callback(config)
			} else {
				fmt.Println("Unknown command. Type 'help' for available commands.")
			}
		}
	}
}

func mapLocation(config *Config) {
	if config.Next == "" {
		fmt.Println("No more locations available.")
		return
	}

	locationResponse, err := httppokedex.GetLocation(config.Next, config.Cache)
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

func mapPreviousLocation(config *Config) {
	if config.Previous == "" {
		fmt.Println("You're on the first page.")
		return
	}

	locationResponse, err := httppokedex.GetLocation(config.Previous, config.Cache)
	if err != nil {
		fmt.Println("Error fetching previous locations:", err)
		return
	}

	fmt.Println("Location Areas:")
	for _, loc := range locationResponse.Results {
		fmt.Println("-", loc.Name)
	}

	config.Next = locationResponse.Next
	config.Previous = locationResponse.Previous
}

func explorePokemonLocation(config *Config) {
	location := config.ExploreLocation
	config.ExploreLocation = fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)
	fmt.Println("url----->", config.ExploreLocation)
	locationResponse, err := httppokedex.GetPokemonLocation(config.ExploreLocation)
	if err != nil {
		fmt.Println("Error fetching previous locations:", err)
		return
	}

	fmt.Println("Location Areas of Pokemon:")
	for _, loc := range locationResponse.PokemonEncounter {
		fmt.Println("-", loc.Pokemon.Name)
	}
}

func commandExit(config *Config) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func commandClearScreen(config *Config) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func commandHelp(config *Config) {
	fmt.Println("Welcome to the Pokedex!\nUsage:k")
	for key, value := range commandsMap {
		fmt.Println(key, " ", value.description)
	}
}
