package pokedex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/paysis/pokedex/internal/pokecache"
)

type LocationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationSingleAreaResponse struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type LocationApi struct {
	cache *pokecache.Cache
}

func NewLocationApi() *LocationApi {
	locApi := &LocationApi{
		cache: pokecache.NewCache(1 * time.Hour),
	}
	return locApi
}

func (la *LocationApi) GetLocationArea(url ...string) (LocationAreaResponse, error) {
	if len(url) == 0 || url[0] == "" {
		url = make([]string, 0, 1)
		url = append(url, locationAreaUrl)
	}

	// cache retrieve
	buf, ok := la.cache.Get(url[0])
	if ok {
		resp, err := DecodeJson[LocationAreaResponse](buf)
		if err != nil {
			return LocationAreaResponse{}, err
		}
		return resp, nil
	}

	req, err := http.NewRequest("GET", url[0], nil)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	client := createClient()
	res, err := client.Do(req)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close()

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}

	locationArea, err := DecodeJson[LocationAreaResponse](bytes.NewBuffer(bodyContent))
	if err != nil {
		return LocationAreaResponse{}, err
	}

	// cache
	if err == nil {
		la.cache.Add(url[0], bytes.NewBuffer(bodyContent))
	} else {
		fmt.Printf("cache is not working properly: %v", err)
	}

	return locationArea, nil
}

func (la *LocationApi) GetSingleLocationArea(name string) (LocationSingleAreaResponse, error) {
	var url string = locationAreaUrl + "/" + name

	// cache retrieve
	buf, ok := la.cache.Get(url)
	if ok {
		resp, err := DecodeJson[LocationSingleAreaResponse](buf)
		if err != nil {
			return LocationSingleAreaResponse{}, err
		}

		return resp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationSingleAreaResponse{}, err
	}

	client := createClient()
	res, err := client.Do(req)
	if err != nil {
		return LocationSingleAreaResponse{}, err
	}
	defer res.Body.Close()

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationSingleAreaResponse{}, err
	}

	locationSingleArea, err := DecodeJson[LocationSingleAreaResponse](bytes.NewBuffer(bodyContent))
	if err != nil {
		return LocationSingleAreaResponse{}, err
	}

	// cache
	if err == nil {
		la.cache.Add(url, bytes.NewBuffer(bodyContent))
	} else {
		fmt.Printf("cache is not working properly: %v", err)
	}

	return locationSingleArea, nil
}

func createClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Second * 5,
			}).DialContext,
			Dial:                nil,
			TLSHandshakeTimeout: time.Second * 5,
		},
		Timeout: time.Second * 10,
	}
	return client
}

func DecodeJson[T any](jsonData *bytes.Buffer) (output T, err error) {
	decoder := json.NewDecoder(jsonData)
	err = decoder.Decode(&output)
	return
}
