package entity

import "time"

// Task represents a task in the domain.
type Task struct {
	ID               string
	UserID           string
	Title            string
	Content          string
	DueDate          time.Time
	Done             bool
	RequestTimestamp time.Time // For out-of-order request conflict detection
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
