package modeldb

import (
	"time"

	"github.com/google/uuid"
)

type Result struct {
	tableName struct{} `pg:"test_results,alias:t,discard_unknown_columns"`

	ID        uuid.UUID `json:"id" db:"id"`
	TestID    uuid.UUID `json:"test_id" db:"test_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Result    string    `json:"result" db:"result"`
	Score     int       `json:"score" db:"score"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	TestName  string    `json:"test_name" db:"test_name"`
	UserName  string    `json:"user_name" db:"user_name"`
}
