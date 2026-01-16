package main

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func handleUpdate(ctx context.Context, b *bot.Bot, update *models.Update, store *MessageStore, yourUserID int64) {
	if update.Message != nil && update.Message.Text != "" {
		if update.Message.Text == "/stats" &&
			update.Message.From != nil &&
			update.Message.From.ID == yourUserID {

			count := store.Count()

			sendNotification(
				ctx,
				b,
				yourUserID,
				fmt.Sprintf("📊 Сообщений в хранилище: <b>%d</b>", count),
			)
			return
		}
	}

	if update.BusinessMessage != nil {
		msg := update.BusinessMessage
		if msg.From != nil && msg.From.ID == yourUserID {
			return
		}
		if msg.Text != "" {
			store.Save(msg.BusinessConnectionID, msg.Chat.ID, msg.ID, msg.Text)
		}
		return
	}

	if update.EditedBusinessMessage != nil {
		edited := update.EditedBusinessMessage
		if edited.From != nil && edited.From.ID == yourUserID {
			return
		}

		chatTitle := getChatTitle(edited.Chat)
		userName := getUserName(edited.From)

		// Получаем оригинальное сообщение
		originalText, exists := store.Get(
			edited.BusinessConnectionID,
			edited.Chat.ID,
			edited.ID,
		)

		var notification string
		if exists && originalText != "" {
			if originalText == edited.Text {
				// Текст не изменился
				notification = fmt.Sprintf(
					"✏️ <b>%s</b> | %s\n"+
						"━━━━━━━━━━━━━━━\n"+
						"<i>Сообщение отредактировано (текст не изменился)</i>",
					userName,
					chatTitle,
				)
			} else {
				// Показываем diff с подсветкой
				diffHTML := generatePrettyDiff(originalText, edited.Text)

				notification = fmt.Sprintf(
					"✏️ <b>%s</b> | %s\n"+
						"━━━━━━━━━━━━━━━\n"+
						"%s",
					userName,
					chatTitle,
					diffHTML,
				)
			}
		} else {
			// Оригинал не найден
			notification = fmt.Sprintf(
				"✏️ <b>%s</b> | %s\n"+
					"━━━━━━━━━━━━━━━\n"+
					"%s",
				userName,
				chatTitle,
				escapeHTML(edited.Text),
			)
		}

		sendNotification(ctx, b, yourUserID, notification)
		store.Save(edited.BusinessConnectionID, edited.Chat.ID, edited.ID, edited.Text)
		return
	}

	if update.DeletedBusinessMessages != nil {
		deleted := update.DeletedBusinessMessages
		bizConnID := deleted.BusinessConnectionID
		chatID := deleted.Chat.ID
		chatTitle := getChatTitle(deleted.Chat)

		for _, messageID := range deleted.MessageIDs {
			originalText, exists := store.Get(bizConnID, chatID, messageID)

			if !exists {
				continue
			}

			if originalText != "" {
				var notification = fmt.Sprintf(
					"🗑 <b>%s</b>\n"+
						"━━━━━━━━━━━━━━━\n"+
						"%s",
					chatTitle,
					escapeHTML(originalText),
				)
				sendNotification(ctx, b, yourUserID, notification)
				store.Delete(bizConnID, chatID, messageID)
			}
		}
	}
}
