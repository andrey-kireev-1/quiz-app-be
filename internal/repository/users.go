package repository

import (
	modelDB "quiz-app-be/internal/model/modelDB"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type Users struct {
	db *pg.DB
}

func NewUsers(db *pg.DB) *Users {
	return &Users{db: db}
}

func (s *Users) CreateUser(user modelDB.User) error {
	_, err := s.db.Model(&user).Insert()
	return err
}

func (s *Users) GetUserByEmail(email string) (modelDB.User, error) {
	user := modelDB.User{}
	err := s.db.Model(&user).Where("email = ?", email).Select()
	return user, err
}

func (s *Users) CheckUserByID(id uuid.UUID) (bool, error) {
	user := modelDB.User{}
	return s.db.Model(&user).Where("id = ?", id).Exists()
}
