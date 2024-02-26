package main

import (
	"fmt"
	"os"
	"os/signal"
	weekbot "weekbot/internal/weekbot"
)

func main() {
	config := weekbot.GetConfig()
	bot, err := weekbot.NewBot(*config)
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}
	fmt.Println("Starting Weekbot...")

	bot.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	// Wait for a signal to stop the bot
	<-stop

	bot.Stop()
	fmt.Println("Shutting down...")
}
