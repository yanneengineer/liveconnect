package controllers

import (
	"html/template"
	"net/http"

	"liveconnect/app/models"
)

func Home(w http.ResponseWriter, r *http.Request) {

	session, _ := Store.Get(r, "session")

	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	basePath := "app/views/templates"

	tmpl := template.Must(template.ParseFiles(
		basePath+"/layout/base.html",
		basePath+"/home/index.html",
	))

	data := map[string]interface{}{
		"Title":    "ホーム",
		"UserName": user.Name,
	}

	tmpl.ExecuteTemplate(w, "base", data)
}
