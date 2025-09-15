package todo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/services/todo/mocks"
)

// TestListTodos tests the ListTodos method of the TodoService.
// It uses a mock TodoStore to simulate the data layer.
func TestListTodos(t *testing.T) {
	// Enable parallel execution of tests
	// This is useful when tests are independent and can run concurrently
	// It speeds up the test suite execution
	// Just make sure that shared resources are handled properly
	t.Parallel()

	// Define the fields of the TodoService struct
	// This allows us to set up different configurations for each test case
	type fields struct {
		Sotre *mocks.TodoStore
	}

	// Define the arguments for the ListTodos method
	// This allows us to pass different contexts for each test case
	type args struct {
		ctx context.Context
	}

	now := time.Now()

	// Define the test cases
	// Each test case has a name, fields, args, expected error flag, and expected result
	// The initMocks function is used to set up the mock behavior for each test case
	// This keeps the test cases clean and focused on the specific scenario being tested
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoService) // Function to initialize mocks
		want      []domain.Todo
	}{
		{
			name:   "success",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				// Set up the expected behavior of the mock store
				// When List is called with the given context, return a predefined list of todos
				store.On("List", ta.ctx).Return([]domain.Todo{
					{ID: 1, Title: "Test Todo 1", Done: false, CreatedAt: now},
					{ID: 2, Title: "Test Todo 2", Done: true, CreatedAt: now},
				}, nil).Once()

				s.Store = store
			},
			wantErr: false, // No error expected
			args: args{
				ctx: context.Background(), // Use a background context for simplicity
			},
			// Expected result
			want: []domain.Todo{
				{ID: 1, Title: "Test Todo 1", Done: false, CreatedAt: now},
				{ID: 2, Title: "Test Todo 2", Done: true, CreatedAt: now},
			},
		},
		{
			name:   "store error",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				store.On("List", ta.ctx).Return(nil, errors.New("not found")).Once()

				s.Store = store
			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tt.fields.Sotre,
			}

			tt.initMocks(t, &tt.args, s)

			res, err := s.ListTodos(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTodos() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// First check length
			if len(res) != len(tt.want) {
				t.Errorf("ListTodos() got = %v, want %v", res, tt.want)
				return
			}

			// Deep compare slices
			for i := range res {
				if res[i] != tt.want[i] {
					t.Errorf("ListTodos() got = %v, want %v", res, tt.want)
				}
			}
		})
	}
}

func TestCreateTodo(t *testing.T) {
	t.Parallel()

	// Define the fields of the TodoService struct
	type fields struct {
		Store *mocks.TodoStore
	}

	// Define the arguments for the CreateTodo method
	type args struct {
		ctx   context.Context
		title string
	}

	now := time.Now()

	// Define the test cases
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoService) // Function to initialize mocks
		want      domain.Todo
	}{
		{
			name:   "success",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				// Set up the expected behavior of the mock store, "ta." means test args
				// When Create is called with the given context and title, return a predefined todo
				store.On("Create", ta.ctx, ta.title).Return(domain.Todo{
					ID:        1,
					Title:     ta.title,
					Done:      false,
					CreatedAt: now,
				}, nil).Once()

				s.Store = store
			},
			wantErr: false,
			args: args{
				ctx:   context.Background(),
				title: "New Todo",
			},
			want: domain.Todo{
				ID:        1,
				Title:     "New Todo",
				Done:      false,
				CreatedAt: now,
			},
		},
		{
			name:   "store error",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				store.On("Create", ta.ctx, ta.title).Return(domain.Todo{}, errors.New("could not create")).Once()

				s.Store = store
			},
			wantErr: true,
			args: args{
				ctx:   context.Background(),
				title: "New Todo",
			},
			want: domain.Todo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tt.fields.Store,
			}

			tt.initMocks(t, &tt.args, s)

			got, err := s.CreateTodo(tt.args.ctx, tt.args.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateTodo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTodo(t *testing.T) {
	t.Parallel()

	// Define the fields of the TodoService struct
	type fields struct {
		Store *mocks.TodoStore
	}

	// Define the arguments for the GetTodo method
	type args struct {
		ctx context.Context
		id  int
	}

	now := time.Now()

	// Define the test cases
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoService) // Function to initialize mocks
		want      domain.Todo
	}{
		{
			name:   "success",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				// Set up the expected behavior of the mock store
				// When Get is called with the given context and id, return a predefined todo
				store.On("Get", ta.ctx, ta.id).Return(domain.Todo{
					ID:        ta.id,
					Title:     "Test Todo",
					Done:      false,
					CreatedAt: now,
				}, nil).Once()

				s.Store = store
			},
			wantErr: false,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: domain.Todo{
				ID:        1,
				Title:     "Test Todo",
				Done:      false,
				CreatedAt: now,
			},
		},
		{
			name:   "not found",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				store.On("Get", ta.ctx, ta.id).Return(domain.Todo{}, errors.New("not found")).Once()

				s.Store = store
			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			want: domain.Todo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tt.fields.Store,
			}

			tt.initMocks(t, &tt.args, s)

			got, err := s.GetTodo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTodo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateTodo(t *testing.T) {

	t.Parallel()

	// Define the fields of the TodoService struct

	type fields struct {
		Store *mocks.TodoStore
	}

	// Define the arguments for the UpdateTodo method
	// This allows us to pass different contexts, ids, titles, and done statuses for each test case
	type args struct {
		ctx   context.Context
		id    int
		title string
		done  bool
	}

	now := time.Now()

	// Define the test cases

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoService) // Function to initialize mocks
		want      domain.Todo
	}{
		{
			name:   "success",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				// Set up the expected behavior of the mock store
				// When Update is called with the given context, id, title, and done status, return a predefined todo
				store.On("Update", ta.ctx, ta.id, ta.title, ta.done).Return(domain.Todo{
					ID:        ta.id,
					Title:     ta.title,
					Done:      ta.done,
					CreatedAt: now,
				}, nil).Once()

				s.Store = store
			},
			wantErr: false,
			args: args{
				ctx:   context.Background(),
				id:    1,
				title: "Updated Todo",
				done:  true,
			},
			want: domain.Todo{
				ID:        1,
				Title:     "Updated Todo",
				Done:      true,
				CreatedAt: now,
			},
		},
		{
			name:   "not found",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				store.On("Update", ta.ctx, ta.id, ta.title, ta.done).Return(domain.Todo{}, errors.New("not found")).Once()

				s.Store = store
			},
			wantErr: true,
			args: args{
				ctx:   context.Background(),
				id:    999,
				title: "Updated Todo",
				done:  true,
			},
			want: domain.Todo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tt.fields.Store,
			}

			tt.initMocks(t, &tt.args, s)

			got, err := s.UpdateTodo(tt.args.ctx, tt.args.id, tt.args.title, tt.args.done)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateTodo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {

	t.Parallel()

	// Define the fields of the TodoService struct
	type fields struct {
		Store *mocks.TodoStore
	}

	// Define the arguments for the DeleteTodo method
	// This allows us to pass different contexts and ids for each test case
	type args struct {
		ctx context.Context
		id  int
	}

	// Define the test cases

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoService) // Function to initialize mocks
	}{
		{
			name:   "success",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				// Set up the expected behavior of the mock store
				// When Delete is called with the given context and id, return nil (no error)
				store.On("Delete", ta.ctx, ta.id).Return(nil).Once()

				s.Store = store
			},
			wantErr: false,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
		},
		{
			name:   "not found",
			fields: fields{},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				store.On("Delete", ta.ctx, ta.id).Return(errors.New("not found")).Once()

				s.Store = store
			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
				id:  999,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tt.fields.Store,
			}

			tt.initMocks(t, &tt.args, s)

			err := s.DeleteTodo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTodo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
