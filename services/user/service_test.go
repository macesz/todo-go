package user

import (
	"context"
	"errors"
	"testing"

	"github.com/macesz/todo-go/domain"
	"github.com/macesz/todo-go/services/user/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	type fields struct {
		Store *mocks.UserStore
	}

	type args struct {
		ctx      context.Context
		name     string
		email    string
		password string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		want      *domain.User
		initMocks func(tt *testing.T, ta *args, s *UserService)
	}{
		{
			name:   "Success",
			fields: fields{},
			args: args{
				ctx:      context.Background(),
				name:     "Test User",
				email:    "test@example.com",
				password: "password",
			},
			wantErr: false,
			want: &domain.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			},
			initMocks: func(tt *testing.T, ta *args, s *UserService) {
				store := mocks.NewUserStore(tt)

				userMatcher := mock.MatchedBy(func(user *domain.User) bool {
					return user.Name == ta.name && user.Email == ta.email && user.Password == ta.password
				})

				store.On("CreateUser", ta.ctx, userMatcher).Return(&domain.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "password",
				}, nil).Once()

				s.UserStore = store
			},
		},
		{
			name:   "Error",
			fields: fields{},
			args: args{
				ctx:      context.Background(),
				name:     "Test User",
				email:    "test@example.com",
				password: "password",
			},
			wantErr: true,
			want:    nil,
			initMocks: func(tt *testing.T, ta *args, s *UserService) {
				store := mocks.NewUserStore(tt)

				store.On("CreateUser", ta.ctx, mock.Anything).Return(nil, errors.New("error")).Once()

				s.UserStore = store
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := &UserService{
				UserStore: mocks.NewUserStore(t),
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.CreateUser(tc.args.ctx, tc.args.name, tc.args.email, tc.args.password)

			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestGetTodo(t *testing.T) {
	t.Parallel()

	// Define the fields of the TodoService struct
	type fields struct {
		Store *mocks.UserStore
	}

	// Define the arguments for the GetTodo method
	type args struct {
		ctx context.Context
		id  int64
	}

	// Define the test cases
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *UserService) // Function to initialize mocks
		want      *domain.User
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			wantErr: false,
			want: &domain.User{
				ID:       1,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			},

			initMocks: func(tt *testing.T, ta *args, s *UserService) {
				store := mocks.NewUserStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Set up the expected behavior of the mock store
				// When Get is called with the given context and id, return a predefined todo
				store.On("GetUser", ta.ctx, ta.id).Return(&domain.User{
					ID:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "password",
				}, nil).Once()

				s.UserStore = store
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
			initMocks: func(tt *testing.T, ta *args, s *UserService) {
				store := mocks.NewUserStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Simulate not found error from the store

				store.On("GetUser", ta.ctx, ta.id).
					Return(nil, errors.
						New("not found")).Once()

				s.UserStore = store
			},
		},
	}

	for _, tc := range tests {
		tt := tc
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &UserService{
				UserStore: mocks.NewUserStore(t),
			}

			tc.initMocks(t, &tc.args, s)

			got, err := s.GetUser(tc.args.ctx, tc.args.id)

			require.Equal(t, tc.want, got)
			require.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestDeleteUser(t *testing.T) {

	t.Parallel()

	// Define the fields of the TodoService struct
	type fields struct {
		Store *mocks.UserStore
	}

	// Define the arguments for the DeleteTodo method
	// This allows us to pass different contexts and ids for each test case
	type args struct {
		ctx context.Context
		id  int64
	}

	// Define the test cases

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		initMocks func(tt *testing.T, ta *args, s *UserService) // Function to initialize mocks
	}{
		{
			name:    "success",
			fields:  fields{},
			wantErr: false,
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			initMocks: func(tt *testing.T, ta *args, s *UserService) {
				store := mocks.NewUserStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				// Set up the expected behavior of the mock store
				// When Delete is called with the given context and id, return nil (no error)
				store.On("Delete", ta.ctx, ta.id).Return(nil).Once()

				s.UserStore = store
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
			initMocks: func(tt *testing.T, ta *args, s *UserService) {
				store := mocks.NewUserStore(tt)
				tt.Cleanup(func() { store.AssertExpectations(tt) })

				store.On("Delete", ta.ctx, ta.id).Return(errors.New("not found")).Once()

				s.UserStore = store
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

						s := &UserService{
							UserStore: mocks.NewUserStore(t),
						}

						tc.initMocks(t, &tc.args, s)

						err := s.DeleteUser(tc.args.ctx, tc.args.id)

						require.Equal(t, tc.wantErr, err != nil)

		})
	}
}

