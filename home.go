package main

import (
	//"encoding/json"
	"html/template"
	"log"
	"net/http"
	_ "github.com/lib/pq"
)


func HomePage(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("/home/sergey/site2/templates/home_page.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// Отображаем страницу
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Println(err)
	}
}