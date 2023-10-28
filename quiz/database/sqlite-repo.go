package database

import (
	"database/sql"
	"errors"
	"me885/fintech-or-furniture/quiz"

	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
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
	query := `
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
        inProgress INTEGER NOT NULL
    );
    `

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) CreateQuestion(question quiz.Question) (*quiz.Question, error) {
	res, err := r.db.Exec("INSERT INTO questions(question, answer) values(?,?)", question.Question, question.Answer)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
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

	game := quiz.Game{Id: newUuid, PlayerName: playerName, QuestionsAnswered: 0, Score: 0, InProgress: true}

	uuidBytes := game.Id
	_, err := r.db.Exec(
		"INSERT INTO games(id, playerName, questionsAnswered, score, inProgress) values(?,?,?,?,?)",
		uuidBytes,
		game.PlayerName,
		game.QuestionsAnswered,
		game.Score,
		game.InProgress)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	return &game, nil
}

func (r *SQLiteRepository) GetGameById(id uuid.UUID) (*quiz.Game, error) {
	row := r.db.QueryRow("SELECT playerName, questionsAnswered, score, inProgress FROM games WHERE id = ?", id)

	var game = quiz.Game{Id: id}
	if err := row.Scan(&game.PlayerName, &game.QuestionsAnswered, &game.Score, &game.InProgress); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	return &game, nil
}

func (r *SQLiteRepository) UpdateGame(game *quiz.Game) (*quiz.Game, error) {
	res, err := r.db.Exec(
		"UPDATE games SET playerName = ?, questionsAnswered = ?, score = ?, inProgress = ? WHERE id = ?",
		game.PlayerName,
		game.QuestionsAnswered,
		game.Score,
		game.InProgress,
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
	rows, err := r.db.Query("SELECT * FROM games")
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

func (r *SQLiteRepository) TopTenCompletedGames() ([]quiz.Game, error) {
	rows, err := r.db.Query("SELECT * FROM games WHERE inProgress=0 ORDER BY score DESC LIMIT 10")
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
