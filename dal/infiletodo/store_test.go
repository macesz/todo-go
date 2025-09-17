package infiletodo

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath" // For thread-safety (like synchronized in Java or mutex in JS)
	"sync"
	"testing"
	"time"

	"github.com/macesz/todo-go/domain"
	"github.com/stretchr/testify/require"
)

func TestSaveToFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "todos.csv")

	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	s := &InFileStore{
		data: map[int]domain.Todo{
			1: {ID: 1, Title: "Test Todo 1", Done: false, CreatedAt: ts},
			2: {ID: 2, Title: "Test Todo 2", Done: true, CreatedAt: ts},
		},
		filePath: file,
	}

	err := s.saveToFileLocked()
	require.NoError(t, err)

	f, err := os.Open(file)
	require.NoError(t, err)
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	rows, err := r.ReadAll()
	require.NoError(t, err)
	require.Len(t, rows, 2)

	for _, row := range rows {
		require.Len(t, row, 4)
		id := row[0]
		title := row[1]
		done := row[2]
		createdAt := row[3]

		_, perr := time.Parse(time.RFC3339, row[3])
		require.NoError(t, perr, fmt.Sprintf("invalid timestamp: %q", row[3]))

		if id == "1" {
			require.Equal(t, "Test Todo 1", title)
			require.Equal(t, "false", done)
			require.Equal(t, ts.Format(time.RFC3339), createdAt)
		} else if id == "2" {
			require.Equal(t, "Test Todo 2", title)
			require.Equal(t, "true", done)
			require.Equal(t, ts.Format(time.RFC3339), createdAt)
		} else {
			t.Errorf("unexpected ID: %s", id)
		}
	}
}

func TestLoadFromFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "todos.csv")

	// Prepare file content
	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC).Format(time.RFC3339)
	content := "1,Todo 1,false," + ts + "\n" + "2,Todo 2,true," + ts + "\n"
	require.NoError(t, os.WriteFile(file, []byte(content), 0o600))

	s := &InFileStore{
		data:     make(map[int]domain.Todo),
		filePath: file,
	}
	err := s.loadFromFile()
	require.NoError(t, err)

	require.Len(t, s.data, 2)
	require.Equal(t, "Todo 1", s.data[1].Title)
	require.Equal(t, false, s.data[1].Done)
	require.Equal(t, "Todo 2", s.data[2].Title)
	require.Equal(t, true, s.data[2].Done)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	type fields struct {
		filePath string
		data     map[int]domain.Todo
	}

	type args struct {
		ctx   context.Context
		title string
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantInFile []byte
		want       domain.Todo
		wantErr    bool
	}{
		{
			name: "Create Todo",
			fields: fields{
				filePath: filepath.Join(os.TempDir(), "test_todos_create.csv"), // Use temp file for testing
				data:     make(map[int]domain.Todo),
			},
			args: args{ctx: context.Background(), title: "Test Todo"},
			want: domain.Todo{
				ID:        0,
				Title:     "Test Todo",
				Done:      false,
				CreatedAt: time.Now(), // We will check this separately
			},
			// We expect the file to contain the CSV representation of the todo
			wantInFile: []byte("0,Test Todo,false,"), // CreatedAt will be appended, so we check prefix only
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &InFileStore{
				data:     tt.fields.data,
				filePath: tt.fields.filePath,
			}

			got, err := s.Create(tt.args.ctx, tt.args.title)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("Create() got ID = %v, want %v", got.ID, tt.want.ID)
			}
			if got.Title != tt.want.Title {
				t.Errorf("Create() got Title = %v, want %v", got.Title, tt.want.Title)
			}
			if got.Done != tt.want.Done {
				t.Errorf("Create() got Done = %v, want %v", got.Done, tt.want.Done)
			}

			// Check file contents
			data, err := os.ReadFile(tt.fields.filePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if !bytes.HasPrefix(data, tt.wantInFile) {
				t.Errorf("file contents = %s, want prefix %s", data, tt.wantInFile)
			}
		})
	}
}

func TestList(t *testing.T) {
	t.Parallel()

	type fields struct {
		filePath string
		data     map[int]domain.Todo
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Todo
		wantErr bool
	}{
		{
			name: "List Todos",
			fields: fields{
				filePath: filepath.Join(os.TempDir(), "test_todos_list.csv"), // Use temp file for testing
				// Initialize with some data to test listing functionality without file I/O complexity here (assuming Load is tested separately)
				data: map[int]domain.Todo{
					1: {ID: 1, Title: "Todo 1", Done: false, CreatedAt: time.Now()},
					2: {ID: 2, Title: "Todo 2", Done: true, CreatedAt: time.Now()},
				},
			},
			args: args{ctx: context.Background()},
			want: []domain.Todo{
				{ID: 1, Title: "Todo 1", Done: false},
				{ID: 2, Title: "Todo 2", Done: true},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &InFileStore{
				data:     tt.fields.data,
				filePath: tt.fields.filePath,
			}

			got, err := s.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("List() got length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i].ID != tt.want[i].ID || got[i].Title != tt.want[i].Title || got[i].Done != tt.want[i].Done {
					t.Errorf("List() got = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestCorruptedFile_ReturnsError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "todos.csv")

	require.NoError(t, os.WriteFile(path, []byte("{not-csv"), 0o600))

	s := &InFileStore{
		filePath: path,
		data:     make(map[int]domain.Todo),
	}

	// If your store loads on demand, call the loader explicitly
	err := s.loadFromFile()
	require.Error(t, err, "expected error due to corrupted file")
}

func TestConcurrentCreateUniqIds(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "todos.csv")

	s := &InFileStore{
		filePath: path,
		data:     make(map[int]domain.Todo),
	}

	const n = 200
	ids := make(chan int, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			todo, err := s.Create(context.Background(), "Concurrent Todo")
			require.NoError(t, err)
			ids <- todo.ID
		}()
	}
	wg.Wait()
	close(ids)

	seen := map[int]struct{}{}
	for id := range ids {
		if _, ok := seen[id]; ok {
			t.Fatalf("duplicate id: %d", id)
		}
		seen[id] = struct{}{}
	}
}
