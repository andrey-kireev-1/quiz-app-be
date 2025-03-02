package modeldb

import (
	"time"

	"github.com/google/uuid"
)

// User - структура для хранения информации о пользователе`
type User struct {
	tableName struct{} `pg:"users,alias:t,discard_unknown_columns"`

	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
