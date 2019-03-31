package database

import (
	"fmt"
	"io/ioutil"

	"github.com/jackc/pgx"
)

//Application - application struct
type Application struct {
	DB *pgx.ConnPool
}

//App - export
var App Application

//OpenConnection - connects to database
func (a *Application) OpenConnection(input string) {
	pgxConfig, _ := pgx.ParseURI(input)

	a.DB, _ = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig: pgxConfig,
		})

	if query, err := ioutil.ReadFile("database/sql/items.sql"); err != nil {
		fmt.Println(err)
	} else {
		if _, err := a.DB.Exec(string(query)); err != nil {
			fmt.Println(err)
		}
	}
}

//CloseConnection - closes database connection
func (a *Application) CloseConnection() {
	a.DB.Close()
}
