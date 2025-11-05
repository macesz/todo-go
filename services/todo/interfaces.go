package todo

import (
	"context"

	"github.com/macesz/todo-go/domain"
)

// TodoStore defines the interface for a todo storage backend. Like a Java interface
type TodoStore interface {
	List(ctx context.Context, userID int64) ([]*domain.Todo, error)
	Create(ctx context.Context, userID int64, title string, priority int64) (*domain.Todo, error)
	Get(ctx context.Context, id int64) (*domain.Todo, error)
	Update(ctx context.Context, id int64, title string, done bool, priority int64) (*domain.Todo, error)
	Delete(ctx context.Context, id int64) error
}

//********************************************************************************************

// A side note about the TodoStore interface, and about a refactor to an UPSERT, and how I faced the DAL Interface Dilemma

//Pros: Fewer methods, less code duplication
// Perfect for database UPSERT operations
// Simpler implementation

//Key Realizations:

// Clear Intent: Method names immediately communicate what's happening
// Type Safety: Create can't accidentally overwrite; Update can't accidentally create
// Different Contracts: Create needs minimal input, Update needs full control
// Go Philosophy: Explicit is better than implicit
// Integration Test Clarity: Each operation tests exactly what it claims
// The Lesson: Sometimes the "clever" solution isn't the best solution. Clear interfaces trump clever implementations, especially when multiple developers (including future you) need to understand the code quickly.

// *****************
//
// The Upsert method is used to create or update a todo item in the storage backend.
//
// type TodoStore interface {
// 	Save(ctx context.Context, todo *domain.Todo) error
// 	List(ctx context.Context) ([]*domain.Todo, error)
// 	Delete(ctx context.Context, id int64) error
// 	Get(ctx context.Context, id int64) (*domain.Todo, error)
// }
