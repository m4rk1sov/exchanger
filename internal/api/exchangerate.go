package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const apiURL = "https://openexchangerates.org/api/latest.json?app_id=%s"

type Response struct {
	Rates     map[string]float64
	Timestamp int64
	Base      string
}

func FectchExchangeRate() (*Response, error) {
	appID := os.Getenv("APP_ID")
	if appID == "" {
		return nil, fmt.Errorf("APP_ID is not set")
	}
	
	url := fmt.Sprintf(apiURL, appID)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)
	
	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	
	return &res, nil
}

func SaveToCache(res *Response) error {
	bytes, _ := json.Marshal(res)
	err := os.WriteFile("cache.json", bytes, 644)
	if err != nil {
		return err
	}
	return nil
}

func LoadFromCache() (*Response, error) {
	bytes, err := os.ReadFile("cache.json")
	if err != nil {
		return nil, err
	}
	var res Response
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return nil, err
	}
	
	return &res, nil
}
