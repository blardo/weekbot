package models

import (
	"testing"
	"weekbot-go/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate the ServerConfig model
	err = db.AutoMigrate(&models.ServerConfig{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func TestGetServerConfig_NotFound(t *testing.T) {
	db := setupTestDB(t)

	config := models.GetServerConfig(db, "test-guild-123")

	if config != nil {
		t.Errorf("Expected nil for non-existent config, got %+v", config)
	}
}

func TestGetOrCreateServerConfig_Creates(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	config := models.GetOrCreateServerConfig(db, guildID)

	if config == nil {
		t.Fatal("Expected config to be created, got nil")
	}

	if config.GuildID != guildID {
		t.Errorf("Expected GuildID %s, got %s", guildID, config.GuildID)
	}

	// All values should be nil (use defaults)
	if config.MinSuggestionsToStart != nil {
		t.Errorf("Expected MinSuggestionsToStart to be nil, got %d", *config.MinSuggestionsToStart)
	}
	if config.MinUpdicksToQualify != nil {
		t.Errorf("Expected MinUpdicksToQualify to be nil, got %d", *config.MinUpdicksToQualify)
	}
	if config.MinBallotsToEndPoll != nil {
		t.Errorf("Expected MinBallotsToEndPoll to be nil, got %d", *config.MinBallotsToEndPoll)
	}
	if config.ReactionEmoji != nil {
		t.Errorf("Expected ReactionEmoji to be nil, got %s", *config.ReactionEmoji)
	}
}

func TestGetOrCreateServerConfig_ReturnsExisting(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	// Create first config
	config1 := models.GetOrCreateServerConfig(db, guildID)
	if config1 == nil {
		t.Fatal("Expected config to be created, got nil")
	}

	// Get again - should return same config
	config2 := models.GetOrCreateServerConfig(db, guildID)
	if config2 == nil {
		t.Fatal("Expected config to be returned, got nil")
	}

	if config1.ID != config2.ID {
		t.Errorf("Expected same config ID, got %d and %d", config1.ID, config2.ID)
	}
}

func TestUpdateServerConfigField_MinSuggestionsToStartPoll(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	err := models.UpdateServerConfigField(db, guildID, "min-suggestions-to-start-poll", 5)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	config := models.GetServerConfig(db, guildID)
	if config == nil {
		t.Fatal("Expected config to exist")
	}

	if config.MinSuggestionsToStart == nil {
		t.Fatal("Expected MinSuggestionsToStart to be set")
	}

	if *config.MinSuggestionsToStart != 5 {
		t.Errorf("Expected MinSuggestionsToStart to be 5, got %d", *config.MinSuggestionsToStart)
	}
}

func TestUpdateServerConfigField_MinUpdicks(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	err := models.UpdateServerConfigField(db, guildID, "min-updicks", 10)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	config := models.GetServerConfig(db, guildID)
	if config == nil {
		t.Fatal("Expected config to exist")
	}

	if config.MinUpdicksToQualify == nil {
		t.Fatal("Expected MinUpdicksToQualify to be set")
	}

	if *config.MinUpdicksToQualify != 10 {
		t.Errorf("Expected MinUpdicksToQualify to be 10, got %d", *config.MinUpdicksToQualify)
	}
}

func TestUpdateServerConfigField_MinBallotsToEndPoll(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	err := models.UpdateServerConfigField(db, guildID, "min-ballots-to-end-poll", 15)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	config := models.GetServerConfig(db, guildID)
	if config == nil {
		t.Fatal("Expected config to exist")
	}

	if config.MinBallotsToEndPoll == nil {
		t.Fatal("Expected MinBallotsToEndPoll to be set")
	}

	if *config.MinBallotsToEndPoll != 15 {
		t.Errorf("Expected MinBallotsToEndPoll to be 15, got %d", *config.MinBallotsToEndPoll)
	}
}

func TestUpdateServerConfigField_InvalidField(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	err := models.UpdateServerConfigField(db, guildID, "invalid-field", 5)
	if err == nil {
		t.Error("Expected error for invalid field, got nil")
	}
}

func TestUpdateServerConfigStringField_ReactionEmoji(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	err := models.UpdateServerConfigStringField(db, guildID, "reaction-emoji", "🔥")
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	config := models.GetServerConfig(db, guildID)
	if config == nil {
		t.Fatal("Expected config to exist")
	}

	if config.ReactionEmoji == nil {
		t.Fatal("Expected ReactionEmoji to be set")
	}

	if *config.ReactionEmoji != "🔥" {
		t.Errorf("Expected ReactionEmoji to be 🔥, got %s", *config.ReactionEmoji)
	}
}

func TestUpdateServerConfigStringField_InvalidField(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	err := models.UpdateServerConfigStringField(db, guildID, "invalid-field", "test")
	if err == nil {
		t.Error("Expected error for invalid field, got nil")
	}
}

func TestMultipleGuilds_IndependentConfigs(t *testing.T) {
	db := setupTestDB(t)
	guild1 := "guild-1"
	guild2 := "guild-2"

	// Set different values for each guild
	err := models.UpdateServerConfigField(db, guild1, "min-updicks", 5)
	if err != nil {
		t.Fatalf("Failed to update guild1 config: %v", err)
	}

	err = models.UpdateServerConfigField(db, guild2, "min-updicks", 20)
	if err != nil {
		t.Fatalf("Failed to update guild2 config: %v", err)
	}

	// Verify each guild has its own config
	config1 := models.GetServerConfig(db, guild1)
	config2 := models.GetServerConfig(db, guild2)

	if config1 == nil || config2 == nil {
		t.Fatal("Expected both configs to exist")
	}

	if *config1.MinUpdicksToQualify != 5 {
		t.Errorf("Expected guild1 MinUpdicksToQualify to be 5, got %d", *config1.MinUpdicksToQualify)
	}

	if *config2.MinUpdicksToQualify != 20 {
		t.Errorf("Expected guild2 MinUpdicksToQualify to be 20, got %d", *config2.MinUpdicksToQualify)
	}
}

func TestUpdateServerConfigField_UpdatesExisting(t *testing.T) {
	db := setupTestDB(t)
	guildID := "test-guild-123"

	// Set initial value
	err := models.UpdateServerConfigField(db, guildID, "min-updicks", 5)
	if err != nil {
		t.Fatalf("Failed to set initial config: %v", err)
	}

	// Update to new value
	err = models.UpdateServerConfigField(db, guildID, "min-updicks", 10)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	config := models.GetServerConfig(db, guildID)
	if config == nil {
		t.Fatal("Expected config to exist")
	}

	if *config.MinUpdicksToQualify != 10 {
		t.Errorf("Expected MinUpdicksToQualify to be 10, got %d", *config.MinUpdicksToQualify)
	}

	// Verify only one config exists for this guild
	var count int64
	db.Model(&models.ServerConfig{}).Where("guild_id = ?", guildID).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 config for guild, got %d", count)
	}
}
