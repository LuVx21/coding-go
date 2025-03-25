package main

import (
	"log"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Test_telegram(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI("YOUR_BOT_TOKEN")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			bot.Send(msg)
		}
	}
}

func Test_telegram_02(t *testing.T) {
	b, err := tb.NewBot(tb.Settings{
		Token:  "YOUR_BOT_TOKEN",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello, world!")
	})

	b.Start()
}
