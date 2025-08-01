package router

import (
	"github.com/lin-snow/ech0/internal/di"
)

// setupTodoRoutes 设置待办事项路由
func setupTodoRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	// Public

	// Auth
	appRouterGroup.AuthRouterGroup.GET("/todo", h.TodoHandler.GetTodoList())
	appRouterGroup.AuthRouterGroup.POST("/todo", h.TodoHandler.AddTodo())
	appRouterGroup.AuthRouterGroup.PUT("/todo/:id", h.TodoHandler.UpdateTodo())
	appRouterGroup.AuthRouterGroup.DELETE("/todo/:id", h.TodoHandler.DeleteTodo())
}
