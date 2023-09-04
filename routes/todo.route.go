package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/salassep/golang-todolist-api/controllers"
	"github.com/salassep/golang-todolist-api/middleware"
)

type TodoRouteController struct {
	todoController controllers.TodoController
}

func NewTodoRouteController(todoController controllers.TodoController) TodoRouteController {
	return TodoRouteController{todoController}
}

func (rc *TodoRouteController) TodoRoute(rg *gin.RouterGroup) {
	router := rg.Group("todo")
	router.Use(middleware.DeserializeUser())
	router.POST("/", rc.todoController.CreateTodo)
	router.GET("/", rc.todoController.FindTodos)
	router.PUT("/:todoId", rc.todoController.UpdateTodo)
	router.GET("/:todoId", rc.todoController.FindTodoById)
	router.DELETE("/:todoId", rc.todoController.DeleteTodo)
}
