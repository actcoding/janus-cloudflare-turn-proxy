package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type CloudflareResponse struct {
	IceServers CloudflareIceServers `json:"iceServers"`
}

type CloudflareIceServers struct {
	Urls     []string `json:"urls"`
	Username string   `json:"username"`
	Password string   `json:"credential"`
}

func requestCredentials(keyId string, params Params, payload Payload) (*CloudflareResponse, error) {
	upstream := fmt.Sprintf("https://rtc.live.cloudflare.com/v1/turn/keys/%s/credentials/generate", keyId)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", upstream, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+params["key"])
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		return nil, errors.New("invalid response: " + string(body))
	}

	var credentials CloudflareResponse
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		return nil, err
	}

	return &credentials, nil
}
