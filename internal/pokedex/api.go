package pokedex

import (
	"encoding/json"
	"net"
	"net/http"
	"time"
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

func GetLocation(url ...string) (LocationAreaResponse, error) {
	if len(url) == 0 || url[0] == "" {
		url = make([]string, 0, 1)
		url = append(url, locationAreaUrl)
	}

	req, err := http.NewRequest("GET", url[0], nil)
	if err != nil {
		return LocationAreaResponse{}, err
	}

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
	res, err := client.Do(req)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close()

	var locationArea LocationAreaResponse
	decoder := json.NewDecoder(res.Body)
	if err = decoder.Decode(&locationArea); err != nil {
		return LocationAreaResponse{}, err
	}

	return locationArea, nil
}
