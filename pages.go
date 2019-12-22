package main

import (
	"fmt"
	"html/template"
	"net/http"
)

//Game handler
func game(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("new_map.html")
	if err != nil {
		panic(err)
	}
	if winner != "" {
		http.Redirect(w, r, "/winner", http.StatusFound)
	}
	id, err := r.Cookie("id")
	if err != nil {
		http.Redirect(w, r, "/wrong_login", http.StatusFound)
		return
	}
	sendPots(id.Value)
	t.Execute(w, nil)
}

//Incorrect login
func badLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Invalid login code")
}

//Win page
func win(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, winner+" wins")
}

//Lose page
func lose(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "YOU LOSE")
}
