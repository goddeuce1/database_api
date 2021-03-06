package main

import (
	"log"
	"park_base/park_db/database"
	"park_base/park_db/router"

	"github.com/valyala/fasthttp"
)

const conn = "postgres://gd1:123@localhost:5432/api_db"
const port = ":5000"

func main() {
	database.App.OpenConnection(conn)
	defer database.App.CloseConnection()

	log.Fatal(fasthttp.ListenAndServe(port, router.Router.Handler))
}
