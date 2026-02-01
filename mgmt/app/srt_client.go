package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	onceSrtClientApi sync.Once
	srtClientApi     *SrtClientApi
)

type SrtClientApi struct {
	BaseURL string
	Client  *http.Client
}

func NewSrtClientApi(baseURL string) *SrtClientApi {
	onceSrtClientApi.Do(func() {
		srtClientApi = &SrtClientApi{
			BaseURL: baseURL,
			Client:  &http.Client{Timeout: 10 * time.Second},
		}
	})
	return srtClientApi
}

func RePlaceSrtClientApi(baseURL string) *SrtClientApi {
	srtClientApi = &SrtClientApi{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
	return srtClientApi
}

func GetSrtClientApi() *SrtClientApi {
	return srtClientApi
}

func (c *SrtClientApi) get(path string, out any) error {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SRS API error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

// =======================
// STREAM APIs
// =======================

func (c *SrtClientApi) GetStreams() (*StreamsResponse, error) {
	var resp StreamsResponse
	err := c.get("/api/v1/streams", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SrtClientApi) GetStream(id string) (*StreamResponse, error) {
	var resp StreamResponse
	err := c.get("/api/v1/streams/"+id, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// =======================
// CLIENT APIs
// =======================

func (c *SrtClientApi) GetClients() (*ClientsResponse, error) {
	var resp ClientsResponse
	err := c.get("/api/v1/clients", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SrtClientApi) GetClient(id string) (*ClientResponse, error) {
	var resp ClientResponse
	err := c.get("/api/v1/clients/"+id, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// =======================
// SUMMARY API
// =======================

func (c *SrtClientApi) GetSummary() (*SummaryResponse, error) {
	var resp SummaryResponse
	err := c.get("/api/v1/summaries", &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
