package models

import "github.com/ichtrojan/go-todo/config"

// Todo represents a task with its properties
type Todo struct {
	Id        int
	Item      string
	Completed int
}

// UpdateTask updates a task with the given ID in the database
func UpdateTask(id string, newContent string) error {
	// Import database package
	db := config.Database()
	defer db.Close()

	// Prepare SQL statement to update the task
	stmt, err := db.Prepare("UPDATE todos SET item = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the update statement
	_, err = stmt.Exec(newContent, id)
	if err != nil {
		return err
	}

	return nil
}