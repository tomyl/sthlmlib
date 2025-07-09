package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
)

const (
	apiEndpoint = "https://biblioteket.stockholm.se/graphql/"
)

// Client is a client for the Stockholm Public Library API.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new authenticated client.
func NewClient(cardNumber, pinCode string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("could not create cookie jar: %w", err)
	}

	client := &Client{
		httpClient: &http.Client{
			Jar: jar,
		},
	}

	if err := client.login(cardNumber, pinCode); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	return client, nil
}

func (c *Client) login(cardNumber, pinCode string) error {
	requestBody := LoginRequest{
		Query: `
  mutation loginPatron($loginInput: InputLogin!) {
    loginPatron(input: $loginInput) {
      status
      statusMessage
    }
  }
`,
		Variables: Variables{
			Operation: "loginPatron",
			LoginInput: LoginInput{
				CardNumber: cardNumber,
				PinCode:    pinCode,
			},
		},
	}

	var loginResponse LoginResponse
	if err := c.query(&requestBody, &loginResponse); err != nil {
		return err
	}

	if loginResponse.Data.LoginPatron.Status != "OK" {
		return fmt.Errorf("login failed with status: %s", loginResponse.Data.LoginPatron.Status)
	}

	return nil
}

func (c *Client) query(requestBody, responseBody interface{}) error {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("could not marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "sthlmlib-go-client")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("could not send request: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
		return fmt.Errorf("could not decode response: %w", err)
	}

	return nil
}
