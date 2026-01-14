package main

import (
	"fmt"
	"sync"
	"time"
)

type StoredMessage struct {
	Text      string
	Timestamp time.Time
}

type MessageStore struct {
	messages map[string]map[int]*StoredMessage
	mu       sync.RWMutex
	ttl      time.Duration
}

func NewMessageStore(ttl time.Duration) *MessageStore {
	store := &MessageStore{
		messages: make(map[string]map[int]*StoredMessage),
		ttl:      ttl,
	}
	go store.cleanup()
	return store
}

func (ms *MessageStore) Save(bizConnID string, chatID int64, messageID int, text string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	key := fmt.Sprintf("%s:%d", bizConnID, chatID)
	if ms.messages[key] == nil {
		ms.messages[key] = make(map[int]*StoredMessage)
	}
	ms.messages[key][messageID] = &StoredMessage{
		Text:      text,
		Timestamp: time.Now(),
	}
}

func (ms *MessageStore) Get(bizConnID string, chatID int64, messageID int) (string, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	key := fmt.Sprintf("%s:%d", bizConnID, chatID)
	if chat, exists := ms.messages[key]; exists {
		if msg, found := chat[messageID]; found {
			return msg.Text, true
		}
	}
	return "", false
}

func (ms *MessageStore) Delete(bizConnID string, chatID int64, messageID int) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	key := fmt.Sprintf("%s:%d", bizConnID, chatID)
	if chat, exists := ms.messages[key]; exists {
		delete(chat, messageID)
	}
}

func (ms *MessageStore) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		ms.mu.Lock()
		now := time.Now()

		for chatKey, messages := range ms.messages {
			for msgID, msg := range messages {
				if now.Sub(msg.Timestamp) > ms.ttl {
					delete(messages, msgID)
				}
			}
			if len(messages) == 0 {
				delete(ms.messages, chatKey)
			}
		}

		ms.mu.Unlock()
	}
}

func (ms *MessageStore) Count() int {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	total := 0
	for _, messages := range ms.messages {
		total += len(messages)
	}
	return total
}
