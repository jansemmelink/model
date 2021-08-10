package system

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/pkg/errors"
)

var (
	dbMutex sync.Mutex
	db      *sql.DB
)

func dbconn() (*sql.DB, error) {
	if db == nil {
		dbMutex.Lock()
		defer dbMutex.Unlock()
		if db == nil {
			var err error
			db, err = sql.Open("mysql", "root:asdf@tcp(127.0.0.1:3306)/stock2shop")
			if err != nil {
				return nil, errors.Wrapf(err, "failed to connect to db")
			}
		}
	}
	return db, nil
}
