package models

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string
	Password string
}

// SignUpRequest requests represents a request
// for a creation of a new user.
type SignUpRequest struct {
	Email    string
	Password string
	Name     string
}
