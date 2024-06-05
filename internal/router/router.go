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

func SetCommands(dsc *discord.DiscordService, guildId string) *discord.DiscordService {


		println(guildId)
		dsc.AddSlashCommand(commands.GetPingCommand(), guildId)
		println("Adding ping command")
		dsc.AddSlashCommand(commands.GetPollCommand(), guildId)
		println("Adding poll command")
		dsc.AddSlashCommand(commands.EndPollCommand(), guildId)
		println("Adding endpoll command")
		
	return dsc
}
