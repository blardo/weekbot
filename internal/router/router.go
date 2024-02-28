package router

import (
	"weekbot/internal/commands"
	"weekbot/internal/handlers"
	"weekbot/internal/models"
)



func ConfigureBot(b *models.Bot) *models.Bot{
	b.DSC.AddHandler(handlers.ParseInteraction)
	b.DSC.AddHandler(handlers.ParseChatCommand)
	b.DSC.AddHandler(handlers.HandleReactions)
	return b
}

func SetCommands(b *models.Bot) *models.Bot{
	b.DSC.AddSlashCommand(commands.GetPingCommand())
	b.DSC.AddSlashCommand(commands.GetPollCommand())
	return b
}


