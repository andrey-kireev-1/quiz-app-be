package repository

import (
	modeldb "quiz-app-be/internal/model/modelDB"
	modelhttp "quiz-app-be/internal/model/modelHttp"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
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
	_, err := s.db.Model(&test).
		Column(testColumns...).
		Insert()
	return err
}

func (s *Tests) GetTestByID(id string) (modeldb.Test, error) {
	test := modeldb.Test{}
	var prefixedColumns []string
	for _, v := range testColumns {
		prefixedColumns = append(prefixedColumns, "tests."+v)
	}
	err := s.db.Model(&test).
		Table("tests").
		Column(prefixedColumns...).
		ColumnExpr("u.name AS author_name").
		Join("LEFT JOIN users AS u ON u.id = tests.author_id").
		Where("tests.id = ?", id).
		Group("tests.id", "u.name").
		Select()
	return test, err
}

func (s *Tests) GetHomeTests(limit, offset int) ([]modeldb.Test, error) {
	var tests []modeldb.Test
	var prefixedColumns []string
	for _, v := range testColumns {
		prefixedColumns = append(prefixedColumns, "tests."+v)
	}
	err := s.db.Model(&tests).
		Table("tests").
		Column(prefixedColumns...).
		ColumnExpr("count(tr.id) as attempts_count").
		ColumnExpr("u.name AS author_name").
		Join("LEFT JOIN test_results AS tr ON tr.test_id = tests.id").
		Join("LEFT JOIN users AS u ON u.id = tests.author_id").
		Where("tests.is_private = ?", false).
		Group("tests.id", "u.name").
		Order("attempts_count DESC", "tests.created_at DESC").
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

func (s *Tests) GetUserTests(authorID uuid.UUID) ([]modeldb.Test, error) {
	var tests []modeldb.Test
	err := s.db.Model(&tests).
		Column(testColumns...).
		Where("author_id = ?", authorID).
		Order("created_at DESC").
		Select()

	return tests, err
}

func (s *Tests) GetFilteredTests(filters modelhttp.TestFilters) ([]modeldb.Test, error) {
	var tests []modeldb.Test
	var prefixedColumns []string
	for _, v := range testColumns {
		prefixedColumns = append(prefixedColumns, "tests."+v)
	}
	query := s.db.Model(&tests).
		Table("tests").
		Column(prefixedColumns...).
		ColumnExpr("count(tr.id) as attempts_count").
		ColumnExpr("u.name AS author_name").
		Join("LEFT JOIN test_results AS tr ON tr.test_id = tests.id").
		Join("LEFT JOIN users AS u ON u.id = tests.author_id").
		Where("tests.is_private = ?", false).
		Group("tests.id", "u.name")

	if filters.TestName != nil {
		query = query.Where("tests.name ILIKE ?", "%"+*filters.TestName+"%")
	}

	if filters.AuthorName != nil {
		query = query.Where("u.name ILIKE ?", "%"+*filters.AuthorName+"%")
	}

	if filters.CreatedFrom != nil {
		query = query.Where("tests.created_at >= ?", *filters.CreatedFrom)
	}

	if filters.CreatedTo != nil {
		query = query.Where("tests.created_at <= ?", *filters.CreatedTo)
	}

	if filters.IsStrict != nil {
		query = query.Where("tests.is_strict = ?", *filters.IsStrict)
	}

	err := query.Order("attempts_count DESC", "tests.created_at DESC").Select()
	return tests, err
}
