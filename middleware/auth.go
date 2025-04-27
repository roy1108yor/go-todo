package middleware

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/ichtrojan/go-todo/models"
	"net/http"
	"os"
)

// Key for user context
type contextKey string

const UserContextKey contextKey = "user"

// Store is a session store for authentication
var Store = sessions.NewCookieStore([]byte(getSessionKey()))

// getSessionKey returns the session key from environment or a default value
func getSessionKey() string {
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		// In production, this should be properly set
		key = "todo-app-secret-key-change-in-production"
	}
	return key
}

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session
		session, _ := Store.Get(r, "auth-session")
		
		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		// Get user ID from session
		userID, ok := session.Values["user_id"].(int)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		// Get user information from database
		user, err := models.GetUserByID(userID)
		if err != nil {
			// If user doesn't exist anymore, clear session and redirect to login
			session.Values["authenticated"] = false
			session.Values["user_id"] = nil
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		
		// Add user to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		
		// Call the next handler with our modified context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUser extracts the user from the request context
func GetUser(r *http.Request) (models.User, bool) {
	user, ok := r.Context().Value(UserContextKey).(models.User)
	return user, ok
}

// RequireAuth is a wrapper for routes that require authentication
func RequireAuth(handler http.HandlerFunc) http.Handler {
	return AuthMiddleware(http.HandlerFunc(handler))
}