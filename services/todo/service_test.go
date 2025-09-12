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
