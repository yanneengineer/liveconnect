package controllers

import (
	"fmt"
	"net/http"
	"time"

	"liveconnect/app/models"
)

type LivesIndexViewData struct {
	Title  string
	Future []models.Live
	Past   []models.Live
	Rank   []models.ArtistCount
}

func LivesIndex(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session")
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	future, past, err := models.GetLivesGrouped(userID)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	artistCounts, err := models.GetArtistCounts(userID)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	data := LivesIndexViewData{
		Title:  "ライブ一覧",
		Future: future,
		Past:   past,
		Rank:   artistCounts,
	}

	render(w, "lives/index.html", data)
}

func LivesNew(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "ライブ登録",
	}

	render(w, "lives/new.html", data)
}

func LivesCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/lives", http.StatusSeeOther)
		return
	}

	session, _ := Store.Get(r, "session")
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	artist := r.FormValue("artist")
	dateStr := r.FormValue("date")
	venue := r.FormValue("venue")

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", 400)
		return
	}

	live := models.Live{
		UserID: userID,
		Artist: artist,
		Date:   parsedDate,
		Venue:  venue,
	}

	err = models.CreateLive(live)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	http.Redirect(w, r, "/lives", http.StatusSeeOther)
}

func LivesDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/lives", http.StatusSeeOther)
		return
	}

	session, _ := Store.Get(r, "session")
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	liveID := r.FormValue("live_id")
	if liveID == "" {
		http.Redirect(w, r, "/lives", http.StatusSeeOther)
		return
	}

	err := models.DeleteLiveByID(liveID, userID)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	http.Redirect(w, r, "/lives", http.StatusSeeOther)
}

func LivesEdit(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "session")
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	liveID := r.URL.Query().Get("id")
	if liveID == "" {
		http.Redirect(w, r, "/lives", http.StatusSeeOther)
		return
	}

	live, err := models.GetLiveByID(liveID, userID)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	data := struct {
		Title string
		models.Live
	}{
		Title: "ライブ編集",
		Live:  live,
	}

	render(w, "lives/edit.html", data)
}

func LivesUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/lives", http.StatusSeeOther)
		return
	}

	session, _ := Store.Get(r, "session")
	userID, ok := session.Values["user_id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	artist := r.FormValue("artist")
	dateStr := r.FormValue("date")
	venue := r.FormValue("venue")

	var liveID int
	fmt.Sscan(id, &liveID)

	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format", 400)
		return
	}

	live := models.Live{
		ID:     liveID,
		UserID: userID,
		Artist: artist,
		Date:   parsedDate,
		Venue:  venue,
	}

	err = models.UpdateLive(live)
	if err != nil {
		http.Error(w, "DB Error", 500)
		return
	}

	http.Redirect(w, r, "/lives", http.StatusSeeOther)
}
