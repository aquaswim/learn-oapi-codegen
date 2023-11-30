package modules

import (
	"context"
	"github.com/golobby/container/v3"
	"todo-codegen/internal/dto"
	errorCode "todo-codegen/internal/error_code"
	"todo-codegen/internal/repositories"
)

type TodoModule interface {
	CreateTodo(ctx context.Context, todoItem dto.TodoItem) (*dto.TodoItemWithId, error)
	FindTodo(ctx context.Context, param dto.TodoItemListParams) (*dto.TodoItemFindResponse, error)
	DeleteTodoById(ctx context.Context, id int) (*dto.TodoItemWithId, error)
	GetTodoById(ctx context.Context, id int) (*dto.TodoItemWithId, error)
	UpdateTodo(ctx context.Context, id int, body dto.TodoItem) (*dto.TodoItemWithId, error)
}

type todoModule struct {
	todoRepo repositories.TodoRepository `container:"type"`
}

func (t todoModule) FindTodo(ctx context.Context, param dto.TodoItemListParams) (*dto.TodoItemFindResponse, error) {
	todos, err := t.todoRepo.FindTodo(ctx, &param)
	if err != nil {
		return nil, err
	}

	total, err := t.todoRepo.CountTodo(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.TodoItemFindResponse{
		Meta: &dto.ListMeta{
			Limit: param.Limit,
			Skip:  param.Skip,
			Total: &total,
		},
		Result: todos,
	}, nil
}

func (t todoModule) CreateTodo(ctx context.Context, todoItem dto.TodoItem) (*dto.TodoItemWithId, error) {
	todo, err := t.todoRepo.CreateTodo(ctx, todoItem)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (t todoModule) DeleteTodoById(ctx context.Context, id int) (*dto.TodoItemWithId, error) {
	todo, err := t.GetTodoById(ctx, id)
	if err != nil {
		return nil, err
	}

	err = t.todoRepo.DeleteTodoById(ctx, id)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (t todoModule) GetTodoById(ctx context.Context, id int) (*dto.TodoItemWithId, error) {
	todo, err := t.todoRepo.GetTodoById(ctx, id)
	if err != nil {
		return nil, errorCode.NewError(errorCode.ErrCodeInternal, err.Error())
	}
	if todo == nil {
		return nil, errorCode.ErrorTodoNotFound
	}
	return todo, nil
}

func (t todoModule) UpdateTodo(ctx context.Context, id int, body dto.TodoItem) (*dto.TodoItemWithId, error) {
	updatedTodo, err := t.todoRepo.UpdateTodoById(ctx, id, body)
	if err != nil {
		return nil, err
	}
	return updatedTodo, nil
}

func init() {
	container.MustSingletonLazy(container.Global, func() TodoModule {
		var module todoModule
		container.MustFill(container.Global, &module)
		return &module
	})
}
