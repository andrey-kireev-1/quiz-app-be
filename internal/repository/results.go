package repository

import (
	modeldb "quiz-app-be/internal/model/modelDB"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type Results struct {
	db *pg.DB
}

func NewResults(db *pg.DB) *Results {
	return &Results{db: db}
}

var resultColumns = []string{
	"id",
	"test_id",
	"user_id",
	"score",
	"result",
	"created_at",
}

func (r *Results) SetResult(result modeldb.Result) error {
	_, err := r.db.Model(&result).Insert()
	return err
}

func (r *Results) GetUserResults(userID uuid.UUID) ([]modeldb.Result, error) {
	var results []modeldb.Result
	err := r.db.Model(&results).
		Column("t.*").
		ColumnExpr("tests.name AS test_name").
		Join("LEFT JOIN tests ON t.test_id = tests.id").
		Where("t.user_id = ?", userID).
		Order("t.created_at DESC").
		Select()

	return results, err
}

func (r *Results) GetResultsForAuthorTests(authorID uuid.UUID) ([]modeldb.Result, error) {
	var results []modeldb.Result
	err := r.db.Model(&results).
		Column("t.*").
		ColumnExpr("tests.name AS test_name").
		ColumnExpr("users.name AS user_name").
		Join("JOIN tests ON tests.id = t.test_id").
		Join("JOIN users ON users.id = t.user_id").
		Where("tests.author_id = ?", authorID).
		Order("t.created_at DESC").
		Select()

	return results, err
}
