package comradetest

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func CreateDbMock() (*sql.DB, error) {
	db, _, err := sqlmock.New()

	return db, err

}
