package discord

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Client struct {
	token         string
	publickey     string
	applicationId string
	path          string
	commands      []Command
	router        *echo.Echo
}

func New(applicationId, publickey, token string) *Client {
	return &Client{
		applicationId: applicationId,
		publickey:     publickey,
		token:         token,
		router:        echo.New(),
		commands:      []Command{},
		path:          "/interactions",
	}
}

func (c *Client) SetPath(path string) {
	c.path = path
}

func (c *Client) AddCommands(commands ...Command) {
	c.commands = append(c.commands, commands...)
}

func (c *Client) RegisterCommands() int {
	var commands []map[string]interface{}
	for _, command := range c.commands {
		commands = append(commands, command.marshal())
	}
	data, _ := json.Marshal(commands)
	r := Request{
		Method:  "PUT",
		Body:    data,
		Token:   c.token,
		Headers: map[string]string{"Content-Type": "application/json"},
		Path:    fmt.Sprintf("/applications/%s/commands", c.applicationId),
	}
	resp, _ := r.Do()
	return resp.StatusCode
}

func (c *Client) Run(port int) {
	if c.token == "" {
		panic("Discord Bot Token is required")
	}
	if c.publickey == "" {
		panic("Discord Bot PublicKey is required")
	}
	if c.applicationId == "" {
		panic("Discord Bot ApplicationID is required")
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

func (c *Client) handleInteractionGET(ctx echo.Context) error {
	w := ctx.Response()
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	return err
}

func (c *Client) handleInteractionPOST(ctx echo.Context) error {
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
	var inter Interaction
	_ = json.Unmarshal(ba, &inter)
	switch inter.Type {
	// ping
	case 1:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp, _ := json.Marshal(map[string]interface{}{"type": 1})
		_, err := w.Write(resp)
		return err
	// application command
	case 2:
		w.WriteHeader(http.StatusOK)
		data, _ := json.Marshal(map[string]interface{}{
			"type": 4,
			"data": map[string]interface{}{
				"content": "pong!",
			},
		})
		r := Request{
			Method:  "POST",
			Body:    data,
			Path:    fmt.Sprintf("/interactions/%s/%s/callback", inter.Id, inter.Token),
			Headers: map[string]string{"Content-Type": "application/json"},
		}
		_, err := r.Do()
		return err
	default:
		w.WriteHeader(http.StatusOK)
		return nil
	}
}
