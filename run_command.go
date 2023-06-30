package haccerinteractions

import (
	"net/http"

	"github.com/bwmarrin/discordgo"
)

// Run a command in a channel
func (hir *haccerInteractionsRunner) GuildChannelRunCommand(c Command, args *[]CommandRunOption, gc GuildChannel) (*discordgo.Message, error) {
	if args == nil {
		args = &[]CommandRunOption{}
	}
	reqData := commandRunRequest{
		Type:      2,
		BotID:     c.ApplicationID,
		GuildID:   gc.GuildID,
		ChannelID: gc.ChannelID,
		SessionID: hir.Session.State.SessionID,
		Data: commandRunData{
			ID:      c.ID,
			Type:    c.Type,
			Name:    c.Name,
			Version: c.Version,
			Options: *args,
		},
	}
	_, err := hir.Session.Request(http.MethodPost, "https://discord.com/api/v9/interactions", reqData)
	if err != nil {
		return nil, err
	}
	cmdMutex := hir.getCommandMutex(c.Name)
	cmdMutex.Lock()
	commandRespChan := make(chan discordgo.Message)
	hir.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Interaction == nil {
			return
		}
		if m.Interaction.User.ID == s.State.User.ID && m.Interaction.Name == c.Name {
			commandRespChan <- *m.Message
			cmdMutex.Unlock()
		}
	})
	cmdRespMsg := <-commandRespChan
	return &cmdRespMsg, nil
}
