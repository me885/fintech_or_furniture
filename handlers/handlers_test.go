package handlers

import (
	"io"
	"me885/fintech-or-furniture/quiz/database"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRootPage(t *testing.T) {
	handler := http.HandlerFunc(RootPage)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if !strings.Contains(html, "Fintech or Furniture") {
		t.Fatal(html)
	}
}

func TestNewGame(t *testing.T) {
	os.Remove("test.db")

	formdata := url.Values{}
	formdata.Set("name", "testname")

	req, err := http.NewRequest("POST", "/new-game/", strings.NewReader(formdata.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")
	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.NewGame)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if !strings.Contains(html, "Is it a Fintech or Furniture?") {
		t.Fatal(html)
	}

	games, _ := testDb.AllGames()

	if games[0].PlayerName != "testname" {
		t.Fatal(games)
	}
}

func TestAnswer_NoCookie(t *testing.T) {
	os.Remove("test.db")

	req, err := http.NewRequest("POST", "/answer/1/?answer=Fintech", nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")
	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.Answer)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != 401 {
		t.Fatal(resp.Code)
	}
}

func TestAnswer_InvalidAnswer(t *testing.T) {
	os.Remove("test.db")

	answer := "Banana"

	req, err := http.NewRequest("POST", "/answer/1/?answer="+answer, nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")
	game, _ := testDb.CreateGame("testname")

	req.AddCookie(&http.Cookie{Name: "sessionId", Value: game.Id.String()})

	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.Answer)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != 400 {
		t.Fatal(resp.Code)
	}
}

func TestAnswer_Correct(t *testing.T) {
	os.Remove("test.db")

	answer := "Furniture"

	req, err := http.NewRequest("POST", "/answer/1/?answer="+answer, nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")
	game, _ := testDb.CreateGame("testname")

	req.AddCookie(&http.Cookie{Name: "sessionId", Value: game.Id.String()})

	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.Answer)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if !strings.Contains(html, "That's Correct!") {
		t.Fatal(html)
	}

	if !strings.Contains(html, "Your current score is: 1") {
		t.Fatal(html)
	}
}

func TestAnswer_Incorrect(t *testing.T) {
	os.Remove("test.db")

	answer := "Fintech"

	req, err := http.NewRequest("POST", "/answer/1/?answer="+answer, nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")
	game, _ := testDb.CreateGame("testname")

	req.AddCookie(&http.Cookie{Name: "sessionId", Value: game.Id.String()})

	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.Answer)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if !strings.Contains(html, "That's Incorrect!") {
		t.Fatal(html)
	}

	if !strings.Contains(html, "Your current score is: 0") {
		t.Fatal(html)
	}
}

func TestAnswer_LastQuestion(t *testing.T) {
	os.Remove("test.db")

	answer := "Fintech"

	req, err := http.NewRequest("POST", "/answer/1/?answer="+answer, nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")
	game, _ := testDb.CreateGame("testname")

	game.QuestionsAnswered = 9
	game.Score = 8

	testDb.UpdateGame(game)

	req.AddCookie(&http.Cookie{Name: "sessionId", Value: game.Id.String()})

	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.Answer)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if !strings.Contains(html, "You achieved the score of:") {
		t.Fatal(html)
	}

	if !strings.Contains(html, "8/10") {
		t.Fatal(html)
	}
}

func TestNextQuestion_NoCookie(t *testing.T) {
	os.Remove("test.db")

	req, err := http.NewRequest("GET", "/next-question/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")

	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.Answer)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	if resp.Code != 401 {
		t.Fatal(resp.Code)
	}
}

func TestNextQuestion(t *testing.T) {
	os.Remove("test.db")

	req, err := http.NewRequest("GET", "/next-question/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")
	game, _ := testDb.CreateGame("testname")

	game.QuestionsAnswered = 1
	testDb.UpdateGame(game)

	req.AddCookie(&http.Cookie{Name: "sessionId", Value: game.Id.String()})

	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.NextQuestion)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if !strings.Contains(html, "Is it a Fintech or Furniture?") {
		t.Fatal(html)
	}
}

func TestLeaderboard(t *testing.T) {
	os.Remove("test.db")

	req, err := http.NewRequest("GET", "/leaderboard/?time-select=start of day", nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")

	game1, _ := testDb.CreateGame("testname1")
	game2, _ := testDb.CreateGame("testname2")

	game1.QuestionsAnswered = 10
	game2.QuestionsAnswered = 10
	game1.Score = 6
	game2.Score = 8
	game1.InProgress = false
	game2.InProgress = false
	game1.Completed = time.Now()
	game2.Completed = time.Now()

	testDb.UpdateGame(game1)
	testDb.UpdateGame(game2)

	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.Leaderboard)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if !strings.Contains(html, "Leaderboard") {
		t.Fatal(html)
	}
	if !strings.Contains(html, "testname1") {
		t.Fatal(html)
	}
	if !strings.Contains(html, "testname2") {
		t.Fatal(html)
	}
}

func TestResult(t *testing.T) {
	os.Remove("test.db")

	req, err := http.NewRequest("POST", "/result/", nil)
	if err != nil {
		t.Fatal(err)
	}

	testDb := database.InitDatabase("test.db")
	game, _ := testDb.CreateGame("testname")

	game.QuestionsAnswered = 10
	game.Score = 8
	game.InProgress = false

	testDb.UpdateGame(game)

	req.AddCookie(&http.Cookie{Name: "sessionId", Value: game.Id.String()})

	handlerContext := Context{DB: testDb}

	handler := http.HandlerFunc(handlerContext.EndPage)

	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	html := string(body)

	if !strings.Contains(html, "You achieved the score of:") {
		t.Fatal(html)
	}

	if !strings.Contains(html, "8/10") {
		t.Fatal(html)
	}
}
