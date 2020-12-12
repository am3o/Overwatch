package client

import (
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramBotClient struct {
	ids []int64
	client *telegram.BotAPI
}

func NewTelegram(token string, ids []int64) (*TelegramBotClient, error) {
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("could not create telegram bot: %w", err)
	}

	bot.Debug = true

	return &TelegramBotClient{
		client: bot,
		ids: ids,
	}, nil
}

func (c *TelegramBotClient) Message(text string) error{
	for _, id := range c.ids {
		_, err := c.client.Send(telegram.NewMessage(id, text))
		if err != nil {
			return err
		}
	}
	return nil
}
