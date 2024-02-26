package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	weekbot "weekbot/internal/weekbot"
)

func main() {
	// Get the bot configuration and create a new bot
	config := weekbot.GetConfig()
	bot, err := weekbot.NewBot(config)
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}
	fmt.Println("Starting Weekbot...")

	// Run the core Discord listener
	// Order here is important so we have a UserID for later
	bot.Run()

	// Setup the bot's commands
	bot.SetupCommands()

	// Setup an interrupt listener
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)

	// Block until we get an interrupt
	<-stop

	// Clean up and stop the bot
	bot.Stop()
	fmt.Println("Shutting down...")
}
