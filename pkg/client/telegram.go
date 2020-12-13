package client

import (
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

type TelegramBotClient struct {
	ids    []int64
	logger *zap.Logger
	client *telegram.BotAPI
}

func NewTelegram(token string, ids []int64, logger *zap.Logger) (*TelegramBotClient, error) {
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("could not create telegram bot: %w", err)
	}

	bot.Debug = true

	return &TelegramBotClient{
		client: bot,
		ids:    ids,
		logger: logger,
	}, nil
}

func (c *TelegramBotClient) Message(text string) error {
	for _, id := range c.ids {
		_, err := c.client.Send(telegram.NewMessage(id, text))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *TelegramBotClient) Update(messages func() []string) error {
	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.client.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("could not initialize telegram callback")
	}

	for update := range updates {

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		c.logger.Info("receive bot message",
			zap.String("message", update.Message.Text),
			zap.String("user", fmt.Sprintf("%s, %s", update.Message.Chat.LastName, update.Message.Chat.FirstName)))

		switch {
		case strings.Contains(strings.ToLower(update.Message.Text), "check"):
			c.client.Send(telegram.NewMessage(update.Message.Chat.ID, "Die aktuellen Angebote sind:"))
			for _, message := range messages() {
				_, err := c.client.Send(telegram.NewMessage(update.Message.Chat.ID, message))
				if err != nil {
					return fmt.Errorf("could not send telegram response: %w", err)
				}
			}
		default:
		}

	}

	return nil
}
