package discord

type Choice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (c *Choice) marshal() map[string]interface{} {
	return map[string]interface{}{
		"name":  c.Name,
		"value": c.Value,
	}
}

type Option struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        OptionType `json:"type"`
	Required    bool       `json:"required"`
	Choices     []Choice   `json:"choices"`
	Value       string     `json:"value"`
	Options     []Option   `json:"options"`
	Focused     bool       `json:"focused"`
}

func (o *Option) marshal() map[string]interface{} {
	d := map[string]interface{}{
		"name":        o.Name,
		"description": o.Description,
		"type":        o.Type,
		"required":    o.Required,
	}
	for _, choice := range o.Choices {
		d["choices"] = append(d["choices"].([]map[string]interface{}), choice.marshal())
	}
	return d
}

type SubCommand struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Options     []Option `json:"options"`
}

func (s *SubCommand) marshal() map[string]interface{} {
	return map[string]interface{}{
		"type":        1,
		"name":        s.Name,
		"description": s.Description,
		"options":     s.Options,
	}
}

type Command struct {
	Id           string                          `json:"id,omitempty"`
	Type         CommandTye                      `json:"type"`
	Name         string                          `json:"name"`
	Description  string                          `json:"description"`
	Options      []Option                        `json:"options"`
	GuildId      string                          `json:"guild_id,omitempty"`
	DefaultPerms int                             `json:"default_member_permissions,omitempty"`
	AllowInDMs   bool                            `json:"dm_permissions,omitempty"`
	NSFW         bool                            `json:"nsfw,omitempty"`
	TargetId     string                          `json:"target_id,omitempty"`
	Resolved     map[string]interface{}          `json:"resolved,omitempty"`
	subcommands  []SubCommand                    `json:"-"`
	Callback     func(ctx *CommandContext) error `json:"-"`
}

func (c *Command) AddSubCommand(subcommand SubCommand) {
	c.subcommands = append(c.subcommands, subcommand)
}

func (c *Command) marshal() map[string]interface{} {
	payload := map[string]interface{}{
		"name":                       c.Name,
		"type":                       c.Type,
		"nsfw":                       c.NSFW,
		"description":                c.Description,
		"dm_permissions":             c.AllowInDMs,
		"default_member_permissions": c.DefaultPerms,
		"options":                    []map[string]interface{}{},
	}
	if c.Type == 0 {
		panic("command type not set!")
	}
	if c.Type != 1 && c.Description != "" {
		panic("command description can only be set for slash commands")
	}
	var ops []map[string]interface{}
	for _, option := range c.Options {
		ops = append(ops, option.marshal())
	}
	payload["options"] = ops
	var scs []map[string]interface{}
	for _, subcommand := range c.subcommands {
		scs = append(scs, subcommand.marshal())
	}
	payload["options"] = append(payload["options"].([]map[string]interface{}), scs...)
	if c.Id != "" {
		payload["id"] = c.Id
	}
	if c.GuildId != "" {
		payload["guild_id"] = c.GuildId
	}
	return payload
}
