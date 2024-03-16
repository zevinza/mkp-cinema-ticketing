package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// client := httpclient.NewHTTPClient(baseURL, authKey)
// HTTPClient struct
type HTTPClient struct {
	baseURL string
	authKey string
	client  *http.Client
}

// NewHTTPClient func
func NewHTTPClient(baseURL string, authKey ...string) *HTTPClient {
	var key string
	if len(authKey) > 0 {
		key = authKey[0]
	}
	return &HTTPClient{
		baseURL: baseURL,
		authKey: key,
		client:  http.DefaultClient,
	}
}

// Get func
func (c *HTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+url, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.client.Do(req)
}

// Post func
func (c *HTTPClient) Post(url string, payload interface{}, contentType ...string) (*http.Response, error) {
	return c.sendRequest(http.MethodPost, url, payload, contentType...)
}

// Put func
func (c *HTTPClient) Put(url string, payload interface{}, contentType ...string) (*http.Response, error) {
	return c.sendRequest(http.MethodPut, url, payload, contentType...)
}

// Delete func
func (c *HTTPClient) Delete(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+url, nil)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req)
	return c.client.Do(req)
}

// setHeaders func
func (c *HTTPClient) setHeaders(req *http.Request, contentType ...string) {
	req.Header.Add(viper.GetString("HEADER_TOKEN_KEY"), viper.GetString("VALUE_TOKEN_KEY"))
	if c.authKey != "" {
		req.Header.Add("Authorization", "Bearer "+c.authKey)
	}
	if len(contentType) > 0 {
		req.Header.Set("Content-Type", contentType[0])
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
}

// sendRequest func
func (c *HTTPClient) sendRequest(method, url string, payload interface{}, contentType ...string) (*http.Response, error) {
	var bodyPayload io.Reader

	switch t := payload.(type) {
	case string:
		bodyPayload = strings.NewReader(t)
	default:
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		bodyPayload = bytes.NewBuffer(jsonPayload)
	}

	req, err := http.NewRequest(method, c.baseURL+url, bodyPayload)
	if err != nil {
		return nil, err
	}

	c.setHeaders(req, contentType...)

	// Log the payload and URL
	log.Printf("Sending request:\nMethod: %s\nURL: %s\nPayload: %s", method, req.URL.String(), payloadToString(payload))

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	// Log the response
	log.Printf("Received response:\nStatus: %s\n", res.Status)

	return res, nil
}

// payloadToString converts the payload to a string for logging purposes
func payloadToString(payload interface{}) string {
	switch t := payload.(type) {
	case string:
		return t
	default:
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return fmt.Sprintf("Error marshaling payload: %v", err)
		}
		return string(jsonPayload)
	}
}
