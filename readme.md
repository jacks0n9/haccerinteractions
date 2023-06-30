# Haccerinteractions
## Welcome to Haccerinteractions, a package for advanced selfbotters
Haccerinteractions is a package that expands the selfbotting capabilities of discordgo by giving you the ability to run slash commands, get their output, and interact with message components like buttons.
## Getting started
After installing and importing Haccerinteractions, follow these steps.
### Create a new discordgo session
Because Haccerinteractions is powered by discordgo, you first need to make a session with this code:
```go
session,_:=discordgo.New("insert_token_here")
```
### Create a Haccerinteractions runner
This will be the thing that runs commands and clicks buttons. Use this code to make a runner:
```go
runner:=haccerinteractions.NewRunner(session)
```
This will configure required parameters within your session as well as opening a gateway connection.
### Get commands
You first must get commands for a channel with this code:
```go
// You can leave the application/bot id empty.
// The limit is ignored when the bot_id is set, however
commands,err:=runner.GuildChannelGetSlashCommands("channel_id_here",10,"bot_id_here")
```
### Run Commands
```go
// Run a command.
// Set arguments to nil to pass no arguments
responseMessage,err:=runner.GuildChannelRunCommand(commands[0],nil,"channel_id_here")
```
### Interact with components
```go
// You can interact with various components by using different request data structs
button:=responseMessage.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.Button)
runner.GuildChannelComponentRequest("channel_id_here",responseMessage.ID,"bot_id_here",button.CustomID,haccerinteractions.ButtonClickRequestData{})
```
