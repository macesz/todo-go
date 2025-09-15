package domain

// TodoDTO is a Data Transfer Object for Todo.
// It's used to transfer data in a format suitable for APIs (like JSON).
// Similar to a Java DTO class or a JS object used in APIs.

type TodoDTO struct {
	ID        int    `json:"id"` // json tag for easy JSON encoding (like @JsonProperty in Java)
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"createdAt"`
}
