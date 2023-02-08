package discord

import (
	"encoding/json"
	"fmt"
)

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
}

func (i *Interaction) Respond(payload map[string]interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r := Request{
		Method:  "POST",
		Body:    data,
		Token:   i.Token,
		Headers: map[string]string{"Content-Type": "application/json"},
		Path:    fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token),
	}
	_, err = r.Do()
	return err
}

func (i *Interaction) Defer(ephemeral bool) error {
	var body []byte
	if ephemeral {
		body = []byte(`{"type": 5, "data": {"flags": 64}}`)
	} else {
		body = []byte(`{"type": 5}`)
	}
	r := Request{
		Method:  "POST",
		Token:   i.Token,
		Headers: map[string]string{"Content-Type": "application/json"},
		Path:    fmt.Sprintf("/interactions/%s/%s/callback", i.Id, i.Token),
		Body:    body,
	}
	_, err := r.Do()
	i.Responded = true
	return err
}

func (i *Interaction) Reply(payload map[string]interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	r := Request{
		Body:    data,
		Method:  "POST",
		Token:   i.Token,
		Headers: map[string]string{"Content-Type": "application/json"},
		Path:    fmt.Sprintf("/webhooks/%s/%s", i.ApplicationId, i.Token),
	}
	_, err = r.Do()
	return err
}
