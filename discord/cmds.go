package discord

type Choice struct {
	Name  string
	Value interface{}
}

func (c *Choice) marshal() map[string]interface{} {
	return map[string]interface{}{
		"name":  c.Name,
		"value": c.Value,
	}
}

type Option struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        int      `json:"type"`
	Required    bool     `json:"required"`
	Choices     []Choice `json:"choices"`
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
	Id           string
	Type         CommandTye
	Name         string
	Description  string
	Options      []Option
	GuildId      string
	DefaultPerms int
	AllowInDMs   bool
	NSFW         bool
	subcommands  []SubCommand
	Callback     func(interaction *Interaction) map[string]interface{}
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
