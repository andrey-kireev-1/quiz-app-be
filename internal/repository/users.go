package repository

import (
	modeldb "quiz-app-be/internal/model/modelDB"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type Users struct {
	db *pg.DB
}

func NewUsers(db *pg.DB) *Users {
	return &Users{db: db}
}

func (s *Users) CreateUser(user modeldb.User) (modeldb.User, error) {
	_, err := s.db.Model(&user).Insert()
	return user, err
}

func (s *Users) GetUserByEmail(email string) (modeldb.User, error) {
	user := modeldb.User{}
	err := s.db.Model(&user).Where("email = ?", email).Select()
	return user, err
}

func (s *Users) CheckUserByID(id uuid.UUID) (bool, error) {
	user := modeldb.User{}
	return s.db.Model(&user).Where("id = ?", id).Exists()
}

func (r *Users) GetUserProfile(userID uuid.UUID) (modeldb.User, error) {
	var user modeldb.User
	err := r.db.Model(&user).
		Column("name", "email").
		Where("id = ?", userID).
		Select()
	return user, err
}
