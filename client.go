package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Client struct {
	DiscordToken string
	DiscordPublicKey string
	DiscordApplicationID string
	InteractionListenerEndpoint string
}

func (c *Client) Run(port int32) error {
	if c.DiscordToken == "" {
		panic("DiscordToken is required")
	}
	if c.DiscordPublicKey == "" {
		panic("DiscordPublicKey is required")
	}
	if c.DiscordApplicationID == "" {
		panic("DiscordApplicationID is required")
	}
	r := mux.NewRouter()
	var ep string
	if c.InteractionListenerEndpoint != "" {
		ep = c.InteractionListenerEndpoint
	} else {
		ep = "/interactions"
	}
	r.HandleFunc(ep, c.HandleInteractions).Methods("POST", "GET")
	http.Handle("/", r)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (c *Client) HandleInteractions(w http.ResponseWriter, r *http.Request) {
	//...
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(c.DiscordApplicationID))
}
