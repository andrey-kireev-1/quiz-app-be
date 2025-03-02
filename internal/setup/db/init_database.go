package db

import (
	"context"
	"quiz-app-be/internal/config"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

func Init(dbConf config.Psql) (pgDB *pg.DB, err error) {
	pgDB = pg.Connect(&pg.Options{
		User:     dbConf.User,
		Password: dbConf.Pass,
		Database: dbConf.Dbname,
		Addr:     dbConf.Host + ":" + dbConf.Port,
		PoolSize: dbConf.MaxOpenConns,
	})
	if err = pgDB.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "cannot ping db")
	}
	return
}
