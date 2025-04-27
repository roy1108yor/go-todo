package models

import (
	"database/sql"
	"errors"
	"github.com/ichtrojan/go-todo/config"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}

// CreateUser creates a new user in the database with encrypted password
func CreateUser(user User) (int64, error) {
	database := config.Database()

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Insert the new user
	result, err := database.Exec(
		`INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)`,
		user.Username,
		user.Email,
		hashedPassword,
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	// Return the ID of the newly created user
	id, err := result.LastInsertId()
	return id, err
}

// FindUserByEmail finds a user by their email address
func FindUserByEmail(email string) (User, error) {
	database := config.Database()
	var user User

	err := database.QueryRow(
		`SELECT id, username, email, password, created_at FROM users WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}

// CheckPassword compares the provided password with the stored hashed password
func (user User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id int) (User, error) {
	database := config.Database()
	var user User

	err := database.QueryRow(
		`SELECT id, username, email, password, created_at FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}

	return user, nil
}

// InitUserTable ensures the users table exists in the database
func InitUserTable() {
	database := config.Database()

	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT,
			username VARCHAR(50) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		);
	`)

	if err != nil {
		// Log error but don't stop execution
		// The table might already exist
		// We don't want to halt the application for this
		// just keep going
	}
}