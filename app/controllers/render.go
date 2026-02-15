package controllers

import (
	"html/template"
	"log"
	"net/http"

	"liveconnect/config"
)

func render(w http.ResponseWriter, page string, data interface{}) {

	basePath := config.Config.Web.Views + "/templates"

	t, err := template.ParseFiles(
		basePath+"/layout/base.html",
		basePath+"/"+page,
	)
	if err != nil {
		log.Println("template parse error:", err)
		http.Error(w, "Template Parse Error", http.StatusInternalServerError)
		return
	}

	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println("template execute error:", err)
	}
}
