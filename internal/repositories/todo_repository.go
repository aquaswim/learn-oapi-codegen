package repositories

import (
	"context"
	"github.com/golobby/container/v3"
	"sync"
	"todo-codegen/internal/dto"
	errorCode "todo-codegen/internal/error_code"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, item dto.TodoItem) (*dto.TodoItemWithId, error)
	FindTodo(ctx context.Context, param *dto.TodoItemListParams) (*[]dto.TodoItemWithId, error)
	CountTodo(ctx context.Context) (int, error)
	GetTodoById(ctx context.Context, id int) (*dto.TodoItemWithId, error)
	DeleteTodoById(ctx context.Context, id int) error
	UpdateTodoById(ctx context.Context, id int, data dto.TodoItem) (*dto.TodoItemWithId, error)
}

type todoRepo struct {
	mu    sync.Mutex
	data  []dto.TodoItemWithId
	curID int
}

func (t *todoRepo) CreateTodo(_ context.Context, item dto.TodoItem) (*dto.TodoItemWithId, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	todoItem := dto.TodoItemWithId{
		Description: item.Description,
		Id:          t.nextID(),
		Title:       item.Title,
	}
	t.data = append(t.data, todoItem)
	return &todoItem, nil
}

func (t *todoRepo) FindTodo(_ context.Context, param *dto.TodoItemListParams) (*[]dto.TodoItemWithId, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	length := len(t.data)
	start := min(*param.Skip, length)
	end := min(*param.Skip+*param.Limit, length)

	slicedData := t.data[start:end]
	return &slicedData, nil
}

func (t *todoRepo) CountTodo(_ context.Context) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return len(t.data), nil
}

func (t *todoRepo) GetTodoById(_ context.Context, id int) (*dto.TodoItemWithId, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, d := range t.data {
		if d.Id == id {
			return &d, nil
		}
	}
	// if not found send null item without error object
	return nil, nil
}

func (t *todoRepo) nextID() int {
	t.curID++
	return t.curID
}

func (t *todoRepo) DeleteTodoById(_ context.Context, id int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	index := -1
	for i, d := range t.data {
		if d.Id == id {
			index = i
		}
	}
	if index == -1 {
		return errorCode.ErrorTodoNotFound
	}
	// delete the item in list
	t.data = append(t.data[:index], t.data[index+1:]...)
	return nil
}

func (t *todoRepo) UpdateTodoById(_ context.Context, id int, data dto.TodoItem) (*dto.TodoItemWithId, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	for i, d := range t.data {
		if d.Id == id {
			t.data[i].Title = data.Title
			t.data[i].Description = data.Description
			return &t.data[i], nil
		}
	}
	return nil, errorCode.ErrorTodoNotFound
}

func init() {
	container.MustSingletonLazy(container.Global, func() TodoRepository {
		return &todoRepo{
			mu:    sync.Mutex{},
			data:  make([]dto.TodoItemWithId, 0, 100),
			curID: 0,
		}
	})
}
