package domain

// TodoDTO is a Data Transfer Object for Todo.
// It's used to transfer data in a format suitable for APIs (like JSON).
// Similar to a Java DTO class or a JS object used in APIs.

type ErrorResponse struct {
	Error string `json:"error"`
}

// TodoList
type TodoListDTO struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`

	Title     string    `json:"title"`
	Color     *string   `json:"color,omitempty"`
	Labels    []string  `json:"labels,omitempty"`
	CreatedAt string    `json:"created_at"`
	Items     []TodoDTO `json:"items,omitempty"`
}

type CreateTodoListRequestDTO struct {
	Title  string   `json:"title"`
	Color  *string  `json:"color,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

type UpdateTodoListRequestDTO struct {
	Title  string   `json:"title,omitempty"`
	Color  *string  `json:"color,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

// TODO
type TodoDTO struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userID"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	Priority  int64  `json:"priority"`
	CreatedAt string `json:"created_at"`
}

type CreateTodoDTO struct {
	Title    string `json:"title" validate:"required,min=1,max=255"`
	Priority int64  `json:"priority" validate:"required,min=1,max=5"`
}

type UpdateTodoDTO struct {
	Title    string `json:"title" validate:"required,min=1,max=255"`
	Done     bool   `json:"done" validate:"required"`
	Priority int64  `json:"priority" validate:"required,min=1,max=5"`
}

// User
type UserDTO struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequestDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=255,containsany=0123456789,containsany=ABCDEFGHIJKLMNOPQRSTUVWXYZ"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}
