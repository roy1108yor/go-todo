package models

import (
	"database/sql"
	"errors"
	"github.com/ichtrojan/go-todo/config"
	"time"
)

type Todo struct {
	Id        int
	Item      string
	Completed int
	CreatedAt time.Time
	UserID    int // Added UserID field to associate todos with users
}

// FormatCreatedAt formats and returns the creation time as a string
func (todo Todo) FormatCreatedAt() string {
	return todo.CreatedAt.Format("2006-01-02 15:04:05")
}

// CreateTodo creates a new todo item in the database
func CreateTodo(todo Todo) error {
	database := config.Database()

	_, err := database.Exec(
		`INSERT INTO todos (item, completed, created_at, user_id) VALUES (?, ?, ?, ?)`,
		todo.Item,
		todo.Completed,
		time.Now(),
		todo.UserID,
	)

	return err
}

// GetAllTodos gets all todo items from the database for a specific user
func GetAllTodos(userID int) ([]Todo, error) {
	database := config.Database()

	rows, err := database.Query(
		`SELECT id, item, completed, created_at, user_id FROM todos WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo

		err = rows.Scan(&todo.Id, &todo.Item, &todo.Completed, &todo.CreatedAt, &todo.UserID)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

// GetTodoByID retrieves a todo by ID and checks user ownership
func GetTodoByID(id int, userID int) (Todo, error) {
	database := config.Database()
	var todo Todo

	err := database.QueryRow(
		`SELECT id, item, completed, created_at, user_id FROM todos WHERE id = ? AND user_id = ?`,
		id, userID,
	).Scan(&todo.Id, &todo.Item, &todo.Completed, &todo.CreatedAt, &todo.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return todo, errors.New("todo not found or you don't have permission")
		}
		return todo, err
	}

	return todo, nil
}

// UpdateTodo updates an existing todo item in the database
func UpdateTodo(todo Todo) error {
	database := config.Database()
	
	// First verify the todo belongs to the user
	_, err := GetTodoByID(todo.Id, todo.UserID)
	if err != nil {
		return err
	}

	_, err = database.Exec(
		`UPDATE todos SET item = ?, completed = ? WHERE id = ? AND user_id = ?`,
		todo.Item, todo.Completed, todo.Id, todo.UserID,
	)

	return err
}

// CompleteTodo marks a todo item as completed
func CompleteTodo(id int, userID int) error {
	database := config.Database()

	// First verify the todo belongs to the user
	_, err := GetTodoByID(id, userID)
	if err != nil {
		return err
	}

	_, err = database.Exec(
		`UPDATE todos SET completed = 1 WHERE id = ? AND user_id = ?`,
		id, userID,
	)

	return err
}

// DeleteTodo removes a todo item from the database
func DeleteTodo(id int, userID int) error {
	database := config.Database()

	// First verify the todo belongs to the user
	_, err := GetTodoByID(id, userID)
	if err != nil {
		return err
	}

	_, err = database.Exec(
		`DELETE FROM todos WHERE id = ? AND user_id = ?`,
		id, userID,
	)

	return err
}

// InitTodoTable ensures the todos table exists and has the correct structure
func InitTodoTable() {
	database := config.Database()

	_, err := database.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INT AUTO_INCREMENT,
			item TEXT NOT NULL,
			completed INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id INT NOT NULL,
			PRIMARY KEY (id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);
	`)

	if err != nil {
		// Log error but don't stop execution
		// The table might already exist
		// We don't want to halt the application for this
		// just keep going
	}
}