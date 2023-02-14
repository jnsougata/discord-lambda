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
	router        *echo.Echo
	queue         []Command
	commands      map[string]Command
	callbacks     map[string]interface{}
}

func New(applicationId, publickey, token string) *Client {
	return &Client{
		applicationId: applicationId,
		publickey:     publickey,
		token:         token,
		router:        echo.New(),
		queue:         []Command{},
		commands:      map[string]Command{},
		callbacks:     map[string]interface{}{},
		path:          "/interactions",
	}
}

func (c *Client) SetPath(path string) {
	c.path = path
}

func (c *Client) DefaultRouter() *echo.Echo {
	return c.router
}

func (c *Client) AddCommands(commands ...Command) {
	for _, command := range commands {
		if command.Id != "" {
			c.commands[command.Id] = command
		}
	}
	c.queue = append(c.queue, commands...)
}

func (c *Client) RegisterCommands() []map[string]interface{} {
	var payload []map[string]interface{}
	for _, cmd := range c.queue {
		payload = append(payload, cmd.marshal())
	}
	data, _ := json.Marshal(payload)
	r := Request{
		Method:  "PUT",
		Body:    data,
		Token:   c.token,
		Headers: map[string]string{"Content-Type": "application/json"},
		Path:    fmt.Sprintf("/applications/%s/commands", c.applicationId),
	}
	resp, _ := r.Do()
	ba, _ := io.ReadAll(resp.Body)
	var res []map[string]interface{}
	_ = json.Unmarshal(ba, &res)
	return res
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
	e.GET("/sync/:token", c.handleSync)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func (c *Client) handleSync(ctx echo.Context) error {
	if ctx.Param("token") != c.token {
		w := ctx.Response()
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("Unauthorized"))
		return err
	}
	data := c.RegisterCommands()
	w := ctx.Response()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	ba, _ := json.MarshalIndent(data, "", "  ")
	_, err := w.Write(ba)
	return err
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
	inter.client = c
	switch inter.Type {
	case 1:
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		resp, _ := json.Marshal(map[string]interface{}{"type": 1})
		_, err := w.Write(resp)
		return err
	case 2:
		w.WriteHeader(http.StatusOK)
		var cmd Command
		cd, _ := json.Marshal(inter.Data)
		_ = json.Unmarshal(cd, &cmd)
		cc := CommandContext{inter: &inter, Command: &cmd}
		f, exists := c.commands[cmd.Id]
		if exists {
			return f.Callback(&cc)
		}
		return nil
	default:
		w.WriteHeader(http.StatusOK)
		return nil
	}
}
