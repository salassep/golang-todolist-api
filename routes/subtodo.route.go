package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/salassep/golang-todolist-api/controllers"
	"github.com/salassep/golang-todolist-api/middleware"
)

type SubTodoRouteController struct {
	subTodoController controllers.SubTodoController
}

func NewSubTodoRouteController(subTodoController controllers.SubTodoController) SubTodoRouteController {
	return SubTodoRouteController{subTodoController}
}

func (rc *SubTodoRouteController) SubTodoRoute(rg *gin.RouterGroup) {
	router := rg.Group("subtodo")
	router.Use(middleware.DeserializeUser())
	router.POST("/", rc.subTodoController.CreateSubTodo)
	router.PUT("/:subTodoId", rc.subTodoController.UpdateSubTodo)
	router.DELETE("/:subTodoId", rc.subTodoController.DeleteSubTodo)
}
