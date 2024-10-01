package externalapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type ExternalApiClient struct {
	ApiURL string
	logger *logrus.Logger
}

func NewExternalApiClient(apiURL string, logger *logrus.Logger) *ExternalApiClient {
	return &ExternalApiClient{
		ApiURL: apiURL,
		logger: logger,
	}
}

type response struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (e *ExternalApiClient) FetchSongInfo(group, song string) (*response, error) {
	url := fmt.Sprintf("%sgroup=%s&song=%s", e.ApiURL, group, song)
	e.logger.WithFields(logrus.Fields{
		"url":   url,
		"group": group,
		"song":  song,
	}).Info("Fetching song info from external API")

	resp, err := http.Get(url)
	if err != nil {
		e.logger.WithFields(logrus.Fields{
			"url":   url,
			"group": group,
			"song":  song,
		}).Error("Error while making HTTP request to external API: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e.logger.WithFields(logrus.Fields{
			"url":        url,
			"group":      group,
			"song":       song,
			"statusCode": resp.StatusCode,
		}).Errorf("Failed to fetch song info: status code %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to fetch song info: status code %d", resp.StatusCode)
	}

	e.logger.WithFields(logrus.Fields{
		"url":        url,
		"group":      group,
		"song":       song,
		"statusCode": resp.StatusCode,
	}).Info("Successfully fetched song info from external API")

	var response response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		e.logger.WithFields(logrus.Fields{
			"url":   url,
			"group": group,
			"song":  song,
			"error": err,
		}).Error("Failed to decode response from external API")
		return nil, err
	}

	e.logger.WithFields(logrus.Fields{
		"group": group,
		"song":  song,
	}).Info("Successfully decoded song info")

	return &response, nil
}
