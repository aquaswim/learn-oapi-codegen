package api_server

import (
	"github.com/golobby/container/v3"
	"github.com/labstack/echo/v4"
	"net/http"
	"todo-codegen/internal/dto"
	"todo-codegen/internal/modules"
)

func New() ServerInterface {
	var server apiServer
	container.MustFill(container.Global, &server)
	return server
}

type apiServer struct {
	todoModule modules.TodoModule `container:"type"`
}

func (a apiServer) HealthCheck(ctx echo.Context) error {
	return ctx.JSON(200, dto.HealthStatus{Healthy: true})
}

func (a apiServer) TodoItemList(ctx echo.Context, params dto.TodoItemListParams) error {
	todo, err := a.todoModule.FindTodo(ctx.Request().Context(), params)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, todo)
}

func (a apiServer) TodoItemCreate(ctx echo.Context) error {
	var reqBody dto.TodoItem
	err := ctx.Bind(&reqBody)
	if err != nil {
		return err
	}
	todo, err := a.todoModule.CreateTodo(ctx.Request().Context(), reqBody)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, todo)
}

func (a apiServer) TodoItemDeleteById(ctx echo.Context, id int) error {
	deletedTodo, err := a.todoModule.DeleteTodoById(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, deletedTodo)
}

func (a apiServer) TodoItemGetById(ctx echo.Context, id int) error {
	todoItem, err := a.todoModule.GetTodoById(ctx.Request().Context(), id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, todoItem)
}

func (a apiServer) TodoItemUpdateById(ctx echo.Context, id int) error {
	var reqBody dto.TodoItem
	err := ctx.Bind(&reqBody)
	if err != nil {
		return err
	}
	todo, err := a.todoModule.UpdateTodo(ctx.Request().Context(), id, reqBody)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, todo)
}
