package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Payload struct {
	Ttl uint32 `json:"ttl"`
}

var config *Config

func main() {
	config = loadConfig()

	http.HandleFunc("/", getRoot)

	address := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	log.Printf("listening on %s", address)

	err := http.ListenAndServe(address, nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server closed")
	} else if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}

// as per https://datatracker.ietf.org/doc/html/draft-uberti-behave-turn-rest-00#section-2.2
type UbertiTurnResponse struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Ttl      uint32   `json:"ttl"`
	Urls     []string `json:"uris"`
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("Incoming: %s %s", r.Method, r.RemoteAddr)

	params := parseParams(r.RequestURI)
	payload := Payload{
		Ttl: config.Cloudflare.Ttl,
	}

	credentials, err := requestCredentials(config.Cloudflare.KeyId, params, payload)
	if err != nil {
		log.Printf("Failed to request credentials: %s", err)
	}

	urls := []string{}
	for _, url := range credentials.IceServers.Urls {
		if strings.HasPrefix(url, "turn:") || strings.HasPrefix(url, "turns:") {
			urls = append(urls, url)
		}
	}

	responsePayload := &UbertiTurnResponse{
		Username: credentials.IceServers.Username,
		Password: credentials.IceServers.Password,
		Ttl:      config.Cloudflare.Ttl,
		Urls:     urls,
	}

	jsonPayload, err := json.Marshal(responsePayload)
	if err != nil {
		log.Printf("Failed to marshal payload: %s\n", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonPayload)
}
