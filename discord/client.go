package discord

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type client struct {
	token         string
	publickey     string
	applicationId string
	path          string
	commands      []Command
	router        *echo.Echo
}

func New(applicationId, publickey, token string) *client {
	return &client{
		applicationId: applicationId,
		publickey:     publickey,
		token:         token,
		router:        echo.New(),
		commands:      []Command{},
		path:          "/interactions",
	}
}

func (c *client) SetPath(path string) {
	c.path = path
}

func (c *client) AddCommands(commands ...Command) {
	c.commands = append(c.commands, commands...)
}

func (c *client) RegisterCommands() error {
	// temporary .. will add bulk register later
	for _, command := range c.commands {
		b, err := json.Marshal(command.marshal())
		if err != nil {
			return err
		}
		req, _ := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/applications/%s/commands", c.applicationId), bytes.NewReader(b))
		req.Header.Set("Authorization", fmt.Sprintf("Bot %s", c.token))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		fmt.Println(resp.StatusCode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) Run(port int) {
	if c.token == "" {
		panic("DiscordToken is required")
	}
	if c.publickey == "" {
		panic("DiscordPublicKey is required")
	}
	if c.applicationId == "" {
		panic("DiscordApplicationID is required")
	}
	var ep string
	if c.path != "" {
		ep = c.path
	} else {
		ep = "/interactions"
	}
	e := echo.New()
	e.POST(ep, c.handleInteractionPOST)
	e.GET(ep, c.handleInteractionGET)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func (c *client) handleInteractionGET(ctx echo.Context) error {
	w := ctx.Response()
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	return err
}

func (c *client) handleInteractionPOST(ctx echo.Context) error {
	r := ctx.Request()
	w := ctx.Response()
	ba, _ := io.ReadAll(r.Body)
	body := string(ba)
	timestamp := r.Header.Get("X-Signature-Timestamp")
	sig, _ := hex.DecodeString(r.Header.Get("X-Signature-Ed25519"))
	pk, _ := hex.DecodeString(c.publickey)
	ok := ed25519.Verify(pk, []byte(timestamp+body), sig)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("invalid request signature"))
		return err
	}
	var data map[string]interface{}
	_ = json.Unmarshal(ba, &data)
	interactionType := data["type"].(float64)
	switch interactionType {
	case 1:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp, _ := json.Marshal(map[string]interface{}{"type": 1})
		_, err := w.Write(resp)
		return err
	default:
		// add interaction struct later
		w.WriteHeader(http.StatusOK)
		d, _ := json.Marshal(map[string]interface{}{
			"type": 4,
			"data": map[string]interface{}{
				"content": "pong!",
			},
		})
		r := Request{
			Method:  "POST",
			Path:    fmt.Sprintf("/interactions/%s/%s/callback", data["id"].(string), data["token"].(string)),
			Body:    d,
			Headers: map[string]string{"Content-Type": "application/json"},
		}
		_, err := r.Do()
		return err
	}
}
