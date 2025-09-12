package infiletodo

import (
	"bufio"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync" // For thread-safety (like synchronized in Java or mutex in JS)
	"time"

	"github.com/macesz/todo-go/domain"
)

// TodoStore manages a collection of Todos in a file.
// It's like a Java HashMap<Integer, Todo> with methods.
type InFileStore struct {
	mu       sync.RWMutex        // Mutex for safe concurrent access (Go's goroutines are like threads)
	nextID   int                 // Auto-increment ID (like a database sequence)
	data     map[int]domain.Todo // map is like Java HashMap or JS object {}
	filePath string              // Path to the file where todos are stored
}

// NewTodoStore creates a new store instance.
// Like a constructor in Java or new Store() in JS.
// NewInFileStore constructs the store and loads existing todos from file.
// If the file doesn't exist, it will be created (empty).
func NewInFileStore(filePath string) (*InFileStore, error) {
	// Ensure the directory exists (e.g., "/todos/")
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}

	// Initialize the store
	store := &InFileStore{
		nextID:   1,
		data:     make(map[int]domain.Todo),
		filePath: filePath,
	}

	// Initialize empty file if not present
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("init data file: %w", err)
		}
		_ = f.Close()
		return store, nil
	} else if err != nil {
		return nil, err
	}

	if err := store.loadFromFile(); err != nil {
		return nil, err
	}
	return store, nil
}

// saveToFileLocked writes the current in-memory data to disk atomically.
// IMPORTANT: Caller must hold s.mu.Lock() (write lock).
func (s *InFileStore) saveToFileLocked() error {
	// Get the directory of the file
	dir := filepath.Dir(s.filePath)

	// Create a temp file in the same directory
	tmp, err := os.CreateTemp(dir, "todos-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	// Write CSV data to temp file
	tmpName := tmp.Name()

	writer := csv.NewWriter(tmp) // buffered writer for efficiency

	// Stable ordering by ID for predictable diffs
	ids := make([]int, 0, len(s.data)) // preallocate slice
	for id := range s.data {
		ids = append(ids, id)
	}

	// Sort IDs to ensure consistent order
	sort.Ints(ids)

	// Write each todo as a CSV record
	for _, id := range ids {
		todo := s.data[id] // get todo by id, (s -> InFileStore)
		rec := []string{   // CSV record as slice of strings
			strconv.Itoa(todo.ID),                     // convert int to string
			todo.Title,                                // Title is already a string
			strconv.FormatBool(todo.Done),             // convert bool to string
			todo.CreatedAt.UTC().Format(time.RFC3339), // format time to string in RFC3339
		}
		// Write the record
		// If write fails, clean up temp file and return error
		if err := writer.Write(rec); err != nil {
			_ = tmp.Close()        // close temp file
			_ = os.Remove(tmpName) // remove temp file
			return fmt.Errorf("write csv: %w", err)
		}
	}
	writer.Flush() // flush buffered data to underlying writer

	// Check for errors during flush
	if err := writer.Error(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("flush csv: %w", err)
	}

	// Ensure data is on disk before rename
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("fsync temp file: %w", err)
	}

	// Close the temp file
	// Closing also flushes, but we already flushed above
	// We check for close errors separately to handle them
	// (e.g., disk full errors may appear on close)
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp file: %w", err)
	}

	// Atomic replace of the original file with the temp file
	// os.Rename is atomic on POSIX systems if source and target are on the same filesystem
	if err := os.Rename(tmpName, s.filePath); err != nil {
		_ = os.Remove(tmpName) // remove temp file on error
		return fmt.Errorf("atomic replace: %w", err)
	}
	return nil
}

// loadFromFile loads all todos into memory.
// Holds the write lock while replacing the map and computing nextID.
func (s *InFileStore) loadFromFile() error {
	s.mu.Lock()         // Hold write lock during load to prevent access to partial data
	defer s.mu.Unlock() // defer ensures unlock happens (like finally in Java)

	// Open the file for reading
	f, err := os.Open(s.filePath) // open for read-only
	if err != nil {
		return fmt.Errorf("open data file: %w", err)
	}
	defer f.Close() // ensure file is closed

	// Use a buffered reader for efficiency
	br := bufio.NewReader(f)

	// Peek to see if the file is empty
	peek, err := br.Peek(1)
	if err != nil && err != io.EOF {
		return fmt.Errorf("peek data file: %w", err)
	}

	// If empty, initialize empty map and return
	if len(peek) == 0 {
		// Empty file, nothing to load
		s.data = make(map[int]domain.Todo) // reset data map
		s.nextID = 1                       // reset nextID
		return nil
	}

	// Reset reader to start of file, since Peek advanced it by 1 byte
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek data file: %w", err)
	}

	// Read CSV records
	r := csv.NewReader(f)  // CSV reader
	r.FieldsPerRecord = -1 // allow variable fields; we will validate manually

	// Read all records at once
	// Could also read line by line with r.Read() in a loop
	// but ReadAll is simpler for small files
	records, err := r.ReadAll() // read all records at once
	if err != nil {
		return fmt.Errorf("read todos file: %w", err)
	}

	// Prepare to load data
	s.data = make(map[int]domain.Todo, len(records)) // reset data map
	s.nextID = 1

	// Parse records
	for i, rec := range records {
		if len(rec) == 0 {
			continue
		}
		if len(rec) != 4 {
			return fmt.Errorf("invalid record on line %d: expected 4 fields, got %d", i+1, len(rec))
		}

		// Parse each field with error handling
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			return fmt.Errorf("parse id on line %d: %w", i+1, err)
		}
		title := rec[1]
		done, err := strconv.ParseBool(rec[2])
		if err != nil {
			return fmt.Errorf("parse done on line %d: %w", i+1, err)
		}
		createdAt, err := time.Parse(time.RFC3339, rec[3])
		if err != nil {
			return fmt.Errorf("parse createdAt on line %d: %w", i+1, err)
		}

		// Add to map
		s.data[id] = domain.Todo{
			ID:        id,
			Title:     title,
			Done:      done,
			CreatedAt: createdAt,
		}
		// Update nextID to be one more than the highest ID seen
		if id >= s.nextID {
			s.nextID = id + 1
		}
	}
	return nil
}

// Create adds a new Todo with the given title.
func (s *InFileStore) Create(_ context.Context, title string) (domain.Todo, error) {
	// Create a new Todo with the given title and default values
	todo := domain.Todo{
		ID:        0,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now().UTC(),
	}
	// Validate the Todo before creating it
	if err := todo.Validate(); err != nil {
		return domain.Todo{}, err
	}

	s.mu.Lock()         // Lock for writing
	defer s.mu.Unlock() // defer ensures unlock happens

	// Assign the next ID and increment
	todo.ID = s.nextID
	s.nextID++
	s.data[todo.ID] = todo

	// Persist to disk immediately
	if err := s.saveToFileLocked(); err != nil {
		// Roll back in-memory state if disk write fails
		delete(s.data, todo.ID)
		s.nextID--
		return domain.Todo{}, err
	}
	return todo, nil
}

// List returns all Todos sorted by ID ascending.
func (s *InFileStore) List(_ context.Context) ([]domain.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]int, 0, len(s.data))
	for id := range s.data {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	todos := make([]domain.Todo, 0, len(ids))
	for _, id := range ids {
		todos = append(todos, s.data[id])
	}
	return todos, nil
}

// Get retrieves a Todo by ID.
func (s *InFileStore) Get(_ context.Context, id int) (domain.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.data[id]
	if !ok {
		return domain.Todo{}, errors.New("todo not found")
	}

	return todo, nil
}

// Update modifies an existing Todo by ID.
func (s *InFileStore) Update(_ context.Context, id int, title string, done bool) (domain.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.data[id]
	if !ok {
		return domain.Todo{}, errors.New("todo not found")
	}

	todo.Title = title
	todo.Done = done
	if err := todo.Validate(); err != nil {
		return domain.Todo{}, err
	}

	s.data[id] = todo

	if err := s.saveToFileLocked(); err != nil {
		return domain.Todo{}, err
	}
	return todo, nil
}

// Delete removes a Todo by ID.
func (s *InFileStore) Delete(_ context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return errors.New("todo not found")
	}
	delete(s.data, id)

	if err := s.saveToFileLocked(); err != nil {
		// Could also consider restoring the item on error
		return err
	}

	return nil
}
