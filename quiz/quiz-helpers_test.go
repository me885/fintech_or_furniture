package quiz

import (
	"testing"

	"github.com/google/uuid"
)

type HandleAnswerTest struct {
	answer                    string
	questionAnswer            Answer
	questionsAnswered         int64
	expectedWasCorrect        bool
	expectedQuestionsAnswered int64
}

var handleAnswerHappyPaths = []HandleAnswerTest{
	{"Fintech", Fintech, 4, true, 5},
	{"Fintech", Furniture, 4, false, 5},
	{"Furniture", Fintech, 4, false, 5},
	{"Furniture", Furniture, 4, true, 5},
}

func TestHandleAnswer(t *testing.T) {
	for _, v := range handleAnswerHappyPaths {
		answer := v.answer
		question := Question{Id: 1, Question: "google", Answer: v.questionAnswer}
		game := &Game{Id: uuid.New(), PlayerName: "bob", QuestionsAnswered: v.questionsAnswered, Score: 4, InProgress: true}

		wasCorrect, err := HandleAnswer(answer, question, game)
		if wasCorrect != v.expectedWasCorrect || game.QuestionsAnswered != v.expectedQuestionsAnswered || err != nil {
			t.Fatal(wasCorrect, game.QuestionsAnswered, err, v)
		}
	}
}

func TestHandleAnswer_InvalidAnswer(t *testing.T) {
	answer := "apple"
	question := Question{Id: 1, Question: "google", Answer: Fintech}
	game := &Game{Id: uuid.New(), PlayerName: "bob", QuestionsAnswered: 4, Score: 4, InProgress: true}

	wasCorrect, err := HandleAnswer(answer, question, game)
	if err == nil {
		t.Fatal(wasCorrect)
	}
}
