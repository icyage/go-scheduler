package utils

import (
	"database/sql"
	"github.com/cenkalti/backoff"
	_ "github.com/lib/pq"
	"github.com/obsurvive/voyager/log"
)

var (
	db *sql.DB
)

func connect() (error, *sql.DB) {
	log.Debug("about to open pq connection")
	dbi, err := sql.Open("postgres", cfg.GetString("postgres_dsn"))
	if err != nil {
		log.Debugf("error opening connection", err)
		return err, nil
	} else {
		log.Debug("pq open connection success")
		return nil, dbi
	}
}

func GetDB() (error, *sql.DB) {
	if db == nil {
		log.Debug("about to start connecting to db with backoff")
		err := backoff.Retry(func() error {
			var retErr error
			retErr, db = connect()
			log.Debugf("got db", db)
			return retErr
		}, bo)
		if err != nil {
			log.Debugf("error with backoff", err)
			return err, nil
		}
		log.Debugf("returning db", db)
		return nil, db
	} else {
		log.Debug("returning db, was already connected")
		return nil, db
	}
}
