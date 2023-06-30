package haccerinteractions

import "github.com/bwmarrin/discordgo"

type Command struct {
	Type          int             `json:"type"`
	ID            string          `json:"id"`
	ApplicationID string          `json:"application_id"`
	Name          string          `json:"name"`
	Version       string          `json:"version"`
	Options       []CommandOption `json:"options"`
}
type GuildChannel struct {
	GuildID   string
	ChannelID string
}

type commandRunRequest struct {
	Type      int            `json:"type"`
	BotID     string         `json:"application_id"`
	GuildID   string         `json:"guild_id"`
	ChannelID string         `json:"channel_id"`
	SessionID string         `json:"session_id"`
	Data      commandRunData `json:"data"`
}
type commandRunData struct {
	ID      string             `json:"id"`
	Type    int                `json:"type"`
	Name    string             `json:"name"`
	Version string             `json:"version"`
	Options []CommandRunOption `json:"options"`
}
type CommandRunOption struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CommandSearchResponse struct {
	Commands []Command `json:"application_commands"`
}

type CommandOption struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}
type ComponentInteractRequest struct {
	Type      discordgo.InteractionType `json:"type"`
	Flags     int                       `json:"message_flags"`
	SessionID string                    `json:"session_id"`
	Nonce     string                    `json:"nonce"`
	GuildID   string                    `json:"guild_id"`
	ChannelID string                    `json:"channel_id"`
	MessageID string                    `json:"message_id"`
	BotID     string                    `json:"application_id"`
	Data      interface{}               `json:"data"`
}

type ButtonClickRequestData struct{}
type SelectMenuSelectRequestData struct {
	Type   discordgo.SelectMenuType `json:"type"`
	Values []string                 `json:"values"`
}
