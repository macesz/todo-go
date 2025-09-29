package domain

// TodoDTO is a Data Transfer Object for Todo.
// It's used to transfer data in a format suitable for APIs (like JSON).
// Similar to a Java DTO class or a JS object used in APIs.

type TodoDTO struct {
	ID        int64  `json:"id"` // json tag for easy JSON encoding (like @JsonProperty in Java)
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"createdAt"`
}

type CreateTodoDTO struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

type UpdateTodoDTO struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
	Done  bool   `json:"done" validate:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
