package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Question struct {
	Text   string
	Answer string
}

type Questions struct {
	Questions map[string]*Question
}

func main() {
	doc, _ := ioutil.ReadFile("questions.json")
	var questions Questions
	err := json.Unmarshal(doc, &questions)
	if err != nil {
		log.Fatal("Invalid settings format:", err)
	}
	newQuestion := &Question{
		Text:   "allo",
		Answer: "hi",
	}
	log.Printf("\n%+v", questions.Questions["1"])
	questions.Questions["2"] = newQuestion
	questions.Questions["1"].Text = "hihihihi"
	questions.Questions["1"].Answer = "fsfsad"
	log.Printf("\n%+v", questions.Questions["1"])
	output, _ := json.MarshalIndent(questions, "", " ")

	ioutil.WriteFile("questions.json", output, 0644)

	log.Printf("\n%s", output)
}
