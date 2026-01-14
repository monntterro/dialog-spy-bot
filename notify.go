package main

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func sendNotification(ctx context.Context, b *bot.Bot, userID int64, text string) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}
