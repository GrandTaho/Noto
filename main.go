package main

import (
	"log"
	"net/http"

	"github.com/GrandTaho/noto/database"
	"github.com/GrandTaho/noto/note"
	_ "github.com/go-sql-driver/mysql"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	note.SetupRoutes(apiBasePath)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
