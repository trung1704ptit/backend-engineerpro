package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// User represents the user data structure
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Age          int    `json:"age,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
}

// Database connection
var db *sql.DB

// Initialize the database and create the users table
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	// Create the users table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		age INTEGER,
		profile_image TEXT
	);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

// Insert a new user into the database
func CreateUser(user User) error {
	_, err := db.Exec("INSERT INTO users (username, email, password, age, profile_image) VALUES (?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Password, user.Age, user.ProfileImage)
	return err
}

// Find a user by username
func GetUserByUsername(username string) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, email, password, age, profile_image FROM users WHERE username = ?", username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.Age, &user.ProfileImage)
	return user, err
}

// Update user profile
func UpdateUser(user User) error {
	_, err := db.Exec("UPDATE users SET email = ?, age = ?, profile_image = ? WHERE username = ?",
		user.Email, user.Age, user.ProfileImage, user.Username)
	return err
}
