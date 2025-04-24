package models

import (
	"github.com/ichtrojan/go-todo/config"
	"time"
)

type Todo struct {
	Id         int
	Item       string
	Completed  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CompletedAt *time.Time
}

// UpdateTodo updates an existing todo item in the database
func UpdateTodo(todo Todo) error {
	database := config.Database()

	_, err := database.Exec(`UPDATE todos SET item = ?, completed = ? WHERE id = ?`, 
		todo.Item, todo.Completed, todo.Id)

	return err
}