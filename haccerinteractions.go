package haccerinteractions

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type haccerInteractionsRunner struct {
	Session        *discordgo.Session
	commandRunLock map[string]*sync.Mutex
}

func (hir *haccerInteractionsRunner) getCommandMutex(name string) *sync.Mutex {
	if _, ok := hir.commandRunLock[name]; !ok {
		newMutex := sync.Mutex{}
		hir.commandRunLock[name] = &newMutex
		return &newMutex
	}
	return hir.commandRunLock[name]
}

// Create a command runner
func NewRunner(s *discordgo.Session) haccerInteractionsRunner {
	s.StateEnabled = true
	s.Identify.Intents = discordgo.IntentsAll
	s.Open()
	newRunner := haccerInteractionsRunner{
		Session:        s,
		commandRunLock: make(map[string]*sync.Mutex),
	}
	return newRunner
}

// Get slash commands in a channel. Limit is ignored if application id is set.
func (hir haccerInteractionsRunner) GuildChannelGetSlashCommands(guildID string, applicationID string) (*[]Command, error) {
	response, err := hir.Session.Request(http.MethodGet, fmt.Sprintf(`https://discord.com/api/v9/guilds/%s/application-command-index`, guildID), nil)
	searchResponse := CommandSearchResponse{}
	json.Unmarshal(response, &searchResponse)
	if err != nil {
		return nil, err
	}
	if applicationID!=""{
		found:=[]Command{}
		for _,command:=range searchResponse.Commands{
			if command.ApplicationID==applicationID{
				found = append(found, command)
			}
		}
		return &found,nil
	}

	return &searchResponse.Commands, nil
}

// Interact with a component
// Remember that the top level components in a message are action rows.
func (hir haccerInteractionsRunner) GuildChannelComponentRequest(gc GuildChannel, messageID string, botID string, customID string, data interface{}) error {
	var dataMap map[string]interface{}
	jsonEncoded, err := json.Marshal(data)
	if err != nil {
		return errors.New("error encoding data to json for component type detection: " + err.Error())
	}
	json.Unmarshal(jsonEncoded, &dataMap)
	var componentType discordgo.ComponentType
	switch data.(type) {
	case ButtonClickRequestData:
		componentType = discordgo.ButtonComponent
	case SelectMenuSelectRequestData:
		componentType = discordgo.ComponentType(dataMap["type"].(float64))
	}

	dataMap["component_type"] = componentType
	dataMap["custom_id"] = customID
	source := rand.NewSource(int64(time.Now().Nanosecond()))
	gen := rand.New(source)
	requestStruct := ComponentInteractRequest{
		Type:      discordgo.InteractionMessageComponent,
		Flags:     0,
		SessionID: hir.Session.State.SessionID,
		Nonce:     fmt.Sprint(gen.Intn(99999999999999)),
		GuildID:   gc.GuildID,
		ChannelID: gc.ChannelID,
		MessageID: messageID,
		BotID:     botID,
		Data:      dataMap,
	}
	_, err = hir.Session.Request(http.MethodPost, "https://discord.com/api/v9/interactions", requestStruct)
	return err
}
