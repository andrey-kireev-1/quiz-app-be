package modeldb

import (
	"time"

	"github.com/google/uuid"
)

type Test struct {
	tableName struct{} `pg:"tests,alias:t,discard_unknown_columns"`

	ID            uuid.UUID `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Description   *string   `json:"description,omitempty" db:"description"`
	Data          string    `json:"data" db:"data"`
	ImageURL      *string   `json:"image_url,omitempty" db:"image_url"`
	IsStrict      bool      `json:"is_strict,omitempty" db:"is_strict"`
	IsPrivate     bool      `json:"is_private,omitempty" db:"is_private"`
	AuthorID      uuid.UUID `json:"author_id" db:"author_id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	AttemptsCount int       `json:"attempts_count" pg:"attempts_count"`
	AuthorName    string    `json:"author_name" pg:"author_name"`
}
