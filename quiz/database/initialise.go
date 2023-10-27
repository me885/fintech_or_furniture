package database

import (
	"database/sql"
	"log"
	"me885/fintech-or-furniture/quiz"
)

func InitDatabase(filename string) *SQLiteRepository {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}

	sqliteRepository := NewSQLiteRepository(db)

	if err := sqliteRepository.Migrate(); err != nil {
		log.Fatal(err)
	}

	questions := [...]quiz.Question{
		{Question: "PAX", Answer: quiz.Furniture},
		{Question: "YAVRIO", Answer: quiz.Fintech},
		{Question: "YPPERLIG", Answer: quiz.Furniture},
		{Question: "ZYNGA", Answer: quiz.Fintech},
		{Question: "SLYP", Answer: quiz.Fintech},
		{Question: "FADO", Answer: quiz.Furniture},
		{Question: "LACK", Answer: quiz.Furniture},
		{Question: "TROFAST", Answer: quiz.Furniture},
		{Question: "ANROK", Answer: quiz.Fintech},
		{Question: "VOXNAN", Answer: quiz.Furniture},
		{Question: "VOWCH", Answer: quiz.Fintech},
		{Question: "CRUX", Answer: quiz.Fintech},
		{Question: "FYSSE", Answer: quiz.Furniture},
		{Question: "STORI", Answer: quiz.Fintech},
		{Question: "KALLAX", Answer: quiz.Furniture},
		{Question: "PAGOS", Answer: quiz.Fintech},
	}

	for _, element := range questions {
		sqliteRepository.CreateQuestion(element)
	}

	return sqliteRepository
}
