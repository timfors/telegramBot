package main

import (
	"encoding/json"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"strings"
)

type Question struct {
	text   string
	answer string
}

type Questions struct {
	Questions map[string]*Question
}

var editQuestionNum string
var botState string
var questions Questions
var progresses map[int64]int
var commands map[string]string

func main() {
	editQuestionNum = "0"
	botState = "idle"
	file, err := ioutil.ReadFile("questions.json")
	err = json.Unmarshal(file, &questions)
	progresses = map[int64]int{}
	commands = map[string]string{"/showAll": "show all the questions", "/addQuestion": "add question", "/removeQuestion": "remove question", "/changeQuestion": "changes question"}
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
				question, _ := questions.Questions["0"]
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, question.text)
				bot.Send(msg)
			}

		}
		if update.Message.Chat.ID == 322726399 {
			AdminAnswer(bot, update)
		} else {
			SimpleAnswer(bot, update)
		}
		log.Printf("\nbotState: %s\n", botState)

	}
}

func SimpleAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	stage := progresses[update.Message.Chat.ID]
	answ := questions.Questions[string(stage)].answer
	if strings.ToLower(update.Message.Text) == strings.ToLower(answ) {
		progresses[update.Message.Chat.ID]++
		if stage-1 < len(questions.Questions) {
			question, _ := questions.Questions[string(stage+1)]
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, question.text)
			bot.Send(msg)
		}
	}
}

func AdminAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Text {
	case "/showAll":
		for i := 1; i <= len(questions.Questions); i++ {
			log.Printf("\n%d\n", i)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, string(i)+". "+questions.Questions[string(i)].text+"\nAnswer: "+questions.Questions[string(i)].answer)
			bot.Send(msg)
		}
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/changeQuestion":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Номер вопроса, уважаемый.")
		botState = "getQuestionNum"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/addQuestion":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Давай вопрос и разойдемся.")
		botState = "addingText"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/?", "/commands":
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
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionText":
			if len(strings.Split(update.Message.Text, " ")) > 0 {
				questions.Questions[editQuestionNum].text = update.Message.Text
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ну а сейчас меняй ответ. (оставь пустым, если не хочешь изменять)")
			bot.Send(msg)
			botState = "editingQuestionAnswer"
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionAnswer":
			if len(strings.Split(update.Message.Text, " ")) > 0 {
				questions.Questions[editQuestionNum].answer = update.Message.Text
			}
			botState = "idle"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Начальник, принимай работу!")
			bot.Send(msg)
			SaveJSON()
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingText":
			questions.Questions[string(len(questions.Questions)+1)].text = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ответик в студию!")
			botState = "addingAnswer"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingAnswer":
			questions.Questions[string(len(questions.Questions))].answer = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хотово!")
			botState = "idle"
			bot.Send(msg)
			SaveJSON()
			log.Printf("\nbotState: %s\n", botState)
			break

		case "idle":
			SimpleAnswer(bot, update)
		}
	}
}

func SaveJSON() {
	output, _ := json.Marshal(questions)
	ioutil.WriteFile("question.json", output, 0644)
}
