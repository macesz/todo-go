package todo

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/services/todo/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var fixedTime = time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

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
		Store *mocks.TodoStore
	}

	// Define the arguments for the ListTodos method
	// This allows us to pass different contexts for each test case
	type args struct {
		ctx    context.Context
		userID int64
	}

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
		want      []*domain.Todo
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{ctx: context.Background()},
			want: []*domain.Todo{
				{ID: 1, UserID: 1, Title: "Test Todo 1", Done: false, Priority: 5, CreatedAt: fixedTime},
				{ID: 2, UserID: 1, Title: "Test Todo 2", Done: true, Priority: 5, CreatedAt: fixedTime},
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				store.On("List", ta.ctx, ta.userID).Return([]*domain.Todo{
					{ID: 1, UserID: 1, Title: "Test Todo 1", Done: false, Priority: 5, CreatedAt: fixedTime},
					{ID: 2, UserID: 1, Title: "Test Todo 2", Done: true, Priority: 5, CreatedAt: fixedTime},
				}, nil).Once()

				s.Store = store
			},
		},
		{
			name:    "store error",
			fields:  fields{},
			args:    args{ctx: context.Background()},
			wantErr: true,
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})
				store.On("List", ta.ctx, ta.userID).Return(nil, errors.New("could not list")).Once()

				s.Store = store
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.ListTodos(tc.args.ctx, tc.args.userID)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
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
		ctx      context.Context
		userId   int64
		title    string
		priority int64
	}

	// Define the test cases
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoService) // Function to initialize mocks
		validate  func(*testing.T, *args, *domain.Todo)
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{ctx: context.Background(), userId: 1, title: "New Todo", priority: 5},
			validate: func(t *testing.T, ta *args, todo *domain.Todo) {
				require.Equal(t, int64(1), todo.ID)
				require.Equal(t, ta.userId, todo.UserID)
				require.Equal(t, ta.title, todo.Title)
				require.Equal(t, ta.priority, todo.Priority)
				require.False(t, todo.Done)
				require.NotZero(t, todo.CreatedAt)
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Set up the expected behavior of the mock store, "ta." means test args
				// When Create is called with the given context and title, return a predefined todo
				store.On("Create", ta.ctx, mock.MatchedBy(
					func(todo *domain.Todo) bool {
						return todo.UserID == ta.userId &&
							todo.Title == ta.title &&
							todo.Priority == ta.priority
					})).Run(func(args mock.Arguments) {
					// Simulate the store setting the ID
					todo := args.Get(1).(*domain.Todo)
					todo.ID = 1
				}).Return(nil).Once()

				s.Store = store
			},
		},
		{
			name:    "store error",
			fields:  fields{},
			args:    args{ctx: context.Background(), title: "New Todo"},
			wantErr: true,
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Simulate an error from the store

				store.On("Create", ta.ctx, mock.MatchedBy(
					func(todo *domain.Todo) bool {
						return todo.UserID == ta.userId &&
							todo.Title == ta.title &&
							todo.Priority == ta.priority
					})).Run(func(args mock.Arguments) {
					// Simulate the store setting the ID
					todo := args.Get(1).(*domain.Todo)
					todo.ID = 1
				}).Return(fmt.Errorf("test error")).Once()
				s.Store = store
			},
			validate: func(t *testing.T, ta *args, todo *domain.Todo) {
				require.Nil(t, todo)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.CreateTodo(tc.args.ctx, tc.args.userId, tc.args.title, tc.args.priority)

			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			if tc.validate != nil {
				tc.validate(t, &tc.args, got)
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
		ctx    context.Context
		userId int64
		id     int64
	}

	now := time.Now()

	// Define the test cases
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoService) // Function to initialize mocks
		want      *domain.Todo
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				ctx:    context.Background(),
				userId: 1,
				id:     1,
			},
			wantErr: false,
			want: &domain.Todo{
				ID:        1,
				UserID:    1,
				Title:     "Test Todo",
				Done:      false,
				CreatedAt: now,
			},

			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Set up the expected behavior of the mock store
				// When Get is called with the given context and id, return a predefined todo
				store.On("Get", ta.ctx, ta.id).Return(&domain.Todo{
					ID:        ta.id,
					UserID:    ta.userId,
					Title:     "Test Todo",
					Done:      false,
					CreatedAt: now,
				}, nil).Once()

				s.Store = store
			},
		},
		{
			name:    "not found",
			fields:  fields{},
			wantErr: true,
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			want: nil,
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Simulate not found error from the store
				store.On("Get", ta.ctx, ta.id).Return(nil, errors.New("not found")).Once()

				s.Store = store
			},
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.GetTodo(tc.args.ctx, tc.args.userId, tc.args.id)

			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantErr, err != nil)
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
		ctx      context.Context
		userId   int64
		id       int64
		title    string
		done     bool
		priority int64
	}

	fixedTime := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	// Define the test cases

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoService) // Function to initialize mocks
		want      *domain.Todo
	}{
		{
			name:    "success",
			fields:  fields{},
			wantErr: false,
			args: args{
				ctx:      context.Background(),
				userId:   1,
				id:       1,
				title:    "Updated Todo",
				done:     true,
				priority: 3,
			},
			want: &domain.Todo{
				UserID:    1,
				ID:        1,
				Title:     "Updated Todo",
				Done:      true,
				Priority:  3,
				CreatedAt: fixedTime,
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				store.On("Get", ta.ctx, ta.id).Return(&domain.Todo{
					ID:        ta.id,
					UserID:    ta.userId,
					Title:     "Test Todo",
					Done:      false,
					CreatedAt: fixedTime,
				}, nil).Once()

				// When Update is called with the given context, id, title, and done status, return a predefined todo
				store.On("Update", ta.ctx, ta.id, ta.title, ta.done, ta.priority).Return(&domain.Todo{
					UserID:    ta.userId,
					ID:        ta.id,
					Title:     ta.title,
					Done:      ta.done,
					Priority:  ta.priority,
					CreatedAt: fixedTime,
				}, nil).Once()

				s.Store = store
			},
		},
		{
			name:    "not found",
			fields:  fields{},
			wantErr: true,
			args: args{
				ctx:   context.Background(),
				id:    999,
				title: "Updated Todo",
				done:  true,
			},
			want: nil,
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				store.On("Get", ta.ctx, ta.id).Return(&domain.Todo{
					ID:     ta.id,
					UserID: ta.userId,
					Title:  "Test Todo",
					Done:   false,
				}, nil).Once()

				store.On("Update", ta.ctx, ta.id, ta.title, ta.done, ta.priority).Return((*domain.Todo)(nil), errors.New("not found")).Once()

				s.Store = store
			},
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.UpdateTodo(tc.args.ctx, tc.args.userId, tc.args.id, tc.args.title, tc.args.done, tc.args.priority)

			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantErr, err != nil)
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
		ctx    context.Context
		userId int64
		id     int64
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
			name:    "success",
			fields:  fields{},
			wantErr: false,
			args: args{
				ctx:    context.Background(),
				userId: 1,
				id:     1,
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				store.On("Get", ta.ctx, ta.id).Return(&domain.Todo{
					ID:     ta.id,
					UserID: ta.userId,
					Title:  "Test Todo",
					Done:   false,
				}, nil).Once()

				// When Delete is called with the given context and id, return nil (no error)
				store.On("Delete", ta.ctx, ta.id).Return(nil).Once()

				s.Store = store
			},
		},
		{
			name:    "not found",
			fields:  fields{},
			wantErr: true,
			args: args{
				ctx: context.Background(),
				id:  999,
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoService) {
				store := mocks.NewTodoStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				store.On("Get", ta.ctx, ta.id).Return(&domain.Todo{
					ID:     ta.id,
					UserID: ta.userId,
					Title:  "Test Todo",
					Done:   false,
				}, nil).Once()

				store.On("Delete", ta.ctx, ta.id).Return(errors.New("not found")).Once()

				s.Store = store
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

			err := s.DeleteTodo(tt.args.ctx, tt.args.userId, tt.args.id)

			require.Equal(t, tt.wantErr, err != nil)
		})
	}
}
