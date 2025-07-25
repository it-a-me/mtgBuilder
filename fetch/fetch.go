// Package fetch contains a wrapper over the scryfall bulk data api
package fetch

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

type BulkData struct {
	Object  string          `json:"object"`
	HasMore bool            `json:"has_more"`
	Data    []BulkDataEntry `json:"data"`
}

type BulkDataEntry struct {
	Object          string `json:"object"`
	ID              string `json:"id"`
	Type            string `json:"type"`
	UpdatedAt       string `json:"updated_at"`
	URI             string `json:"uri"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Size            int    `json:"size"`
	DownloadURI     string `json:"download_uri"`
	ContentType     string `json:"content_type"`
	ContentEncoding string `json:"content_encoding"`
}

const (
	bulkDataURI                     = "https://api.scryfall.com/bulk-data"
	userAgent                       = "mtgBuilder"
	requestsPerSecond time.Duration = 1
)

var lastRequest time.Time

// globalRatelimit blocks until it is legal to make a request to scryfall
func globalRatelimit() {
	nextRequest := lastRequest.Add(time.Second / requestsPerSecond)
	delay := time.Until(nextRequest)
	if delay > 0 {
		slog.Debug("internal ratelimit", "delayingFor", delay.String())
	}
	time.Sleep(delay)
	lastRequest = time.Now()
}

// scryfallGet is a wrapper around http.Get that complies with scryfall's rest api rules
// see https://scryfall.com/docs/api
func scryfallGet(url string) (*http.Response, error) {
	slog.Debug("Fetching from scryfall", "url", url)
	globalRatelimit()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "*/*")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// fetchBulkData fetches avalibable bulk data entries from scryfall
func fetchBulkData() (BulkData, error) {
	bulkData, err := scryfallGet(bulkDataURI) // required header see https://scryfall.com/docs/api
	if err != nil {
		return BulkData{}, err
	}
	defer bulkData.Body.Close()

	body, err := io.ReadAll(bulkData.Body)
	if err != nil {
		return BulkData{}, err
	}

	var entries BulkData
	err = json.Unmarshal(body, &entries)
	return entries, err
}

// GetOracleCardsJSON fetches the latest available oracle-cards.json from scryfall's bulk data api
func GetOracleCardsJSON() ([]byte, error) {
	bulkData, err := fetchBulkData()
	if err != nil {
		return nil, err
	}

	i := slices.IndexFunc(bulkData.Data, func(e BulkDataEntry) bool { return e.Type == "oracle_cards" })
	oracleCardsURI := bulkData.Data[i].DownloadURI

	oracleCards, err := scryfallGet(oracleCardsURI)
	if err != nil {
		return nil, err
	}
	defer oracleCards.Body.Close()

	bytes, err := io.ReadAll(oracleCards.Body)
	return bytes, err
}
