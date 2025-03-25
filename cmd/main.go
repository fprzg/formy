package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type App struct {
	db *sql.DB
}

var (
	app = App{}
)

func main() {
	db, err := sql.Open("sqlite", "./forms.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app = App{
		db: db,
	}

	http.HandleFunc("/form/create", formCreateHandler())
	http.HandleFunc("/submit", submitHandler())
}

func formCreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var form interface{}
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		form.ID = uuid.New().String()
		form.CreatedAt = time.Now().Unix()
		form.Enabled = true

		_, err := db.Exec("INSERT INTO forms (id, name, enabled, created_at) VALUES (?, ?, ?, ?)",
			form.ID, form.Name, form.Enabled, form.CreatedAt)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(form)
	}
}

func submitHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var submission interface{}
		if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		dataJSON, _ := json.Marshall(submission.Data)
		_, err := db.Exec("INSERT INTO submissions (form_id, submitted_at, data) VALUES (?, ?, ?)",
			submission.FormID, time.Now().Unix(), dataJSON)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// anti-abuse checks
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
