package repository

import (
	modeldb "quiz-app-be/internal/model/modelDB"

	"github.com/go-pg/pg/v10"
)

type Tests struct {
	db *pg.DB
}

func NewTests(db *pg.DB) *Tests {
	return &Tests{db: db}
}

func (s *Tests) CreateTest(test modeldb.Test) error {
	_, err := s.db.Model(&test).Insert()
	return err
}

func (s *Tests) GetTestByID(id string) (modeldb.Test, error) {
	test := modeldb.Test{}
	err := s.db.Model(&test).Where("id = ?", id).Select()
	return test, err
}
