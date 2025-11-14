package todolist

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/macesz/todo-go/domain"
)

func (s *TodoListService) List(ctx context.Context, userID int64) ([]*domain.TodoList, error) {
	todoLists, err := s.Store.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list todo lists: %w", err)
	}

	return todoLists, nil
}

func (s *TodoListService) GetListByID(ctx context.Context, userID int64, id int64) (*domain.TodoList, error) {
	todoList, err := s.Store.GetListByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrListNotFound
		}
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	if todoList.UserID != userID {
		return nil, domain.ErrListNotFound
	}

	return todoList, nil
}

func (s *TodoListService) Create(ctx context.Context, userID int64, title string, color string, labels []string) (*domain.TodoList, error) {
	if title == "" {
		title = "Title"
	}

	createdAt := time.Now()

	todolist := &domain.TodoList{
		UserID:    userID,
		Title:     title,
		Color:     color,
		Labels:    labels,
		CreatedAt: createdAt,
	}

	err := s.Store.Create(ctx, todolist)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo list: %w", err)
	}

	return todolist, err
}

func (s *TodoListService) Update(ctx context.Context, userID int64, id int64, title string, color string, labels []string) (*domain.TodoList, error) {
	_, err := s.GetListByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	updated, err := s.Store.Update(ctx, userID, title, color, labels)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrListNotFound
		}
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	return updated, nil
}

func (s *TodoListService) Delete(ctx context.Context, userID int64, id int64) error {
	if _, err := s.GetListByID(ctx, userID, id); err != nil {
		return err
	}

	err := s.Store.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrListNotFound
		}
		return fmt.Errorf("failed to delete list: %w", err)
	}
	return nil
}
