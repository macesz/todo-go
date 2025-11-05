package inmemorytodo

import (
	"context"
	"errors"
	"sync" // For thread-safety (like synchronized in Java or mutex in JS)
	"time"

	"github.com/macesz/todo-go/domain"
)

// TodoStore manages a collection of Todos in memory.
// It's like a Java HashMap<Integer, Todo> with methods.
type InMemoryStore struct {
	mu     sync.RWMutex          // Mutex for safe concurrent access (Go's goroutines are like threads)
	nextID int64                 // Auto-increment ID (like a database sequence)
	data   map[int64]domain.Todo // map is like Java HashMap or JS object {}
}

// NewTodoStore creates a new store instance.
// Like a constructor in Java or new Store() in JS.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{nextID: 1, data: make(map[int64]domain.Todo)} // make() initializes the map
}

//Here starts all the receiver methods on *TodoStore (pointer for modifications)

// Create adds a new Todo.
// Returns the created Todo or an error.
func (s *InMemoryStore) Create(ctx context.Context, title string) (*domain.Todo, error) {
	// Create a new Todo with the given title and default values
	t := domain.Todo{ID: 0, Title: title, Done: false, CreatedAt: time.Now().UTC()} // time.Now() like new Date() in JS

	// Validate the Todo before creating it
	if err := t.Validate(); err != nil { // Call the receiver method
		return nil, err
	}

	s.mu.Lock()         // Lock for writing (like synchronized block in Java)
	defer s.mu.Unlock() // defer ensures unlock happens (like finally in Java)
	t.ID = s.nextID     // assign the next ID to the Todo
	s.nextID++          // increment the next ID
	s.data[t.ID] = t    // store the Todo in the map
	return &t, nil      // return the created Todo and no error
}

// List returns all Todos
func (s *InMemoryStore) List(ctx context.Context) ([]*domain.Todo, error) {
	s.mu.RLock()                                  // Read lock (like synchronized block in Java)
	defer s.mu.RUnlock()                          // defer ensures unlock happens (like finally in Java)
	todos := make([]*domain.Todo, 0, len(s.data)) // Todo is a slice of Todo structs like an array in JS
	for _, t := range s.data {                    // range is like for (let key in obj) in JS
		todos = append(todos, &t) // append() is like push() in JS
	}
	return todos, nil
}

// Get retrieves a Todo by ID
func (s *InMemoryStore) Get(ctx context.Context, id int64) (*domain.Todo, error) {
	s.mu.RLock()         // Read lock (like synchronized block in Java)
	defer s.mu.RUnlock() // defer ensures unlock happens (like finally in Java)
	t, ok := s.data[id]  // map lookup is like obj[key] in JS, ok is true if the key exists
	if !ok {
		return nil, errors.New("does not exists")
	}
	return &t, nil
}

//Update modifies an existing Todo

func (s *InMemoryStore) Update(ctx context.Context, id int64, title string, done bool) (*domain.Todo, error) {
	s.mu.Lock()         // Write lock (like synchronized block in Java)
	defer s.mu.Unlock() // defer ensures unlock happens (like finally in Java)
	t, ok := s.data[id] // map lookup is like obj[key] in JS, ok is true if the key exists
	if !ok {
		return nil, errors.New("todo not found")
	}
	t.Title = title
	t.Done = done
	if err := t.Validate(); err != nil { // Call the receiver method
		return nil, err
	}
	s.data[id] = t // update the Todo in the map
	return &t, nil // return the updated Todo and no error
}

// Delete removes a Todo by ID

func (s *InMemoryStore) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()         // Write lock (like synchronized block in Java)
	defer s.mu.Unlock() // defer ensures unlock happens (like finally in Java)
	if _, ok := s.data[id]; !ok {
		return errors.New("cant delete")
	}
	delete(s.data, id) // delete() is like delete() in JS or .remove() in Java
	return nil
}
