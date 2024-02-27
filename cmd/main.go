package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"weekbot/internal/models"
	"weekbot/internal/router"
	"weekbot/internal/services"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Connect to the database
	db, err := gorm.Open(sqlite.Open("weekbot.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Poll{})

	// Get the bot configuration and create a new bot
	config := services.GetConfig(db)
	router := router.NewRouter()
	bot, err := models.NewBot(config, router)
	if err != nil {
		fmt.Println("Error creating bot:", err)
		return
	}
	fmt.Println("Starting Weekbot...")

	// Start the core Discord listener
	// Order here is important so we have a UserID for later
	bot.Start()

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
