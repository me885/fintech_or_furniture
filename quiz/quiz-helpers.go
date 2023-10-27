package quiz

import (
	"errors"
)

func HandleAnswer(answer string, question Question, game *Game) (bool, error) {

	game.QuestionsAnswered++

	if answer == "Fintech" && question.Answer == Fintech {
		game.Score++
		return true, nil
	} else if answer == "Furniture" && question.Answer == Furniture {
		game.Score++
		return true, nil
	} else if answer == "Fintech" && question.Answer == Furniture {
		return false, nil
	} else if answer == "Furniture" && question.Answer == Fintech {
		return false, nil
	}

	return false, errors.New("Answer should be 'Fintech' or 'Furniture'")
}

func IsGameComplete(game *Game) bool {
	if game.QuestionsAnswered < 10 {
		return false

	} else {
		game.InProgress = false

		return true
	}
}
