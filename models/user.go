package models

// User represents a user entity on the platform.
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`

	isAdmin string `json:"-"`
}
