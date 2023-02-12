package discord

import (
	"encoding/json"
	"fmt"
)

type AppCommand struct {
	Id          string   `json:"id"`
	Application string   `json:"application_id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Options     []Option `json:"options"`
	GuildId     string   `json:"guild_id"`
	TargetId    string   `json:"target_id"`
}

type CommandContext struct {
	inter     *Interaction
	Command   *AppCommand
	Responded bool
}

func (cc *CommandContext) Respond(payload map[string]interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r := Request{
		Method:  "POST",
		Body:    data,
		Headers: map[string]string{"Content-Type": "application/json"},
		Path:    fmt.Sprintf("/interactions/%s/%s/callback", cc.inter.Id, cc.inter.Token),
	}
	_, err = r.Do()
	return err
}

func (cc *CommandContext) Defer(ephemeral bool) error {
	var body []byte
	if ephemeral {
		body = []byte(`{"type": 5, "data": {"flags": 64}}`)
	} else {
		body = []byte(`{"type": 5}`)
	}
	r := Request{
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/json"},
		Path:    fmt.Sprintf("/interactions/%s/%s/callback", cc.inter.Id, cc.inter.Token),
		Body:    body,
	}
	_, err := r.Do()
	cc.Responded = true
	return err
}

func (cc *CommandContext) Reply(payload map[string]interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r := Request{
		Body:    data,
		Method:  "POST",
		Headers: map[string]string{"Content-Type": "application/json"},
		Path:    fmt.Sprintf("/webhooks/%s/%s", cc.inter.ApplicationId, cc.inter.Token),
	}
	_, err = r.Do()
	return err
}

type Interaction struct {
	Id             string                 `json:"id"`
	Type           int                    `json:"type"`
	Token          string                 `json:"token"`
	Version        int                    `json:"version"`
	ApplicationId  string                 `json:"application_id"`
	GuildId        string                 `json:"guild_id"`
	ChannelId      string                 `json:"channel_id"`
	AppPermissions string                 `json:"app_permissions"`
	Locale         string                 `json:"locale"`
	GuildLocale    string                 `json:"guild_locale"`
	Data           map[string]interface{} `json:"data"`
	User           map[string]interface{} `json:"user"`
	Member         map[string]interface{} `json:"member"`
	Responded      bool
	client         *Client
}
