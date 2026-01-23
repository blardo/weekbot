package handlers

import (
	"log"
	"sync"
	"time"
)

// MessageTracker prevents duplicate processing of Discord events
type MessageTracker struct {
	processedMessages map[string]time.Time
	processedReactions map[string]time.Time
	mutex             sync.RWMutex
	cleanupInterval   time.Duration
	cleanupTicker     *time.Ticker
}

var globalTracker *MessageTracker
var trackerOnce sync.Once

// GetMessageTracker returns the global message tracker singleton
func GetMessageTracker() *MessageTracker {
	trackerOnce.Do(func() {
		globalTracker = &MessageTracker{
			processedMessages:   make(map[string]time.Time),
			processedReactions:  make(map[string]time.Time),
			cleanupInterval:     5 * time.Minute,
		}
		globalTracker.startCleanup()
	})
	return globalTracker
}

// startCleanup starts a periodic cleanup goroutine
func (mt *MessageTracker) startCleanup() {
	mt.cleanupTicker = time.NewTicker(mt.cleanupInterval)
	go func() {
		for range mt.cleanupTicker.C {
			mt.cleanup()
		}
	}()
}

// cleanup removes old entries from the tracking maps
func (mt *MessageTracker) cleanup() {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()
	
	cutoff := time.Now().Add(-mt.cleanupInterval)
	
	for msgID, timestamp := range mt.processedMessages {
		if timestamp.Before(cutoff) {
			delete(mt.processedMessages, msgID)
		}
	}
	
	for reactionID, timestamp := range mt.processedReactions {
		if timestamp.Before(cutoff) {
			delete(mt.processedReactions, reactionID)
		}
	}
}

// ShouldProcessMessage returns true if the message should be processed (not a duplicate)
func (mt *MessageTracker) ShouldProcessMessage(messageID string) bool {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()
	
	if _, exists := mt.processedMessages[messageID]; exists {
		log.Printf("Message %s already processed, skipping", messageID)
		return false
	}
	
	mt.processedMessages[messageID] = time.Now()
	return true
}

// ShouldProcessReaction returns true if the reaction should be processed (not a duplicate)
func (mt *MessageTracker) ShouldProcessReaction(channelID, messageID, userID, emoji string) bool {
	mt.mutex.Lock()
	defer mt.mutex.Unlock()
	
	reactionID := channelID + ":" + messageID + ":" + userID + ":" + emoji
	if _, exists := mt.processedReactions[reactionID]; exists {
		log.Printf("Reaction %s already processed, skipping", reactionID)
		return false
	}
	
	mt.processedReactions[reactionID] = time.Now()
	return true
}

// Stop stops the cleanup ticker (for testing/cleanup)
func (mt *MessageTracker) Stop() {
	if mt.cleanupTicker != nil {
		mt.cleanupTicker.Stop()
	}
}
