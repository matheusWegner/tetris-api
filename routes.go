package main

import (
	"database/sql"
	"net/http"
)

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Conectando ao banco de dados
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/tetris")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//users, err := getUsers(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//json, err := json.Marshal(users)
}
