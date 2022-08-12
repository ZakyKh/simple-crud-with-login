package main

import (
	"mzaky/simple-crud-with-login/tasks"

	"fmt"
	"net/http"
	"log"
	"os"
	"database/sql"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func HandleHealthz(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "OK\n")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env: " + err.Error())
	}

	db, err := connectPgSql()
	if err != nil {
		log.Fatal("Failed to connect to database: " + err.Error())
		return
	}

	router := httprouter.New()
	router.GET("/healthz", HandleHealthz)

	taskHandler := tasks.NewHandler(db)
	router.GET("/tasks", taskHandler.GetTasks)
	router.GET("/tasks/:id", taskHandler.GetTask)
	router.POST("/tasks", taskHandler.CreateTask)

	fmt.Printf("%s is now listening on port %s\n", os.Getenv("APP_NAME"), os.Getenv("APP_PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), router))
}

func connectPgSql() (*sql.DB, error) {
	pgUser := os.Getenv("POSTGRESQL_USERNAME")
	pgPass := os.Getenv("POSTGRESQL_PASSWORD")
	pgHost := os.Getenv("POSTGRESQL_HOST")
	pgPort := os.Getenv("POSTGRESQL_PORT")
	pgDb := os.Getenv("POSTGRESQL_DATABASE")
	pgSSL := os.Getenv("POSTGRESQL_SSL")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", pgUser, pgPass, pgHost, pgPort, pgDb, pgSSL)
	return sql.Open("postgres", connStr)
}