package main

import (
	"context"
	"errors"
	"github.com/Syfaro/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	Number int
	Text   string
	Answer []string
	Hint   []string
}

type BotAnswer struct {
	Number int
	Text   string
}

type HintTimer struct {
	Time int
}
type Token struct {
	Token string
}

type Progress struct {
	Id       string
	Progress int
}

var botState string

var answers []*BotAnswer
var newAnswer *BotAnswer

var questions []*Question
var newQuestion *Question

var progresses []*Progress
var newProgress Progress

var hintTimer HintTimer

var token Token

var commands map[string]string
var client *mongo.Client

func ConnectToDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(""))
	if err != nil {
		log.Fatal(err)
	}

	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connection to DB done!")
	return client
}

func GetCollection(name string) *mongo.Collection {
	return client.Database("Data").Collection(name)
}

func AddProgress(progress *Progress) {
	_, err := GetCollection("Progresses").InsertOne(context.TODO(), progress)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Added!")
}
func ChangeProgress(progress *Progress) {
	filter := bson.D{{"id", progress.Id}}
	update := bson.D{{"$set", progress}}
	_, err := GetCollection("Progresses").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Changed!")
}

func UpdateProgresses() []*Progress {
	var progresses []*Progress
	options := options.Find()
	filter := bson.M{}
	cur, err := GetCollection("Progresses").Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(context.TODO()) {
		var elem Progress
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		progresses = append(progresses, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	log.Println("Updated: Progresses")
	return progresses
}

func FindProgress(userId string) (*Progress, error) {
	for _, progress := range progresses {
		if progress.Id == userId {
			return progress, nil
		}
	}
	return &Progress{}, errors.New("No such progress")
}
func ChangeToken(token Token) {
	filter := bson.D{}
	update := bson.D{{"$set", token}}
	_, err := GetCollection("Token").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Changed")
}

func UpdateToken() Token {
	var token Token
	filter := bson.D{}
	err := GetCollection("Token").FindOne(context.TODO(), filter).Decode(&token)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Updated: Token")
	return token
}

func ChangeHintTimer(hintTimer HintTimer) {
	filter := bson.D{}
	update := bson.D{{"$set", hintTimer}}
	_, err := GetCollection("HintTimer").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Changed!")
}

func UpdateHintTimer() HintTimer {
	var hintTimer HintTimer
	filter := bson.D{}
	err := GetCollection("HintTimer").FindOne(context.TODO(), filter).Decode(&hintTimer)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Updated: HintTimer")
	return hintTimer
}

func AddBotAnswer(answer *BotAnswer) {
	_, err := GetCollection("BotAnswers").InsertOne(context.TODO(), answer)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Added!")
}

func ChangeBotAnswer(answer *BotAnswer) {
	filter := bson.D{{"number", answer.Number}}
	update := bson.D{{"$set", answer}}
	_, err := GetCollection("BotAnswers").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Changed")
}

func RemoveLastBotAnswer() {
	answerCount, err := GetCollection("BotAnswers").EstimatedDocumentCount(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{{"number", int(answerCount)}}
	_, err = GetCollection("BotAnswers").DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Removed!")
}

func UpdateBotAnswers() []*BotAnswer {
	var answers []*BotAnswer
	options := options.Find()
	filter := bson.D{}

	cur, err := GetCollection("BotAnswers").Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem BotAnswer
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		answers = append(answers, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.TODO())
	log.Println("Updated: BotAnswers")
	return answers
}

func FindBotAnswer(num int) (*BotAnswer, error) {
	for _, answer := range answers {
		if answer.Number == num {
			return answer, nil
		}
	}
	return &BotAnswer{}, errors.New("No such bot answers")
}

func AddQuestion(question *Question) {
	_, err := GetCollection("Questions").InsertOne(context.TODO(), question)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Added!")
}

func ChangeQuestion(question *Question) {
	filter := bson.D{{"number", question.Number}}
	update := bson.D{{"$set", question}}
	_, err := GetCollection("Questions").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Changed!")
}

func RemoveLastQuestion() {
	questionCount, err := GetCollection("Questions").EstimatedDocumentCount(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	filter := bson.D{{"number", int(questionCount)}}
	_, err = GetCollection("BotAnswers").DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Removed!")
}

func UpdateQuestions() []*Question {
	var questions []*Question
	options := options.Find()
	filter := bson.D{}

	cur, err := GetCollection("Questions").Find(context.TODO(), filter, options)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem Question
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		questions = append(questions, &elem)
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
	}
	cur.Close(context.TODO())
	log.Println("Updated: Questions")
	return questions
}

func FindQuestion(num int) (*Question, error) {
	for _, question := range questions {
		if question.Number == num {
			return question, nil
		}
	}
	return &Question{}, errors.New("No such questions")
}

func SetHintTimer(bot *tgbotapi.BotAPI, chatId int64, userName string, progress int, hintCount int) {
	time.Sleep(time.Duration(hintTimer.Time) * time.Minute)
	currentProgress, _ := FindProgress(userName)
	question, _ := FindQuestion(progress)
	if currentProgress.Progress == progress {
		msg := tgbotapi.NewMessage(chatId, question.Hint[hintCount])
		bot.Send(msg)
	}
	if hintCount < len(question.Hint) {
		hintCount++
		SetHintTimer(bot, chatId, userName, progress, hintCount)
	}
}

func TokenGenerator(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(57)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

func CheckAnswer(userAnswer string, correctAnswers []string) bool {
	for _, correctAnswer := range correctAnswers {
		if strings.ToLower(strings.TrimSpace(userAnswer)) == strings.ToLower(strings.TrimSpace(correctAnswer)) {
			return true
		}
	}
	return false
}

func main() {
	botState = "idle"
	client = ConnectToDB()
	questions = UpdateQuestions()
	progresses = UpdateProgresses()
	answers = UpdateBotAnswers()
	token = UpdateToken()
	hintTimer = UpdateHintTimer()

	commands = map[string]string{"/showQ": "show all the questions",
		"/addQ": "add question", "/removeLastQ": "remove last question",
		"/changeQ": "changes question", "/changeA": "change bot answer",
		"/addA": "add bot answer", "/removeLastA": "remove last bot answer",
		"/showA": "show all the bot answer", "/showH": "show hint timer",
		"/changeHintTimer": "change hint timer", "/token": "generate new token",
		"/showT": "show the token"}
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)

	for update := range updates {
		userName := update.Message.Chat.UserName
		userId := update.Message.Chat.ID
		if userName == "timkhab" || userId == 479731828 {
			AdminAnswer(bot, update)
		} else {
			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "reset_progress", "start":
					newProgress, err := FindProgress(userName)
					newProgress = &Progress{userName, 1}
					if err != nil {
						AddProgress(newProgress)
					} else {
						ChangeProgress(newProgress)
					}
					progresses = UpdateProgresses()
					question, err1 := FindQuestion(1)
					if err1 != nil {
						log.Fatal(err)
					}
					msg := tgbotapi.NewMessage(userId, question.Text)
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
	input := update.Message.Text
	userName := update.Message.Chat.UserName
	userId := update.Message.Chat.ID
	progress, _ := FindProgress(userName)
	if progress.Progress == 1 {
		if input == token.Token {
			progress.Progress++
			ChangeProgress(progress)
			go SetHintTimer(bot, userId, userName, progress.Progress, 0)
			question, err := FindQuestion(2)
			if err != nil {
				log.Fatal(err)
			}
			msg := tgbotapi.NewMessage(userId, question.Text)
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(userId, "Пароль неверный!")
			bot.Send(msg)
		}
	} else {
		stage := progress.Progress
		question, err := FindQuestion(stage)
		answ := question.Answer
		if err != nil {
			log.Fatal(err)
		}
		if CheckAnswer(input, answ) {
			progress, _ = FindProgress(userName)
			progress.Progress++
			ChangeProgress(progress)
			if stage < len(questions) {
				go SetHintTimer(bot, userId, userName, progress.Progress, 0)
				question, err := FindQuestion(stage + 1)
				if err != nil {
					log.Fatal(err)
				}
				msg := tgbotapi.NewMessage(userId, question.Text)
				bot.Send(msg)
			}
		} else {
			answerNum := rand.Intn(len(answers)) + 1
			botAnswer, _ := FindBotAnswer(answerNum)
			msg := tgbotapi.NewMessage(userId, botAnswer.Text)
			bot.Send(msg)
		}
	}
}

func AdminAnswer(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	input := update.Message.Text
	userName := update.Message.Chat.UserName
	userId := update.Message.Chat.ID
	switch input {
	case "/token":
		token.Token = TokenGenerator(10)
		ChangeToken(token)
		msg := tgbotapi.NewMessage(userId, "Новый токен: "+token.Token)
		bot.Send(msg)

	case "/changeHintTimer":
		msg := tgbotapi.NewMessage(userId, "Выкладывай, сколько ждать перед подсказкой?")
		botState = "changeHintTimer"
		bot.Send(msg)

	case "/removeLastA":
		RemoveLastBotAnswer()
		answers = UpdateBotAnswers()
		msg := tgbotapi.NewMessage(userId, "Пока, ответик!")
		bot.Send(msg)
		break

	case "/removeLastQ":
		RemoveLastQuestion()
		questions = UpdateQuestions()
		msg := tgbotapi.NewMessage(userId, "Минус бомжара!")
		bot.Send(msg)
		break

	case "/showA":
		for i := 1; i <= len(answers); i++ {
			botAnswer, _ := FindBotAnswer(i)
			msg := tgbotapi.NewMessage(userId, strconv.Itoa(i)+". "+botAnswer.Text)
			bot.Send(msg)
		}
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/showQ":
		for i := 1; i <= len(questions); i++ {
			question, _ := FindQuestion(i)
			output := strconv.Itoa(i) + ". " + question.Text + "\nAnswers: "
			for _, answer := range question.Answer {
				output += answer + " | "
			}
			output += "\nHints:"
			for _, hint := range question.Hint {
				output += hint + " | "
			}
			msg := tgbotapi.NewMessage(userId, output)
			bot.Send(msg)
		}
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/showT":
		msg := tgbotapi.NewMessage(userId, "Токен: "+token.Token)
		bot.Send(msg)
		break

	case "/showH":
		msg := tgbotapi.NewMessage(userId, "Ожидание перед подсказкой: "+strconv.Itoa(hintTimer.Time))
		bot.Send(msg)
		break

	case "/changeA":
		msg := tgbotapi.NewMessage(userId, "Номер ответа бота, уважаемый.")
		botState = "getAnswerNum"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/changeQ":
		msg := tgbotapi.NewMessage(userId, "Номер вопроса, уважаемый.")
		botState = "getQuestionNum"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/addA":
		newAnswer = &BotAnswer{}
		msg := tgbotapi.NewMessage(userId, "Давай ответ и чики брики.")
		botState = "addingAnswerText"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/addQ":
		newQuestion = &Question{}
		msg := tgbotapi.NewMessage(userId, "Давай вопрос и разойдемся.")
		botState = "addingQuestionText"
		bot.Send(msg)
		log.Printf("\nbotState: %s\n", botState)
		break

	case "/?", "/commands":
		for command, description := range commands {
			msg := tgbotapi.NewMessage(userId, command+" - "+description)
			bot.Send(msg)
		}
		break

	case "/reset_progress", "/start":
		newProgress, err := FindProgress(userName)
		newProgress = &Progress{userName, 1}
		if err != nil {
			AddProgress(newProgress)
		} else {
			ChangeProgress(newProgress)
		}
		progresses = UpdateProgresses()
		go SetHintTimer(bot, userId, userName, newProgress.Progress, 0)
		question, _ := FindQuestion(1)
		msg := tgbotapi.NewMessage(userId, question.Text)
		bot.Send(msg)

	default:
		switch botState {
		case "changeHintTimer":
			num, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(userId, "Это не число!Сворачиваемся!")
				bot.Send(msg)
			} else {
				hintTimer.Time = int(num)
				ChangeHintTimer(hintTimer)
				msg := tgbotapi.NewMessage(userId, "Окес!")
				bot.Send(msg)
			}
			botState = "idle"

		case "getAnswerNum":
			editAnswerNum, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				msg := tgbotapi.NewMessage(userId, "Это не число!Сворачиваемся!")
				bot.Send(msg)
				botState = "idle"
			} else {
				if int(editAnswerNum) > len(answers) {
					botState = "idle"
					msg := tgbotapi.NewMessage(userId, "Нету такого ответа!")
					bot.Send(msg)
				} else {
					botState = "editingAnswerText"
					newAnswer, _ = FindBotAnswer(int(editAnswerNum))
					msg := tgbotapi.NewMessage(userId, "Теперь на какой текст меняем? (пиши -, если не хочешь изменять)")
					bot.Send(msg)
				}
			}
			log.Printf("\nbotState: %s\n", botState)
			break

		case "getQuestionNum":
			editQuestionNum, err := strconv.ParseInt(input, 10, 64)
			if err != nil {
				botState = "idle"
				msg := tgbotapi.NewMessage(userId, "Не число...это...")
				bot.Send(msg)
			} else {
				if editQuestionNum > int64(len(questions)) {
					botState = "idle"
					msg := tgbotapi.NewMessage(userId, "Вопрос под таким номером отсутсвует")
					bot.Send(msg)
				} else {
					newQuestion, _ = FindQuestion(int(editQuestionNum))
					botState = "editingQuestionText"
					msg := tgbotapi.NewMessage(userId, "Теперь на какой текст меняем? (пиши -, если не хочешь изменять)")
					bot.Send(msg)
				}
			}
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingAnswerText":
			msg := tgbotapi.NewMessage(userId, "Ну все, дело сделано!")
			if input == "-" {
				msg := tgbotapi.NewMessage(userId, "Отмена операции!")
				bot.Send(msg)
			} else {
				newAnswer.Text = input
				ChangeBotAnswer(newAnswer)
			}
			answers = UpdateBotAnswers()
			bot.Send(msg)
			botState = "idle"
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionText":
			if input == "-" {
				msg := tgbotapi.NewMessage(userId, "Допустим не меняем текст!")
				bot.Send(msg)
			} else {
				newQuestion.Text = input
			}
			msg := tgbotapi.NewMessage(userId, "Ну а сейчас меняй ответ.Если хочешь несколько , то ставь между ответами | (пиши -, если не хочешь изменять)")
			bot.Send(msg)
			botState = "editingQuestionAnswer"
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionAnswer":
			if input == "-" {
				msg := tgbotapi.NewMessage(userId, "Ладно, оставляем ответ.")
				bot.Send(msg)
			} else {
				newQuestion.Answer = strings.Split(input, "|")
			}
			botState = "editingQuestionHint"
			msg := tgbotapi.NewMessage(userId, "Теперь подсказка.Если хочешь несколько , то ставь между подсказками | (пиши -, если не хочешь изменять)")
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "editingQuestionHint":
			if input == "-" {
				msg := tgbotapi.NewMessage(userId, "Ладно,  подсказку не трогаем.")
				bot.Send(msg)
			} else {
				newQuestion.Hint = strings.Split(input, "|")
			}
			botState = "idle"
			msg := tgbotapi.NewMessage(userId, "Начальник, принимай работу!")
			bot.Send(msg)
			ChangeQuestion(newQuestion)
			questions = UpdateQuestions()
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingAnswerText":
			newAnswer.Number = len(answers) + 1
			newAnswer.Text = input
			AddBotAnswer(newAnswer)
			answers = UpdateBotAnswers()
			msg := tgbotapi.NewMessage(userId, "Вроде все.")
			botState = "idle"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingQuestionText":
			newQuestion.Number = len(questions) + 1
			newQuestion.Text = input
			msg := tgbotapi.NewMessage(userId, "Ответик в студию!Если хочешь несколько , то ставь между ответами |")
			botState = "addingQuestionAnswer"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingQuestionAnswer":
			newQuestion.Answer = strings.Split(input, "|")
			msg := tgbotapi.NewMessage(userId, "Давай теперь подсказку!Если хочешь несколько , то ставь между подсказками |")
			botState = "addingQuestionHint"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break

		case "addingQuestionHint":
			newQuestion.Hint = strings.Split(input, "|")
			AddQuestion(newQuestion)
			questions = UpdateQuestions()
			msg := tgbotapi.NewMessage(userId, "Вопрос создан!")
			botState = "idle"
			bot.Send(msg)
			log.Printf("\nbotState: %s\n", botState)
			break
		case "idle":
			SimpleAnswer(bot, update)
		}
	}
}
