package dbrepo

import (
	"database/sql"

	"github.com/xtophe02/bookings-go/internal/config"
	"github.com/xtophe02/bookings-go/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

//CREATE A DATABASE REPOSITORY SO THE APP HAS ACCCESS ANYWHERE
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{App: a, DB: conn}
}
func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{App: a}
}

