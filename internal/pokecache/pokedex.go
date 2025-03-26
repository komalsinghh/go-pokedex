package pokecache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type LocationResponse struct {
	Results  []Location `json:"results"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
}

type Location struct {
	Name string `json:"name"`
}

type PokemonEncounterResult struct {
	PokemonEncounter []Pokemon `json:"pokemon_encounters"`
}

type Pokemon struct {
	Pokemon Location `json:"pokemon"`
}

type CatchPokemon struct {
	BaseExperience int `json:"base_experience"`
}

func GetLocation(url string, cache *Cache) (LocationResponse, error) {
	var locationResponse LocationResponse
	err := fetchAndCache(url, cache, &locationResponse)
	if err != nil {
		return LocationResponse{}, err
	}
	return locationResponse, nil
}

func GetPokemonLocation(url string, cache *Cache) (PokemonEncounterResult, error) {
	var pokemonEncounterResult PokemonEncounterResult
	err := fetchAndCache(url, cache, &pokemonEncounterResult)
	if err != nil {
		return PokemonEncounterResult{}, err
	}
	return pokemonEncounterResult, nil
}

func fetchAndCache[T any](url string, cache *Cache, target *T) error {
	if cacheData, found := cache.Get(url); found {
		fmt.Println("!!!!!!!!!!!!!!! Cache Hit !!!!!!!!!!!!!!!!!!!:", url)
		err := json.Unmarshal(cacheData, target)
		if err != nil {
			return errors.New("error unmarshalling cached data")
		}
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error retrieving data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("error reading response body")
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		return errors.New("error unmarshalling JSON response")
	}

	fmt.Println("!!!!!!!!!!!!!!! Cache Miss !!!!!!!!!!!!!!!!!!!:", url)
	cache.Add(url, data)
	return nil
}

func GetCatchPokemon(url string) (CatchPokemon, error) {
	resp, err := http.Get(url)
	if err != nil {
		return CatchPokemon{}, fmt.Errorf("error retrieving data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CatchPokemon{}, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return CatchPokemon{}, errors.New("error reading response body")
	}

	var catchPokemon CatchPokemon
	err = json.Unmarshal(data, &catchPokemon)
	if err != nil {
		return CatchPokemon{}, errors.New("error unmarshalling JSON response")
	}

	return catchPokemon, nil
}
