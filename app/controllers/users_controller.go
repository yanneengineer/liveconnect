package controllers

import (
	"net/http"
	"regexp"

	"liveconnect/app/models"
)

type RegisterViewData struct {
	Title  string
	Error  string
	UserID string
	Name   string
}

type UsersIndexViewData struct {
	Title string
	Users []models.User
}

func UserLives(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	lives, err := models.GetLivesByUser(userID)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	render(w, "users/lives.html", lives)
}

func UsersIndex(w http.ResponseWriter, r *http.Request) {

	session, _ := Store.Get(r, "session")
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	users, err := models.GetAllUsers(userID)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	data := UsersIndexViewData{
		Title: "ユーザー一覧",
		Users: users,
	}

	render(w, "users/index.html", data)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	render(w, "users/register.html", RegisterViewData{
		Title: "ユーザー登録",
	})
}

func RegisterSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	name := r.FormValue("name")
	password := r.FormValue("password")

	data := RegisterViewData{
		Title:  "ユーザー登録",
		UserID: id,
		Name:   name,
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{1,8}$`, id)
	if !matched {
		data.Error = "ユーザーIDは英数字8文字以内で入力してください"
		renderRegister(w, data)
		return
	}

	exists, err := models.ExistsUserID(id)
	if err != nil {
		data.Error = "登録に失敗しました"
		renderRegister(w, data)
		return
	}
	if exists {
		data.Error = "このユーザーIDはすでに使われています"
		renderRegister(w, data)
		return
	}

	exists, err = models.ExistsUserName(name)
	if err != nil {
		data.Error = "登録に失敗しました"
		renderRegister(w, data)
		return
	}
	if exists {
		data.Error = "このユーザー名はすでに使われています"
		renderRegister(w, data)
		return
	}

	err = models.CreateUser(id, name, password)
	if err != nil {
		data.Error = "登録に失敗しました"
		renderRegister(w, data)
		return
	}

	user, err := models.GetUserByID(id)
	if err != nil {
		data.Error = "ログイン処理に失敗しました"
		renderRegister(w, data)
		return
	}

	session, _ := Store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Save(r, w)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func renderRegister(w http.ResponseWriter, data RegisterViewData) {

	render(w, "users/register.html", data)
}
