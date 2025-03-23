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

func GetLocation(url string, cache *Cache) (LocationResponse, error) {
	if cacheData, found := cache.Get(url); found {
		fmt.Println("!!!!!!!!!!!!!!!Cache Hit!!!!!!!!!!!!!!!!!!!:", url)
		var locationResponse LocationResponse
		err := json.Unmarshal(cacheData, &locationResponse)
		if err != nil {
			return LocationResponse{}, errors.New("error unmarshalling JSON response")
		}
		return locationResponse, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return LocationResponse{}, fmt.Errorf("error in retreiving data %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return LocationResponse{}, errors.New("received non-200 response")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationResponse{}, errors.New("error reading response body")
	}

	var locationResponse LocationResponse
	err = json.Unmarshal(data, &locationResponse)
	if err != nil {
		return LocationResponse{}, errors.New("error unmarshalling JSON response")
	}

	fmt.Println("!!!!!!!!!!!!!!!Cache Miss!!!!!!!!!!!!!!!!!!!:", url)
	cache.Add(url, data)

	return locationResponse, nil
}

func GetPokemonLocation(url string, cache *Cache) (PokemonEncounterResult, error) {
	if cacheData, found := cache.Get(url); found {
		fmt.Println("!!!!!!!!!!!!!!!Cache Hit!!!!!!!!!!!!!!!!!!!:", url)
		var pokemonEncounterResult PokemonEncounterResult
		err := json.Unmarshal(cacheData, &pokemonEncounterResult)
		if err != nil {
			return PokemonEncounterResult{}, errors.New("error unmarshalling JSON response")
		}
		return pokemonEncounterResult, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return PokemonEncounterResult{}, fmt.Errorf("error in retreiving data %v", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return PokemonEncounterResult{}, errors.New("error reading response body")
	}

	var location PokemonEncounterResult
	err = json.Unmarshal(data, &location)
	if err != nil {
		return PokemonEncounterResult{}, errors.New("error unmarshalling JSON response")
	}

	fmt.Println("!!!!!!!!!!!!!!!Cache Miss!!!!!!!!!!!!!!!!!!!:", url)
	cache.Add(url, data)

	return location, nil
}
