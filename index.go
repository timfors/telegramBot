package main

import (
	"encoding/json"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Text   string
	Answer string
	Hint   string
}

type Data struct {
	Questions map[string]*Question
	Answers   map[string]string
	HintTimer int
}

var incorrectAnsw map[int]string
var editQuestionNum string
var editAnswerNum string
var botState string
var data Data
var progresses map[int64]int
var commands map[string]string

func main() {
	editQuestionNum = "0"
	botState = "idle"
	file, err := ioutil.ReadFile("data.json")
	err = json.Unmarshal(file, &data)
	progresses = map[int64]int{}
	commands = map[string]string{"/showQ": "show all the questions",
		"/addQ": "add question", "/removeLastQ": "remove last question",
		"/changeQ": "changes question", "/changeA": "change bot answer",
		"/addA": "add bot answer", "/removeLastA": "remove last bot answer",
		"/showA": "show all the bot answer", "/showH": "show hint timer",
		"/changeHintTimer": "change hint timer"}
	for _, question := range data.Questions {
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
		if update.Message.Chat.ID == 322726399 || update.Message.Chat.ID == 479731828 {
			AdminAnswer(bot, update)
		} else {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "reset_progress", "start":
					progresses[update.Message.Chat.ID] = 1
					go hintTimer(bot, update.Message.Chat.ID, progresses[update.Message.Chat.ID])
					question := data.Questions["1"]
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
	answ := data.Questions[strconv.Itoa(stage)].Answer
	if strings.ToLower(update.Message.Text) == strings.ToLower(answ) {
		progresses[update.Message.Chat.ID]++
		go hintTimer(bot, update.Message.Chat.ID, progresses[update.Message.Chat.ID])
		if stage-1 < len(data.Questions) {
			question, _ := data.Questions[strconv.Itoa(stage+1)]
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, question.Text)
			bot.Send(msg)
		}
	} else {
		answerNum := rand.Intn(len(data.Answers)) + 1
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, data.Answers[strconv.Itoa(answerNum)])
		bot.Send(msg)
	}
}

func AdminAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Text {

	case "/changeHintTimer":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выкладывай, сколько ждать перед подсказкой?")
		botState = "changeHintTimer"
		bot.Send(msg)

	case "/removeLastA":
		delete(data.Answers, strconv.Itoa(len(data.Answers)))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пока, ответик!")
		SaveJSON()
		bot.Send(msg)
		break

	case "/removeLastQ":
		delete(data.Questions, strconv.Itoa(len(data.Questions)))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Минус бомжара!")
		SaveJSON()
		bot.Send(msg)
		break

	case "/showA":
		for i := 1; i <= len(data.Answers); i++ {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.Itoa(i)+". "+data.Answers[strconv.Itoa(i)])
			bot.Send(msg)
		}
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/showQ":
		for i := 1; i <= len(data.Questions); i++ {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.Itoa(i)+". "+data.Questions[strconv.Itoa(i)].Text+"\nAnswer: "+data.Questions[strconv.Itoa(i)].Answer+"\nHint: "+data.Questions[strconv.Itoa(i)].Hint)
			bot.Send(msg)
		}
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/showH":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ожидание перед подсказкой: "+strconv.Itoa(data.HintTimer))
		bot.Send(msg)
		break

	case "/changeA":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Номер ответа бота, уважаемый.")
		botState = "getAnswerNum"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/changeQ":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Номер вопроса, уважаемый.")
		botState = "getQuestionNum"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/addA":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Давай ответ и чики брики.")
		botState = "addingAnswerText"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/addQ":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Давай вопрос и разойдемся.")
		botState = "addingQuestionText"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/?", "/commands":
		for command, description := range commands {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, command+" - "+description)
			bot.Send(msg)
		}
		break

	case "/reset_progress", "/start":
		progresses[update.Message.Chat.ID] = 1
		go hintTimer(bot, update.Message.Chat.ID, progresses[update.Message.Chat.ID])
		question := data.Questions["1"]
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, question.Text)
		bot.Send(msg)

	default:
		switch botState {
		case "changeHintTimer":
			num, err := strconv.ParseInt(update.Message.Text, 10, 64)
			data.HintTimer = int(num)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Это не число!Сворачиваемся!")
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Окес!")
				SaveJSON()
				bot.Send(msg)
			}
			botState = "idle"
		case "getAnswerNum":
			editAnswerNum = update.Message.Text
			num, _ := strconv.ParseInt(editAnswerNum, 10, 64)
			if num > int64(len(data.Answers)) {
				botState = "idle"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нету такого ответа!")
				bot.Send(msg)
			} else {
				botState = "editingAnswerText"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь на какой текст меняем? (пиши ОТСТАНЬ, если не хочешь изменять)")
				bot.Send(msg)
			}
			log.Printf("\nbotState: %s\n", botState)
			break

		case "getQuestionNum":
			editQuestionNum = update.Message.Text
			num, _ := strconv.ParseInt(editQuestionNum, 10, 64)
			if num > int64(len(data.Questions)) {
				botState = "idle"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вопрос под таким номером отсутсвует")
				bot.Send(msg)
			} else {
				botState = "editingQuestionText"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Теперь на какой текст меняем? (пиши ОТСТАНЬ, если не хочешь изменять)")
				bot.Send(msg)
			}
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingAnswerText":
			if strings.ToLower(update.Message.Text) == "отстань" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отмена операции!")
				bot.Send(msg)
			} else {
				data.Answers[editAnswerNum] = update.Message.Text
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ну все, дело сделано!")
			bot.Send(msg)
			botState = "idle"
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionText":
			if strings.ToLower(update.Message.Text) == "отстань" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Допустим не меняем текст!")
				bot.Send(msg)
			} else {
				data.Questions[editQuestionNum].Text = update.Message.Text
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ну а сейчас меняй ответ. (пиши ОТСТАНЬ, если не хочешь изменять)")
			bot.Send(msg)
			botState = "editingQuestionAnswer"
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionAnswer":
			if strings.ToLower(update.Message.Text) == "отстань" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ладно, оставляем ответ.")
				bot.Send(msg)
			} else {
				data.Questions[editQuestionNum].Answer = update.Message.Text
			}
			botState = "idle"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Начальник, принимай работу!")
			bot.Send(msg)
			SaveJSON()
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionHint":
			if strings.ToLower(update.Message.Text) == "отстань" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ладно, оставляем ответ.")
				bot.Send(msg)
			} else {
				data.Questions[editQuestionNum].Answer = update.Message.Text
			}
			botState = "idle"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Начальник, принимай работу!")
			bot.Send(msg)
			SaveJSON()
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingAnswerText":
			num := len(data.Answers) + 1
			data.Answers[strconv.Itoa(num)] = ""
			data.Answers[strconv.Itoa(num)] = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вроде все.")
			botState = "idle"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingQuestionText":
			num := len(data.Questions) + 1
			data.Questions[strconv.Itoa(num)] = &Question{Text: "", Answer: ""}
			data.Questions[strconv.Itoa(num)].Text = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ответик в студию!")
			botState = "addingQuestionAnswer"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingQuestionAnswer":
			data.Questions[strconv.Itoa(len(data.Questions))].Answer = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Давай теперь подсказку!(пиши ОТСТАНЬ, если не хочешь добавлять)")
			botState = "addingQuestionHint"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingQuestionHint":
			data.Questions[strconv.Itoa(len(data.Questions))].Answer = update.Message.Text
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вопрос создан!")
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
	output, _ := json.MarshalIndent(data, "", " ")
	ioutil.WriteFile("data.json", output, 0644)
}

func hintTimer(bot *tgbotapi.BotAPI, chatId int64, progress int) {
	time.Sleep(time.Duration(data.HintTimer) * time.Minute)
	if progresses[chatId] == progress {
		msg := tgbotapi.NewMessage(chatId, "Подсказка:"+data.Questions[strconv.Itoa(progress)].Hint)
		bot.Send(msg)
	}
}
