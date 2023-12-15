package database

import (
	"database/sql"
	"errors"
	"me885/fintech-or-furniture/quiz"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `--sql

    CREATE TABLE IF NOT EXISTS questions(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        question TEXT NOT NULL UNIQUE,
        answer INTEGER NOT NULL
    );

	CREATE TABLE IF NOT EXISTS games(
        id BLOB PRIMARY KEY,
        playerName TEXT NOT NULL,
        questionsAnswered INTEGER NOT NULL,
        score INTEGER NOT NULL,
        inProgress INTEGER NOT NULL,
		created BLOB,
		completed BLOB
    );

	CREATE TABLE IF NOT EXISTS gameQuestions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		gameId BLOB NOT NULL,
		QuestionId INTEGER NOT NULL
	);
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) AddGameQuestion(gameId uuid.UUID, questionId int64) error {
	_, err := r.db.Exec("INSERT INTO gameQuestions(id, gameId, questionId) values(NULL,?,?)", gameId, questionId)

	return err
}

func (r *SQLiteRepository) RemoveGameQuestions(gameId uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM gameQuestions WHERE gameId = ?", gameId)

	return err
}

func (r *SQLiteRepository) GetUnansweredQuestions(gameId uuid.UUID) ([]quiz.Question, error) {
	rows, err := r.db.Query(`--sql
		SELECT id, question, answer
		FROM questions
		WHERE id NOT IN (
			SELECT questionId
			FROM gameQuestions
			WHERE gameId = ?
		)`,
		gameId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var all []quiz.Question
	for rows.Next() {
		var question quiz.Question
		if err := rows.Scan(&question.Id, &question.Question, &question.Answer); err != nil {
			return nil, err
		}
		all = append(all, question)
	}
	return all, nil
}

func (r *SQLiteRepository) CreateQuestion(question quiz.Question) (*quiz.Question, error) {
	res, err := r.db.Exec("INSERT INTO questions(question, answer) values(?,?)", question.Question, question.Answer)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	question.Id = id

	return &question, nil
}

func (r *SQLiteRepository) GetQuestionById(id int64) (*quiz.Question, error) {
	row := r.db.QueryRow("SELECT * FROM questions WHERE id = ?", id)

	var question quiz.Question
	if err := row.Scan(&question.Id, &question.Question, &question.Answer); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &question, nil
}

func (r *SQLiteRepository) CountQuestions() int64 {
	row := r.db.QueryRow("SELECT COUNT(*) FROM Products;")

	var count int64
	row.Scan(&count)
	return count
}

func (r *SQLiteRepository) CreateGame(playerName string) (*quiz.Game, error) {
	newUuid, _ := uuid.NewUUID()

	game := quiz.Game{Id: newUuid, PlayerName: playerName, QuestionsAnswered: 0, Score: 0, InProgress: true, Created: time.Now()}

	uuidBytes := game.Id
	_, err := r.db.Exec(
		"INSERT INTO games(id, playerName, questionsAnswered, score, inProgress, created) values(?,?,?,?,?,?)",
		uuidBytes,
		game.PlayerName,
		game.QuestionsAnswered,
		game.Score,
		game.InProgress,
		game.Created)

	if err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *SQLiteRepository) GetGameById(id uuid.UUID) (*quiz.Game, error) {
	row := r.db.QueryRow("SELECT playerName, questionsAnswered, score, inProgress, created, completed FROM games WHERE id = ?", id)

	var createdStr *string
	var completedStr *string

	var game = quiz.Game{Id: id}
	if err := row.Scan(&game.PlayerName, &game.QuestionsAnswered, &game.Score, &game.InProgress, &createdStr, &completedStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	if createdStr != nil {
		game.Created, _ = time.Parse("2006-01-02 15:04:05+00:00", *createdStr)
	}
	if completedStr != nil {
		game.Completed, _ = time.Parse("2006-01-02 15:04:05+00:00", *completedStr)
	}

	return &game, nil
}

func (r *SQLiteRepository) UpdateGame(game *quiz.Game) (*quiz.Game, error) {
	res, err := r.db.Exec(
		"UPDATE games SET playerName = ?, questionsAnswered = ?, score = ?, inProgress = ?, created = ?, completed = ? WHERE id = ?",
		game.PlayerName,
		game.QuestionsAnswered,
		game.Score,
		game.InProgress,
		game.Created,
		game.Completed,
		game.Id)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return game, nil
}

func (r *SQLiteRepository) AllGames() ([]quiz.Game, error) {
	rows, err := r.db.Query("SELECT id, playerName, questionsAnswered, score, inProgress FROM games")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []quiz.Game
	for rows.Next() {
		var game quiz.Game
		if err := rows.Scan(&game.Id, &game.PlayerName, &game.QuestionsAnswered, &game.Score, &game.InProgress); err != nil {
			return nil, err
		}
		all = append(all, game)
	}
	return all, nil
}

func (r *SQLiteRepository) TopTenCompletedGames(sinceTime string) ([]quiz.Game, error) {
	rows, err := r.db.Query(`--sql
	SELECT id, playerName, questionsAnswered, score, inProgress 
	FROM games 
	WHERE inProgress=0 AND completed > date('now', ?)
	ORDER BY score 
	DESC LIMIT 10
	`, sinceTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []quiz.Game
	for rows.Next() {
		var game quiz.Game
		if err := rows.Scan(&game.Id, &game.PlayerName, &game.QuestionsAnswered, &game.Score, &game.InProgress); err != nil {
			return nil, err
		}
		all = append(all, game)
	}
	return all, nil
}
