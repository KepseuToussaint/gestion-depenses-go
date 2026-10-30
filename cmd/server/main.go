package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"gestion-depenses/internal/handlers"
	"gestion-depenses/internal/repository"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "expenses.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := repository.Migrate(db); err != nil {
		log.Fatal(err)
	}

	repo := repository.New(db)
	h := handlers.New(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /expenses", h.ListExpenses)
	mux.HandleFunc("GET /expenses/new", h.NewExpense)
	mux.HandleFunc("POST /expenses", h.CreateExpense)
	mux.HandleFunc("GET /expenses/{id}/edit", h.EditExpense)
	mux.HandleFunc("POST /expenses/{id}/update", h.UpdateExpense)
	mux.HandleFunc("POST /expenses/{id}/delete", h.DeleteExpense)
	mux.HandleFunc("GET /about", h.About)
	mux.HandleFunc("GET /ping", h.Ping)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	addr := ":8080"
	log.Printf("Serveur démarré sur http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
