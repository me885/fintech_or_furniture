package quiz

import "github.com/google/uuid"

type Answer int64

const (
	Fintech   Answer = 0
	Furniture Answer = 1
)

type Question struct {
	Id       int64
	Question string
	Answer   Answer
}

type Game struct {
	Id                uuid.UUID
	PlayerName        string
	QuestionsAnswered int64
	Score             int64
	InProgress        bool
}

type QuestionPageStruct struct {
	Question Question
	Game     Game
}

type NextQuestionModalStruct struct {
	Correct bool
	Score   int64
}
