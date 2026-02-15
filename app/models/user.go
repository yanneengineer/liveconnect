package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type User struct {
	ID       string
	Name     string
	Password string
}

func GetAllUsers(id string) ([]User, error) {
	rows, err := DB.Query(`
		SELECT id, name
		FROM users
		WHERE NOT id = ?
		ORDER BY id
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

func AuthenticateUser(id, password string) (*User, error) {
	var user User

	err := DB.QueryRow(`
		SELECT id, name, password
		FROM users
		WHERE id = ?
	`, id).Scan(&user.ID, &user.Name, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return &user, nil
}

func CreateUser(id, name, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO users (id, name, password)
	VALUES (?, ?, ?)
	`
	_, err = DB.Exec(query, id, name, string(hashedPassword))

	return err
}

func GetUserByID(id string) (User, error) {
	var user User
	err := DB.QueryRow(`
		SELECT id, name, password
		FROM users
		WHERE id = ?
	`, id).Scan(&user.ID, &user.Name, &user.Password)
	return user, err
}

func ExistsUserID(id string) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM users WHERE id = ?
	`, id).Scan(&count)
	return count > 0, err
}

func ExistsUserName(name string) (bool, error) {
	var count int
	err := DB.QueryRow(`
		SELECT COUNT(*) FROM users WHERE name = ?
	`, name).Scan(&count)
	return count > 0, err
}
