package main

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
	Id          string
	Type        CommandTye
	Name        string
	Description string
	Options     []Option
	subcommands []SubCommand
}

func (c *Command) AddSubCommand(subcommand SubCommand) {
	c.subcommands = append(c.subcommands, subcommand)
}

func (c *Command) marshal() map[string]interface{} {
	d := map[string]interface{}{
		"name":        c.Name,
		"description": c.Description,
		"options":     []map[string]interface{}{},
	}
	for _, option := range c.Options {
		d["options"] = append(d["options"].([]map[string]interface{}), option.marshal())
	}
	for _, subcommand := range c.subcommands {
		d["options"] = append(d["options"].([]map[string]interface{}), subcommand.marshal())
	}
	return d
}
