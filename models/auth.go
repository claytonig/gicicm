package models

// LoginRequest represents a login requst
type LoginRequest struct {
	Username string
	Password string
}

// SignUpRequest requests represents a request
// for a creation of a new user.
type SignUpRequest struct {
	Email    string
	Password string
	Name     string
}
