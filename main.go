package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("Переменная окружения BOT_TOKEN не установлена")
	}

	yourUserIDStr := os.Getenv("YOUR_USER_ID")
	if yourUserIDStr == "" {
		log.Fatal("Переменная окружения YOUR_USER_ID не установлена")
	}

	yourUserID, err := strconv.ParseInt(yourUserIDStr, 10, 64)
	if err != nil {
		log.Fatal("Неверный формат YOUR_USER_ID:", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			handleUpdate(ctx, b, update, yourUserID)
		}),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		log.Fatal("Ошибка создания бота:", err)
	}

	b.Start(ctx)
}

func handleUpdate(ctx context.Context, b *bot.Bot, update *models.Update, yourUserID int64) {
	if update.Message != nil && update.Message.Text != "" {
		if update.Message.Text == "/test" && update.Message.From != nil {
			if update.Message.From.ID != yourUserID {
				return
			}

			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: yourUserID,
				Text:   "✅ Bot is working",
			})
			return
		}
	}

	if update.EditedBusinessMessage != nil {
		edited := update.EditedBusinessMessage

		if edited.From != nil && edited.From.ID == yourUserID {
			return
		}

		chatTitle := getChatTitle(edited.Chat)
		userName := getUserName(edited.From)

		notification := fmt.Sprintf(
			"✏️ <b>%s</b> | %s\n"+
				"━━━━━━━━━━━━━━━\n"+
				"%s",
			userName,
			chatTitle,
			escapeHTML(edited.Text),
		)

		sendNotification(ctx, b, yourUserID, notification)
		return
	}

	if update.DeletedBusinessMessages != nil {
		deleted := update.DeletedBusinessMessages
		chatTitle := getChatTitle(deleted.Chat)

		notification := fmt.Sprintf(
			"🗑 <b>%s</b>\n"+
				"━━━━━━━━━━━━━━━\n"+
				"Удалено сообщений: %d",
			chatTitle,
			len(deleted.MessageIDs),
		)

		sendNotification(ctx, b, yourUserID, notification)
		return
	}
}

func getChatTitle(chat models.Chat) string {
	if chat.Title != "" {
		return chat.Title
	}
	if chat.Username != "" {
		return "@" + chat.Username
	}
	name := chat.FirstName
	if chat.LastName != "" {
		name += " " + chat.LastName
	}
	if name != "" {
		return name
	}
	return fmt.Sprintf("Chat %d", chat.ID)
}

func getUserName(user *models.User) string {
	if user.Username != "" {
		return "@" + user.Username
	}
	name := user.FirstName
	if user.LastName != "" {
		name += " " + user.LastName
	}
	if name != "" {
		return name
	}
	return fmt.Sprintf("User %d", user.ID)
}

func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

func sendNotification(ctx context.Context, b *bot.Bot, userID int64, text string) {
	_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    userID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
	})
}
