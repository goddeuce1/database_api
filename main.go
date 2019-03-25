package main

import (
	"fmt"
	"log"

	"./database"
	"./router"

	"github.com/valyala/fasthttp"
)

const conn = "postgres://gd1:123@localhost:5432/api_db"
const port = ":5000"

func main() {
	database.App.OpenConnection(conn)
	defer database.App.CloseConnection()

	fmt.Println("Starting server...")
	log.Fatal(fasthttp.ListenAndServe(port, router.Router.Handler))
}
