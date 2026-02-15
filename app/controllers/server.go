package controllers

import (
	"net/http"

	"liveconnect/config"
)

func StartMainServer() error {

	files := http.FileServer(http.Dir(config.Config.Web.Static))
	http.Handle("/static/", http.StripPrefix("/static/", files))

	http.HandleFunc("/", Login)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/register", RegisterPage)
	http.HandleFunc("/register/submit", RegisterSubmit)
	http.HandleFunc("/lives", LivesIndex)
	http.HandleFunc("/lives/new", LivesNew)
	http.HandleFunc("/lives/create", LivesCreate)
	http.HandleFunc("/lives/delete", LivesDelete)
	http.HandleFunc("/lives/edit", LivesEdit)
	http.HandleFunc("/lives/update", LivesUpdate)
	http.HandleFunc("/users", UsersIndex)
	http.HandleFunc("/users/lives", UserLives)

	return http.ListenAndServe(":"+config.Config.Web.Port, nil)
}
