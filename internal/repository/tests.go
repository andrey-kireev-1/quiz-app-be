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

var testColumns = []string{
	"id",
	"name",
	"description",
	"data",
	"image_url",
	"is_strict",
	"is_private",
	"author_id",
	"created_at",
	"updated_at",
}

var testExtraColumns = []string{
	"attempts_count",
}

func (s *Tests) CreateTest(test modeldb.Test) error {
	_, err := s.db.Model(&test).Insert()
	return err
}

func (s *Tests) GetTestByID(id string) (modeldb.Test, error) {
	test := modeldb.Test{}
	err := s.db.Model(&test).Column(testColumns...).Where("id = ?", id).Select()
	return test, err
}

func (s *Tests) GetAllTests(limit, offset int) ([]modeldb.Test, error) {
	var tests []modeldb.Test
	err := s.db.Model(&tests).TableExpr("tests AS test").
		Column("test.*").
		ColumnExpr("count(tr.id) as attempts_count").
		Join("LEFT JOIN test_results AS tr ON tr.test_id = test.id").
		Where("test.is_private = ?", false).
		Group("test.id").
		Order("attempts_count DESC", "test.created_at DESC").
		Limit(limit).
		Offset(offset).
		Select()
	return tests, err
}

func (s *Tests) CountAllPublicTests() (int, error) {
	var tests []modeldb.Test
	cnt, err := s.db.Model(&tests).
		Column(testColumns...).
		Where("is_private = ?", false).
		Count()
	return cnt, err
}
