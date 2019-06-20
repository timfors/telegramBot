package main

import (
	"encoding/json"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Question struct {
	Text   string
	Answer string
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
	commands = map[string]string{"/show": "show all the questions", "/add": "add question", "/removeLast": "remove last question", "/change": "changes question"}
	for _, question := range questions.Questions {
		log.Printf("\n%+v\n", question)
	}
	bot, err := tgbotapi.NewBotAPI("866951564:AAHdOQgN6ZrypN0uraxAijmrDmDGln7bw48")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)

	for update := range updates {
		if update.Message.Chat.ID == 322726399 {
			AdminAnswer(bot, update)
		} else {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
				case "resetProgress":
					progresses[update.Message.Chat.ID] = 1
					question := questions.Questions["1"]
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, question.Text)
					bot.Send(msg)
				}

			} else {
				SimpleAnswer(bot, update)
			}
		}

		log.Printf("\nbotState: %s\n", botState)

	}
}

func SimpleAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if progresses[update.Message.Chat.ID] == 0 {
		progresses[update.Message.Chat.ID] = 1
	}
	stage := progresses[update.Message.Chat.ID]
	answ := questions.Questions[strconv.Itoa(stage)].Answer
	if strings.ToLower(update.Message.Text) == strings.ToLower(answ) {
		progresses[update.Message.Chat.ID]++
		if stage-1 < len(questions.Questions) {
			question, _ := questions.Questions[strconv.Itoa(stage+1)]
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, question.Text)
			bot.Send(msg)
		}
	}
}

func AdminAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Text {
	case "/removeLast":
		delete(questions.Questions, strconv.Itoa(len(questions.Questions)))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Минус бомжара!")
		SaveJSON()
		bot.Send(msg)
		break

	case "/show":
		for num, question := range questions.Questions {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, num+". "+question.Text+"\nAnswer: "+question.Answer)
			bot.Send(msg)
		}
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/change":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Номер вопроса, уважаемый.")
		botState = "getQuestionNum"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/add":
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

	case "/start":
	case "/resetProgress":
		progresses[update.Message.Chat.ID] = 1
		question := questions.Questions["1"]
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, question.Text)
		bot.Send(msg)

	default:
		switch botState {

		case "getQuestionNum":
			editQuestionNum = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь на какой текст меняем? (оставь пустым, если не хочешь изменять)")
			bot.Send(msg)
			num, _ := strconv.ParseInt(editQuestionNum, 10, 64)
			if num > int64(len(questions.Questions)) {
				botState = "addingText"
			} else {
				botState = "editingQuestionText"
			}
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionText":
			if len(strings.Split(update.Message.Text, " ")) > 0 {
				questions.Questions[editQuestionNum].Text = update.Message.Text
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ну а сейчас меняй ответ. (оставь пустым, если не хочешь изменять)")
			bot.Send(msg)
			botState = "editingQuestionAnswer"
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionAnswer":
			if len(strings.Split(update.Message.Text, " ")) > 0 {
				questions.Questions[editQuestionNum].Answer = update.Message.Text
			}
			botState = "idle"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Начальник, принимай работу!")
			bot.Send(msg)
			SaveJSON()
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingText":
			num := len(questions.Questions) + 1
			questions.Questions[strconv.Itoa(num)] = &Question{Text: "", Answer: ""}
			questions.Questions[strconv.Itoa(num)].Text = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ответик в студию!")
			botState = "addingAnswer"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingAnswer":
			questions.Questions[strconv.Itoa(len(questions.Questions))].Answer = update.Message.Text
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
	output, _ := json.MarshalIndent(questions, "", " ")
	ioutil.WriteFile("questions.json", output, 0644)
}
