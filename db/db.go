package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// Creates connection
// Make maxOpenConns, maxIdleConns, maxIdleTimeMins to -1 to use default values
func OpenDB(login string, password string, host string, port int, database string, maxOpenConns int, maxIdleConns int, maxIdleTimeMins int) (*sql.DB, error) {

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", login, password, host, strconv.Itoa(port), database)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	if maxOpenConns != -1 {
		db.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConns != -1 {
		db.SetMaxIdleConns(maxIdleConns)
	}
	if maxIdleTimeMins != -1 {
		db.SetConnMaxIdleTime(time.Minute * time.Duration(maxIdleTimeMins))
	}

	return db, nil
}

func OpenDBWithConnString(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
