package models

import (
	"time"
)

type Live struct {
	ID     int
	UserID string
	Name   string
	Artist string
	Date   string
	Venue  string
}

type ArtistCount struct {
	Artist string
	Count  int
}

func CreateLive(live Live) error {
	query := `
	INSERT INTO lives (user_id, artist, date, venue)
	VALUES (?, ?, ?, ?)
	`

	_, err := DB.Exec(
		query,
		live.UserID,
		live.Artist,
		live.Date,
		live.Venue,
	)

	return err
}

func GetArtistCounts(userID string) ([]ArtistCount, error) {

	today := time.Now().Format("2006-01-02")

	rows, err := DB.Query(`
		SELECT artist, COUNT(*) as count
		FROM lives
		WHERE user_id = ?
		AND date < ?
		GROUP BY artist
		ORDER BY count DESC
	`, userID, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ArtistCount

	for rows.Next() {
		var ac ArtistCount
		err := rows.Scan(&ac.Artist, &ac.Count)
		if err != nil {
			return nil, err
		}
		results = append(results, ac)
	}

	return results, nil
}

func GetLivesGrouped(userID string) ([]Live, []Live, error) {
	rows, err := DB.Query(`
		SELECT id, user_id, artist, date, venue
		FROM lives
		WHERE user_id = ?
		ORDER BY date
	`, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var future []Live
	var past []Live

	today := time.Now().Truncate(24 * time.Hour)

	for rows.Next() {
		var l Live
		err := rows.Scan(
			&l.ID,
			&l.UserID,
			&l.Artist,
			&l.Date,
			&l.Venue,
		)
		if err != nil {
			return nil, nil, err
		}

		liveDate, err := time.Parse("2006-01-02", l.Date)
		if err != nil {
			return nil, nil, err
		}

		if liveDate.Before(today) {
			past = append(past, l)
		} else {
			future = append(future, l)
		}
	}

	return future, past, nil
}

func GetLivesByUser(userID string) (map[string]interface{}, error) {

	var userName string
	err := DB.QueryRow(
		"SELECT name FROM users WHERE id = ?",
		userID,
	).Scan(&userName)
	if err != nil {
		return nil, err
	}

	rows, err := DB.Query(`
        SELECT l.id, l.user_id, u.name, l.artist, l.date, l.venue
        FROM lives l
        INNER JOIN users u ON l.user_id = u.id
        WHERE l.user_id = ?
        ORDER BY l.date
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var futureLives []Live
	var pastLives []Live

	today := time.Now().Truncate(24 * time.Hour)

	for rows.Next() {
		var live Live
		err := rows.Scan(
			&live.ID,
			&live.UserID,
			&live.Name,
			&live.Artist,
			&live.Date,
			&live.Venue,
		)
		if err != nil {
			return nil, err
		}

		liveDate, err := time.Parse("2006-01-02", live.Date)
		if err != nil {
			return nil, err
		}

		if liveDate.Before(today) {
			pastLives = append(pastLives, live)
		} else {
			futureLives = append(futureLives, live)
		}
	}

	return map[string]interface{}{
		"Name":   userName,
		"Future": futureLives,
		"Past":   pastLives,
	}, nil
}

func DeleteLiveByID(liveID string, userID string) error {
	query := `
		DELETE FROM lives
		WHERE id = ? AND user_id = ?
	`
	_, err := DB.Exec(query, liveID, userID)
	return err
}

func GetLiveByID(liveID string, userID string) (Live, error) {
	var live Live

	query := `
		SELECT id, user_id, artist, date, venue
		FROM lives
		WHERE id = ? AND user_id = ?
	`

	err := DB.QueryRow(query, liveID, userID).Scan(
		&live.ID,
		&live.UserID,
		&live.Artist,
		&live.Date,
		&live.Venue,
	)

	return live, err
}

func UpdateLive(live Live) error {
	query := `
		UPDATE lives
		SET artist = ?, date = ?, venue = ?
		WHERE id = ? AND user_id = ?
	`

	_, err := DB.Exec(
		query,
		live.Artist,
		live.Date,
		live.Venue,
		live.ID,
		live.UserID,
	)

	return err
}
