package main // All files are in "main" package

import (
	"errors" // For error handling (like Java's Exception)
	"time"   // For timestamps (like JS Date or Java LocalDateTime)
)

// Todo is a struct representing a single todo item.
// It's like a Java class with fields, or a JS object.
type Todo struct {
	ID        int       `json:"id"` // json tag for easy JSON encoding (like @JsonProperty in Java)
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

// Validate is a receiver method (attached to Todo).
// In Java: like public void validate() in Todo class.
// In JS: like Todo.prototype.validate = function() { ... }
func (t *Todo) Validate() error {
	if len(t.Title) == 0 { // len() is like .length in JS
		return errors.New("title is required") // errors.New is like new Error() in JS
	}
	return nil // nil is like null in Java/JS
}
