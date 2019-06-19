package main

import (
	"github.com/Syfaro/telegram-bot-api"
	"log"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("866951564:AAHdOQgN6ZrypN0uraxAijmrDmDGln7bw48")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)

	for {
		select {
		case update := <-updates:
			userName := update.Message.From.UserName
			chatID := update.Message.Chat.ID
			text := update.Message.Text
			log.Printf("[%s] %d %s", userName, chatID, text)
			reply := "Nu darova"
			msg := tgbotapi.NewMessage(chatID, reply)
			bot.Send(msg)
		}
	}
}
