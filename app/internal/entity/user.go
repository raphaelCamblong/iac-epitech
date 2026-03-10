package entity

// User represents a user in the domain.
type User struct {
	ID           string
	Email        string
	PasswordHash string
}
