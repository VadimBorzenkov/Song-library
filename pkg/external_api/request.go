package externalapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ExternalApiClient struct {
	ApiURL string
}

func NewExternalApiClient(apiURL string) *ExternalApiClient {
	return &ExternalApiClient{ApiURL: apiURL}
}

type response struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (e *ExternalApiClient) FetchSongInfo(group, song string) (*response, error) {
	url := fmt.Sprintf("%sgroup=%s&song=%s", e.ApiURL, group, song)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch song info: status code %d", resp.StatusCode)
	}

	var response response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
