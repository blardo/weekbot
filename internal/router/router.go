package router

import (
	"weekbot/internal/commands"
	"weekbot/internal/handlers"
	"weekbot/internal/services/discord"
)

func ConfigureHandlers(dsc *discord.DiscordService) *discord.DiscordService {
	dsc.AddHandler(handlers.ParseInteraction)
	dsc.AddHandler(handlers.ParseChatCommand)
	dsc.AddHandler(handlers.HandleReactions)
	return dsc
}

func SetCommands(dsc *discord.DiscordService) *discord.DiscordService {
	dsc.AddSlashCommand(commands.GetPingCommand())
	dsc.AddSlashCommand(commands.GetPollCommand())
	return dsc
}
