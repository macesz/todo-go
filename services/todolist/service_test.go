package todolist

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/services/todolist/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var fixedTime = time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

func TestListTodos(t *testing.T) {
	t.Parallel()

	type fields struct {
		Store *mocks.TodoListStore
	}

	type args struct {
		ctx    context.Context
		userID int64
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoListService)
		want      []*domain.TodoList
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{ctx: context.Background()},
			want: []*domain.TodoList{
				{ID: 1, UserID: 1, Title: "Shopping", Color: "white", Labels: nil, CreatedAt: fixedTime, Items: nil},
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				store.On("List", ta.ctx, ta.userID).Return([]*domain.TodoList{
					{ID: 1, UserID: 1, Title: "Shopping", Color: "white", Labels: nil, CreatedAt: fixedTime, Items: nil},
				}, nil).Once()

				s.Store = store
			},
		}, {
			name:    "store error",
			fields:  fields{},
			args:    args{ctx: context.Background()},
			wantErr: true,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

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

			s := &TodoListService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.List(tc.args.ctx, tc.args.userID)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestCreateList(t *testing.T) {
	t.Parallel()

	type fields struct {
		Store *mocks.TodoListStore
	}

	type args struct {
		ctx    context.Context
		userId int64
		title  string
		color  string
		labels []string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *TodoListService) // Function to initialize mocks
		validate  func(*testing.T, *args, *domain.TodoList)
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{ctx: context.Background(), userId: 1, title: "Shopping", color: "white", labels: nil},
			validate: func(t *testing.T, ta *args, todoList *domain.TodoList) {
				require.Equal(t, ta.userId, todoList.UserID)
				require.Equal(t, ta.title, todoList.Title)
				require.Equal(t, ta.color, todoList.Color)
				require.Equal(t, ta.labels, todoList.Labels)
				require.NotZero(t, todoList.CreatedAt)
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Set up the expected behavior of the mock store, "ta." means test args
				// When Create is called with the given context and title, return a predefined todo
				store.On("Create", ta.ctx, mock.MatchedBy(
					func(todo *domain.TodoList) bool {
						return todo.UserID == ta.userId &&
							todo.Title == ta.title &&
							todo.Color == ta.color
					})).Run(func(args mock.Arguments) {
					// Simulate the store setting the ID
					todo := args.Get(1).(*domain.TodoList)
					todo.ID = 1
				}).Return(nil).Once()

				s.Store = store
			},
		}, {
			name:    "store error",
			fields:  fields{},
			args:    args{ctx: context.Background(), userId: 1, title: "Shopping", color: "white", labels: nil},
			wantErr: true,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Simulate an error from the store

				store.On("Create", ta.ctx, mock.MatchedBy(
					func(todoList *domain.TodoList) bool {
						return todoList.UserID == ta.userId &&
							todoList.Title == ta.title &&
							todoList.Color == ta.color &&
							todoList.Labels == nil
					})).Run(func(args mock.Arguments) {
					// Simulate the store setting the ID
					todo := args.Get(1).(*domain.TodoList)
					todo.ID = 1
				}).Return(fmt.Errorf("test error")).Once()
				s.Store = store
			},
			validate: func(t *testing.T, ta *args, todoList *domain.TodoList) {
				require.Nil(t, todoList)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoListService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.Create(tc.args.ctx, tc.args.userId, tc.args.title, tc.args.color, tc.args.labels)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tc.validate != nil {
				tc.validate(t, &tc.args, got)
			}
		})
	}

}

func TestGetListByID(t *testing.T) {
	t.Parallel()

	type fields struct {
		Store *mocks.TodoListStore
	}

	type args struct {
		ctx    context.Context
		userID int64
		id     int64
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantedErr error
		initMocks func(tt *testing.T, ta *args, s *TodoListService)
		want      *domain.TodoList
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{ctx: context.Background(), userID: 1, id: 1},
			want: &domain.TodoList{
				ID:        1,
				UserID:    1,
				Title:     "Shopping",
				Color:     "white",
				Labels:    nil,
				Items:     nil,
				CreatedAt: fixedTime,
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				store.On("GetListByID", ta.ctx, ta.id).Return(&domain.TodoList{
					ID:        1,
					UserID:    1,
					Title:     "Shopping",
					Color:     "white",
					Labels:    nil,
					Items:     nil,
					CreatedAt: fixedTime,
				}, nil).Once()

				s.Store = store
			},
		},
		{
			name:      "list not found - sql.ErrNoRows",
			fields:    fields{},
			args:      args{ctx: context.Background(), userID: 1, id: 999},
			wantErr:   true,
			wantedErr: domain.ErrListNotFound,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				store.On("GetListByID", ta.ctx, ta.id).Return(nil, sql.ErrNoRows).Once()

				s.Store = store
			},
		},
		{
			name:    "store error",
			fields:  fields{},
			args:    args{ctx: context.Background(), userID: 1, id: 1},
			wantErr: true,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				store.On("GetListByID", ta.ctx, ta.id).Return(nil, errors.New("database error")).Once()

				s.Store = store
			},
		},
		{
			name:      "list belongs to different user",
			fields:    fields{},
			args:      args{ctx: context.Background(), userID: 1, id: 2},
			wantErr:   true,
			wantedErr: domain.ErrListNotFound,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				// Store returns a list, but it belongs to userID 2, not 1
				store.On("GetListByID", ta.ctx, ta.id).Return(&domain.TodoList{
					ID:        2,
					UserID:    2, // Different user!
					Title:     "Someone else's list",
					Color:     "blue",
					Labels:    nil,
					Items:     nil,
					CreatedAt: fixedTime,
				}, nil).Once()

				s.Store = store
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoListService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.GetListByID(tc.args.ctx, tc.args.userID, tc.args.id)
			if tc.wantErr {
				require.Error(t, err)
				if tc.wantedErr != nil {
					require.ErrorIs(t, err, tc.wantedErr)
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	type fields struct {
		Store *mocks.TodoListStore
	}

	type args struct {
		ctx    context.Context
		userID int64
		id     int64
		title  string
		color  string
		labels []string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantedErr error
		initMocks func(tt *testing.T, ta *args, s *TodoListService)
		want      *domain.TodoList
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				ctx:    context.Background(),
				userID: 1,
				id:     1,
				title:  "Updated Shopping",
				color:  "blue",
				labels: []string{"urgent", "groceries"},
			},
			want: &domain.TodoList{
				ID:        1,
				UserID:    1,
				Title:     "Updated Shopping",
				Color:     "blue",
				Labels:    []string{"urgent", "groceries"},
				Items:     nil,
				CreatedAt: fixedTime,
			},
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				// Mock GetListByID - verify list exists and belongs to user
				store.On("GetListByID", ta.ctx, ta.id).Return(&domain.TodoList{
					ID:        1,
					UserID:    1,
					Title:     "Shopping",
					Color:     "white",
					Labels:    nil,
					Items:     nil,
					CreatedAt: fixedTime,
				}, nil).Once()

				// Mock Update
				store.On("Update", ta.ctx, ta.id, ta.title, ta.color, ta.labels).Return(&domain.TodoList{
					ID:        1,
					UserID:    1,
					Title:     "Updated Shopping",
					Color:     "blue",
					Labels:    []string{"urgent", "groceries"},
					Items:     nil,
					CreatedAt: fixedTime,
				}, nil).Once()

				s.Store = store
			},
		},
		{
			name:      "list not found",
			fields:    fields{},
			args:      args{ctx: context.Background(), userID: 1, id: 999, title: "Test", color: "red", labels: nil},
			wantErr:   true,
			wantedErr: domain.ErrListNotFound,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				// GetListByID returns not found
				store.On("GetListByID", ta.ctx, ta.id).Return(nil, sql.ErrNoRows).Once()

				s.Store = store
			},
		},
		{
			name:      "list belongs to different user",
			fields:    fields{},
			args:      args{ctx: context.Background(), userID: 1, id: 2, title: "Hacked", color: "red", labels: nil},
			wantErr:   true,
			wantedErr: domain.ErrListNotFound,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				// GetListByID returns a list belonging to user 2
				store.On("GetListByID", ta.ctx, ta.id).Return(&domain.TodoList{
					ID:        2,
					UserID:    2, // Different user!
					Title:     "Someone else's list",
					Color:     "blue",
					Labels:    nil,
					Items:     nil,
					CreatedAt: fixedTime,
				}, nil).Once()

				s.Store = store
			},
		},
		{
			name:      "store update returns sql.ErrNoRows",
			fields:    fields{},
			args:      args{ctx: context.Background(), userID: 1, id: 1, title: "Test", color: "red", labels: nil},
			wantErr:   true,
			wantedErr: domain.ErrListNotFound,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				// GetListByID succeeds
				store.On("GetListByID", ta.ctx, ta.id).Return(&domain.TodoList{
					ID:        1,
					UserID:    1,
					Title:     "Shopping",
					Color:     "white",
					Labels:    nil,
					Items:     nil,
					CreatedAt: fixedTime,
				}, nil).Once()

				// But Update fails with ErrNoRows (race condition scenario)
				store.On("Update", ta.ctx, ta.id, ta.title, ta.color, ta.labels).Return(nil, sql.ErrNoRows).Once()

				s.Store = store
			},
		},
		{
			name:    "store update error",
			fields:  fields{},
			args:    args{ctx: context.Background(), userID: 1, id: 1, title: "Test", color: "red", labels: nil},
			wantErr: true,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				// GetListByID succeeds
				store.On("GetListByID", ta.ctx, ta.id).Return(&domain.TodoList{
					ID:        1,
					UserID:    1,
					Title:     "Shopping",
					Color:     "white",
					Labels:    nil,
					Items:     nil,
					CreatedAt: fixedTime,
				}, nil).Once()

				// Update fails with generic error
				store.On("Update", ta.ctx, ta.id, ta.title, ta.color, ta.labels).Return(nil, errors.New("database error")).Once()

				s.Store = store
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoListService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.Update(tc.args.ctx, tc.args.userID, tc.args.id, tc.args.title, tc.args.color, tc.args.labels)
			if tc.wantErr {
				require.Error(t, err)
				if tc.wantedErr != nil {
					require.ErrorIs(t, err, tc.wantedErr)
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	type fields struct {
		Store *mocks.TodoListStore
	}

	type args struct {
		ctx    context.Context
		userID int64
		id     int64
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantedErr error
		initMocks func(tt *testing.T, ta *args, s *TodoListService)
	}{
		{
			name:   "success",
			fields: fields{},
			args:   args{ctx: context.Background(), userID: 1, id: 1},
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				// Mock GetListByID - verify list exists and belongs to user
				store.On("GetListByID", ta.ctx, ta.id).Return(&domain.TodoList{
					ID:        1,
					UserID:    1,
					Title:     "Shopping",
					Color:     "white",
					Labels:    nil,
					Items:     nil,
					CreatedAt: fixedTime,
				}, nil).Once()

				// Mock Delete
				store.On("Delete", ta.ctx, ta.id).Return(nil).Once()

				s.Store = store
			},
		},
		{
			name:      "list not found",
			fields:    fields{},
			args:      args{ctx: context.Background(), userID: 1, id: 999},
			wantErr:   true,
			wantedErr: domain.ErrListNotFound,
			initMocks: func(tt *testing.T, ta *args, s *TodoListService) {
				store := mocks.NewTodoListStore(tt)

				tt.Cleanup(func() {
					store.AssertExpectations(tt)
				})

				// GetListByID returns not found
				store.On("GetListByID", ta.ctx, ta.id).Return(nil, sql.ErrNoRows).Once()

				s.Store = store
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &TodoListService{
				Store: tc.fields.Store,
			}

			tc.initMocks(t, &tc.args, s)

			err := s.Delete(tc.args.ctx, tc.args.userID, tc.args.id)
			if tc.wantErr {
				require.Error(t, err)
				if tc.wantedErr != nil {
					require.ErrorIs(t, err, tc.wantedErr)
				}
				return
			}
			require.NoError(t, err)
		})
	}
}
