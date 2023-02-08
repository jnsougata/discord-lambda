package discord

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
}
