package haccerinteractions

import (
	"errors"
	"net/http"
	"strings"
	"time"

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

	cmdMutex := hir.getCommandMutex(c.Name)
	cmdMutex.Lock()
	commandRespChan := make(chan discordgo.Message)
	deleteHand := hir.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Interaction == nil {
			return
		}
		words := strings.Split(m.Interaction.Name, " ")
		if len(words) < 1 {
			return
		}
		if m.Interaction.User.ID == s.State.User.ID && words[0] == c.Name && m.Author.ID == c.ApplicationID {
			commandRespChan <- *m.Message
			cmdMutex.Unlock()
		}
	})
	_, err := hir.Session.Request(http.MethodPost, "https://discord.com/api/v9/interactions", reqData)
	if err != nil {
		return nil, err
	}
	timer := time.NewTimer(time.Second * 30)
	select {
	case cmdRespMsg := <-commandRespChan:
		{
			deleteHand()
			return &cmdRespMsg, nil
		}
	case <-timer.C:
		{
			deleteHand()
			return nil, errors.New("timeout looking for response")
		}
	}
}
