package router

import (
	"weekbot/internal/handlers"
	"weekbot/internal/services/discord"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

// Setup adds the bot's handlers to the Discord client
func (r *Router) Setup(ds *discord.DiscordService) {
	ds.AddHandler(handlers.ParseInteraction)
	ds.AddHandler(handlers.ParseChatCommand)
	ds.AddHandler(handlers.HandleReactions)
}
