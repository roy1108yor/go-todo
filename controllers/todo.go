package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ichtrojan/go-todo/config"
	"github.com/ichtrojan/go-todo/middleware"
	"github.com/ichtrojan/go-todo/models"
	"html/template"
	"net/http"
	"strconv"
)

var (
	view     = template.Must(template.ParseFiles("./views/index.html"))
	database = config.Database()
)

// GetTodosHandler displays all todos for the current user
func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUser(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get todos for the current user
	todos, err := models.GetAllTodos(user.ID)
	if err != nil {
		fmt.Println(err)
	}

	// Prepare view data
	data := models.View{
		Todos: todos,
		User:  user,
	}

	_ = view.Execute(w, data)
}

// CreateTodoHandler adds a new todo item for the current user
func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUser(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	item := r.FormValue("item")

	// Create todo with user association
	todo := models.Todo{
		Item: item,
		Completed: 0,
		UserID: user.ID,
	}

	err := models.CreateTodo(todo)
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// DeleteTodoHandler removes a todo item if it belongs to the current user
func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUser(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	// Delete todo, this will verify ownership
	err = models.DeleteTodo(id, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// CompleteTodoHandler marks a todo item as complete if it belongs to the current user
func CompleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUser(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	// Complete todo, this will verify ownership
	err = models.CompleteTodo(id, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// UpdateTodoHandler handles the PUT request to update an existing todo item
func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUser(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
		return
	}

	// Set response content type
	w.Header().Set("Content-Type", "application/json")

	// Get the id parameter from request URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid ID format"})
		return
	}

	// Parse request body
	var todoData struct {
		Item string `json:"item"`
	}

	// Decode JSON body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&todoData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Try to get the existing todo to verify ownership
	existingTodo, err := models.GetTodoByID(id, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// Create todo object with updated data
	todo := models.Todo{
		Id: id,
		Item: todoData.Item,
		Completed: existingTodo.Completed,
		CreatedAt: existingTodo.CreatedAt, // Preserve the original creation time
		UserID: user.ID,                // Set the user ID for ownership verification
	}

	// Update the todo in database
	err = models.UpdateTodo(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update todo"})
		fmt.Println(err)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo updated successfully"})
}
