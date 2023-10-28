package handlers

import (
	"errors"
	"log"
	"me885/fintech-or-furniture/quiz"
	"me885/fintech-or-furniture/quiz/database"
	"net/http"
	"regexp"
	"strconv"
	"text/template"

	"github.com/google/uuid"
)

type Context struct {
	DB *database.SQLiteRepository
}

func RootPage(writer http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("./templates/index.html"))

	template.Execute(writer, nil)
}

func (context Context) NewGame(writer http.ResponseWriter, request *http.Request) {

	game, err := context.DB.CreateGame(request.PostFormValue("name"))

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	question, err := context.DB.GetQuestionById(int64(game.QuestionsAnswered + 1))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
		return
	}

	cookie := http.Cookie{Name: "sessionId", Value: game.Id.String(), HttpOnly: true, SameSite: http.SameSiteLaxMode, Path: "/"}

	http.SetCookie(writer, &cookie)

	template := template.Must(template.ParseFiles("./templates/quizQuestion.html"))
	template.Execute(writer, quiz.QuestionPageStruct{Question: *question, Game: *game})
}

func (context Context) Answer(writer http.ResponseWriter, request *http.Request) {

	game, err := getGameIfAuthed(request, context.DB)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	if !game.InProgress {
		http.Error(writer, "Game is finished. Connot answer more questions", http.StatusUnauthorized)
		return
	}

	regex := regexp.MustCompile(`answer/([0-9]*)/`)
	questionId, err := strconv.ParseInt(regex.FindStringSubmatch(request.URL.Path)[1], 10, 64)
	if err != nil {
		http.Error(writer, "QuestionId must specified in URL path", http.StatusBadRequest)
		return
	}

	question, err := context.DB.GetQuestionById(questionId)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	answer := request.URL.Query().Get("answer")

	wasCorrect, err := quiz.HandleAnswer(answer, *question, game)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if !quiz.IsGameComplete(game) {
		template := template.Must(template.ParseFiles("./templates/nextQuestion.html"))
		template.Execute(writer, quiz.NextQuestionModalStruct{Correct: wasCorrect, Score: game.Score})

	} else {
		template := template.Must(template.ParseFiles("./templates/endPage.html"))
		template.Execute(writer, game)
	}

	context.DB.UpdateGame(game)
}

func (context Context) NextQuestion(writer http.ResponseWriter, request *http.Request) {

	game, err := getGameIfAuthed(request, context.DB)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	if !game.InProgress {
		http.Error(writer, "Game is finished. Connot answer more questions", http.StatusUnauthorized)
		return
	}

	question, err := context.DB.GetQuestionById(int64(game.QuestionsAnswered + 1))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	template := template.Must(template.ParseFiles("./templates/quizQuestion.html"))
	template.Execute(writer, quiz.QuestionPageStruct{Question: *question, Game: *game})
}

func (context Context) Leaderboard(writer http.ResponseWriter, request *http.Request) {
	games, err := context.DB.TopTenCompletedGames()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	template := template.Must(template.ParseFiles("./templates/leaderboard.html"))
	template.Execute(writer, games)
}

func (context Context) EndPage(writer http.ResponseWriter, request *http.Request) {
	game, err := getGameIfAuthed(request, context.DB)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	template := template.Must(template.ParseFiles("./templates/endPage.html"))
	template.Execute(writer, game)
}

func getGameIfAuthed(request *http.Request, db *database.SQLiteRepository) (*quiz.Game, error) {
	cookie, err := request.Cookie("sessionId")
	if err != nil {
		return nil, errors.New("sessionId cookie required")
	}

	gameId, err := uuid.Parse(cookie.Value)
	if err != nil {
		return nil, errors.New("sessionId should be valid uuid")
	}

	game, err := db.GetGameById(gameId)
	if err != nil {
		return nil, errors.New("sessionId not found")
	}

	return game, nil
}
