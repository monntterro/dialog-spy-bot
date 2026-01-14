package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN is not set")
	}

	yourUserIDStr := os.Getenv("YOUR_USER_ID")
	if yourUserIDStr == "" {
		log.Fatal("YOUR_USER_ID is not set")
	}

	yourUserID, err := strconv.ParseInt(yourUserIDStr, 10, 64)
	if err != nil {
		log.Fatal("YOUR_USER_ID must be int64:", err)
	}

	ttlHours := 24
	if ttlStr := os.Getenv("MESSAGE_TTL_HOURS"); ttlStr != "" {
		if parsed, err := strconv.Atoi(ttlStr); err == nil && parsed > 0 {
			ttlHours = parsed
		}
	}
	messageTTL := time.Duration(ttlHours) * time.Hour

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	store := NewMessageStore(messageTTL)

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			handleUpdate(ctx, b, update, store, yourUserID)
		}),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		return
	}

	b.Start(ctx)
}
