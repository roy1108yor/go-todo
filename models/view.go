package models

type View struct {
	Todos   []Todo
	User    User
	Error   bool
	Success bool
	Flash   string
	Values  map[string]string
}