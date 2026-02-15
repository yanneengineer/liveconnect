package controllers

import (
	"errors"
	"html/template"
	"liveconnect/app/models"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmpl := template.Must(template.ParseFiles(
			"app/views/templates/layout/base.html",
			"app/views/templates/login.html",
		))

		data := map[string]interface{}{
			"Title": "ログイン",
		}

		tmpl.ExecuteTemplate(w, "base", data)
		return
	}

	userID := r.FormValue("id")
	password := r.FormValue("password")

	data := map[string]interface{}{
		"Title": "ログイン",
		"ID":    userID,
	}

	user, err := models.AuthenticateUser(userID, password)
	if err != nil {

		if errors.Is(err, models.ErrInvalidCredentials) {
			data["Error"] = "IDまたはパスワードが違います"
		} else {
			data["Error"] = "システムエラーが発生しました"
		}

		tmpl := template.Must(template.ParseFiles(
			"app/views/templates/layout/base.html",
			"app/views/templates/login.html",
		))

		tmpl.ExecuteTemplate(w, "base", data)
		return
	}

	session, _ := Store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	session, _ := Store.Get(r, "session")

	delete(session.Values, "user_id")

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
