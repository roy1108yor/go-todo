package controllers

import (
	"fmt"
	"github.com/ichtrojan/go-todo/middleware"
	"github.com/ichtrojan/go-todo/models"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var (
	loginTemplate    = template.Must(template.ParseFiles("./views/login.html"))
	registerTemplate = template.Must(template.ParseFiles("./views/register.html"))
)

// RegisterFormHandler displays the registration form
func RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	data := models.View{}
	_ = registerTemplate.Execute(w, data)
}

// RegisterHandler processes the registration form submission
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusInternalServerError)
		return
	}

	// Get form values
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// Basic validation
	var errorMessages []string

	if username == "" {
		errorMessages = append(errorMessages, "Username is required")
	}

	if email == "" {
		errorMessages = append(errorMessages, "Email is required")
	}

	if password == "" {
		errorMessages = append(errorMessages, "Password is required")
	} else if len(password) < 6 {
		errorMessages = append(errorMessages, "Password must be at least 6 characters")
	}

	if password != confirmPassword {
		errorMessages = append(errorMessages, "Passwords do not match")
	}

	// Check if email already exists
	if email != "" {
		_, err = models.FindUserByEmail(email)
		if err == nil {
			errorMessages = append(errorMessages, "Email already registered")
		}
	}

	// If validation fails, show form with errors
	if len(errorMessages) > 0 {
		data := models.View{
			Error:  true,
			Flash:  strings.Join(errorMessages, "<br>"),
			Values: map[string]string{"username": username, "email": email},
		}
		_ = registerTemplate.Execute(w, data)
		return
	}

	// Create user
	user := models.User{
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	userID, err := models.CreateUser(user)
	if err != nil {
		data := models.View{
			Error: true,
			Flash: "Error creating user: " + err.Error(),
			Values: map[string]string{"username": username, "email": email},
		}
		_ = registerTemplate.Execute(w, data)
		return
	}

	// Create session
	session, _ := middleware.Store.Get(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["user_id"] = int(userID)
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LoginFormHandler displays the login form
func LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	// Check if already logged in
	session, _ := middleware.Store.Get(r, "auth-session")
	if auth, ok := session.Values["authenticated"].(bool); ok && auth {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := models.View{}

	// Check for redirect messages
	if r.URL.Query().Get("registered") == "1" {
		data.Success = true
		data.Flash = "Registration successful! Please log in."
	}

	_ = loginTemplate.Execute(w, data)
}

// LoginHandler processes the login form submission
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Could not parse form", http.StatusInternalServerError)
		return
	}

	// Get form values
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	// Basic validation
	if email == "" || password == "" {
		data := models.View{
			Error: true,
			Flash: "Email and password are required",
			Values: map[string]string{"email": email},
		}
		_ = loginTemplate.Execute(w, data)
		return
	}

	// Check credentials
	user, err := models.FindUserByEmail(email)
	if err != nil || !user.CheckPassword(password) {
		data := models.View{
			Error: true,
			Flash: "Invalid email or password",
			Values: map[string]string{"email": email},
		}
		_ = loginTemplate.Execute(w, data)
		return
	}

	// Create session
	session, _ := middleware.Store.Get(r, "auth-session")
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Get session
	session, _ := middleware.Store.Get(r, "auth-session")

	// Remove authentication values
	session.Values["authenticated"] = false
	session.Values["user_id"] = nil

	// Save session
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// HomeHandler shows the protected home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	user, ok := middleware.GetUser(r)
	if !ok {
		fmt.Println("Failed to get user from context")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Get all todos for the current user
	todos, err := models.GetAllTodos(user.ID)
	if err != nil {
		fmt.Println(err)
	}

	// Prepare view data
	data := models.View{
		Todos: todos,
		User:  user,
	}

	// Render template
	indexTemplate := template.Must(template.ParseFiles("./views/index.html"))
	_ = indexTemplate.Execute(w, data)
}
