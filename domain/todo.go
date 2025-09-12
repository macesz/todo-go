package domain

import (
	"errors" // For error handling (like Java's Exception)
	"time"   // For timestamps (like JS Date or Java LocalDateTime)
)

// Todo is a struct representing a single todo item.
// It's like a Java class with fields, or a JS object.
type Todo struct {
	ID        int
	Title     string
	Done      bool
	CreatedAt time.Time
}

// Validate is a receiver method (attached to Todo).
// In Java: like public void validate() in Todo class.
// In JS: like Todo.prototype.validate = function() { ... }
func (t *Todo) Validate() error {
	if len(t.Title) == 0 { // len() is like .length in JS
		return errors.New("title is required") // errors.New is like throw new Error() in JS or Java
	}
	return nil
}
