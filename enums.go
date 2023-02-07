package main

type CommandTye int
type OptionType int

const (
	// Command Types
	SLASH_COMMAND CommandTye = 1
	USER_COMMAND  CommandTye = 2
	MESSAGE_COMMAND CommandTye = 3	

	// Option Types
	SUBCOMMAND OptionType = 1
	SUBCOMMAND_GROUP OptionType = 2
	STRING_OPTION OptionType = 3
	BOOLEAN_OPTION OptionType = 5
	USER_OPTION OptionType = 6
	CHANNEL_OPTION OptionType = 7
	ROLE_OPTION OptionType = 8
	MENTIONABLE_OPTION OptionType = 9
	NUMBER_OPTION OptionType = 10
	ATTACHMENT_OPTION OptionType = 11
)

