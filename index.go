package main

import (
	"encoding/json"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	editQuestionNum := "0"
	botState := "idle"
	file, err := ioutil.ReadFile("questions.json")
	questions := map[string]map[string]string{}
	err = json.Unmarshal(file, &questions)
	progresses := map[int64]int{}
	commands := map[string]string{"/showAll": "show all the questions", "/addQuestion": "add question", "/removeQuestion": "remove question", "/changeQuestion": "changes question"}
	bot, err := tgbotapi.NewBotAPI("866951564:AAHdOQgN6ZrypN0uraxAijmrDmDGln7bw48")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)

	for update := range updates {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
			case "resetProgress":
				progresses[update.Message.Chat.ID] = 1
				question, _ := questions["0"]
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, question["text"])
				bot.Send(msg)
			}

		}
		if update.Message.Chat.ID == 322726399 {
			go AdminAnswer(bot, update, progresses, questions, botState, commands, editQuestionNum)
		} else {
			go SimpleAnswer(bot, update, progresses, questions, botState)
		}
		log.Printf("botState: %s\n", botState)

	}
}

func SimpleAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update, progresses map[int64]int, questions map[string]map[string]string, botState string) {
	stage := progresses[update.Message.Chat.ID]
	answ := questions[string(stage)]["answ"]
	if strings.ToLower(update.Message.Text) == strings.ToLower(answ) {
		progresses[update.Message.Chat.ID]++
		if stage-1 < len(questions) {
			question, _ := questions[string(stage+1)]
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, question["text"])
			bot.Send(msg)
		}
	}
}

func AdminAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update, progresses map[int64]int, questions map[string]map[string]string, botState string, commands map[string]string, editQuestionNum string) {
	switch update.Message.Text {
	case "/showAll":
		for i := 1; i <= len(questions); i++ {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(i)+". "+questions[string(i)]["text"]+"\nAnswer: "+questions[string(i)]["answ"])
			bot.Send(msg)
		}
		break

	case "/editQuest":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Номер вопроса, уважаемый.")
		botState = "getQuestionNum"
		bot.Send(msg)
		break

	case "/addQuestion":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Давай вопрос и разойдемся.")
		botState = "addingText"
		bot.Send(msg)
		break

	case "/commands":
	case "/?":
		for command, description := range commands {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, command+" - "+description)
			bot.Send(msg)
		}
		break

	default:
		switch botState {

		case "getQuestionNum":
			editQuestionNum = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь на какой текст меняем? (оставь пустым, если не хочешь изменять)")
			bot.Send(msg)
			botState = "editingQuestionText"
			break

		case "editingQuestionText":
			if len(strings.Split(update.Message.Text, " ")) > 0 {
				questions[editQuestionNum]["text"] = update.Message.Text
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ну а сейчас меняй ответ. (оставь пустым, если не хочешь изменять)")
			bot.Send(msg)
			botState = "editingQuestionAnswer"
			break

		case "editingQuestionAnswer":
			if len(strings.Split(update.Message.Text, " ")) > 0 {
				questions[editQuestionNum]["answ"] = update.Message.Text
			}
			botState = "idle"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Начальник, принимай работу!")
			bot.Send(msg)
			SaveJSON(questions)
			break

		case "addingText":
			questions[string(len(questions))]["text"] = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ответик в студию!")
			botState = "addingAnswer"
			bot.Send(msg)
			break

		case "addingAnswer":
			questions[string(len(questions))]["answer"] = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хотово!")
			botState = "idle"
			bot.Send(msg)
			SaveJSON(questions)
			break

		case "idle":
			SimpleAnswer(bot, update, progresses, questions, botState)
		}
	}
}

func SaveJSON(questions map[string]map[string]string) {
	output, _ := json.Marshal(questions)
	ioutil.WriteFile("question.json", output, 0644)
}
