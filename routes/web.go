package routes

import (
	"github.com/gorilla/mux"
	"github.com/ichtrojan/go-todo/controllers"
	"github.com/ichtrojan/go-todo/middleware"
	"net/http"
)

func Init() *mux.Router {
	route := mux.NewRouter()

	// Public routes
	route.HandleFunc("/login", controllers.LoginFormHandler).Methods("GET")
	route.HandleFunc("/login", controllers.LoginHandler).Methods("POST")
	route.HandleFunc("/register", controllers.RegisterFormHandler).Methods("GET")
	route.HandleFunc("/register", controllers.RegisterHandler).Methods("POST")

	// Root path redirects to home
	route.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/todos", http.StatusSeeOther)
	})

	// Protected routes
	todoRoutes := route.PathPrefix("/todos").Subrouter()
	todoRoutes.Use(middleware.AuthMiddleware)
	todoRoutes.HandleFunc("", controllers.GetTodosHandler).Methods("GET")
	todoRoutes.HandleFunc("", controllers.CreateTodoHandler).Methods("POST")
	todoRoutes.HandleFunc("/{id}", controllers.UpdateTodoHandler).Methods("PUT")
	todoRoutes.HandleFunc("/delete/{id}", controllers.DeleteTodoHandler)
	todoRoutes.HandleFunc("/complete/{id}", controllers.CompleteTodoHandler)

	// Logout route (requires authentication)
	authRoutes := route.PathPrefix("").Subrouter()
	authRoutes.Use(middleware.AuthMiddleware)
	authRoutes.HandleFunc("/logout", controllers.LogoutHandler).Methods("GET")

	return route
}