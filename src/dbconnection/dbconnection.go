package dbconnection

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

type DatabadeConnection struct{}

//go:generate mockgen --source=dbconnection.go --package=mocks --destination=./mocks/dbconnection.go  DbConnectionService
type DbConnectionService interface {
	GetDb() (*sql.DB, error)
}

func NewDatabadeConnection() *DatabadeConnection {
	return &DatabadeConnection{}
}

func (dbc *DatabadeConnection) GetDb() (*sql.DB, error) {
	source := "user=" + os.Getenv("POSTGRES_USER") +
		" dbname=" + os.Getenv("POSTGRES_NAME") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" host=" + os.Getenv("POSTGRES_HOST") +
		" sslmode=d" + os.Getenv("POSTGRES_SSLMODE")

	return sql.Open("postgres", source)
}
