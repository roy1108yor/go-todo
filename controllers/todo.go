package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ichtrojan/go-todo/config"
	"github.com/ichtrojan/go-todo/models"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	id        int
	item      string
	completed int
	view      = template.Must(template.ParseFiles("./views/index.html"))
	database  = config.Database()
)

func Show(w http.ResponseWriter, r *http.Request) {
	statement, err := database.Query(`SELECT * FROM todos`)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	todos := make([]models.Todo, 0)

	for statement.Next() {
		err = statement.Scan(&id, &item, &completed)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		todo := models.Todo{
			Id:        id,
			Item:      item,
			Completed: completed,
		}

		todos = append(todos, todo)
	}

	data := models.View{
		Todos: todos,
	}

	if err = view.Execute(w, data); err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	item := r.FormValue("item")

	// Validate empty input
	if item == "" || len(strings.TrimSpace(item)) == 0 {
		http.Redirect(w, r, "/", 301)
		return
	}

	_, err := database.Exec(`INSERT INTO todos (item) VALUE (?)`, item)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec(`DELETE FROM todos WHERE id = ?`, id)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

func Complete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec(`UPDATE todos SET completed = 1 WHERE id = ?`, id)

	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", 301)
}

// Update handles updating an existing todo task
func Update(w http.ResponseWriter, r *http.Request) {
	// Parse the ID from URL parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// Read request body for updated task content
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Parse body as JSON
	var requestData struct {
		Item string `json:"item"`
	}
	
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate that the updated task is not empty
	updatedItem := requestData.Item
	if updatedItem == "" || len(strings.TrimSpace(updatedItem)) == 0 {
		http.Error(w, "Task content cannot be empty", http.StatusBadRequest)
		return
	}

	// Update the task in the database
	_, err = database.Exec(`UPDATE todos SET item = ? WHERE id = ?`, updatedItem, id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Task updated successfully"})
}