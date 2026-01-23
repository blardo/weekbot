package handlers

import (
	"fmt"
	"sync"
	"weekbot-go/internal/services/discord"
)

// ChannelCache caches channel IDs to avoid repeated API calls
type ChannelCache struct {
	cache map[string]string // guildID:channelName -> channelID
	mutex sync.RWMutex
	ds    *discord.DiscordService
}

var channelCache *ChannelCache
var cacheOnce sync.Once

// GetChannelCache returns the global channel cache singleton
func GetChannelCache() *ChannelCache {
	cacheOnce.Do(func() {
		ds, err := discord.GetDiscordService()
		if err != nil {
			// Will be initialized later when service is available
			ds = nil
		}
		channelCache = &ChannelCache{
			cache: make(map[string]string),
			ds:    ds,
		}
	})
	return channelCache
}

// SetDiscordService updates the discord service reference
func (cc *ChannelCache) SetDiscordService(ds *discord.DiscordService) {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	cc.ds = ds
}

// GetChannelIDByName gets a channel ID by name, using cache if available
func (cc *ChannelCache) GetChannelIDByName(guildID, channelName string) (string, error) {
	cacheKey := guildID + ":" + channelName
	
	// Check cache first
	cc.mutex.RLock()
	if channelID, exists := cc.cache[cacheKey]; exists {
		cc.mutex.RUnlock()
		return channelID, nil
	}
	cc.mutex.RUnlock()
	
	// Not in cache, fetch from Discord
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	
	// Double-check after acquiring write lock
	if channelID, exists := cc.cache[cacheKey]; exists {
		return channelID, nil
	}
	
	if cc.ds == nil {
		return "", fmt.Errorf("discord service not initialized")
	}
	
	channelID, err := cc.ds.GetChannelIDByName(guildID, channelName)
	if err != nil {
		return "", err
	}
	
	// Cache the result
	cc.cache[cacheKey] = channelID
	return channelID, nil
}

// ClearCache clears the channel cache (useful for testing or when channels change)
func (cc *ChannelCache) ClearCache() {
	cc.mutex.Lock()
	defer cc.mutex.Unlock()
	cc.cache = make(map[string]string)
}
