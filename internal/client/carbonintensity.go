package client

import (
	"app/internal/client/converter"
	"app/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) FetchCarbonForecast(ctx context.Context) ([]domain.Carbon, error) {
	currentTime := time.Now().UTC().Format("2006-01-02T15:04") + "Z"
	requestURL := fmt.Sprintf("%s/intensity/%s/fw24h", c.cfg.CarbonClient.URL, currentTime)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		c.log.Errorf("Failed to create request: %v", err)
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.log.Errorf("Failed to fetch carbon intensity data: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiError struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			c.log.Errorf("API responded with status %d but failed to decode error message: %v", resp.StatusCode, err)
			return nil, err
		}

		c.log.Errorf("API error: %s - %s", apiError.Error.Code, apiError.Error.Message)
		return nil, errors.New(apiError.Error.Message)
	}

	var data converter.CarbonIntensityResponse
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.log.Errorf("Failed to decode API response: %v", err)
		return nil, err
	}

	return converter.ToDomain(data), nil
}
