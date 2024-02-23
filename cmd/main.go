package main

import (
	"fmt"
	weekbot "weekbot/internal/weekbot"
)

func main() {
	config := weekbot.GetConfig()
	bot := weekbot.NewBot(config)

	fmt.Println("Starting Weekbot...")
	bot.Run()
}
