package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/salassep/golang-todolist-api/models"
	"gorm.io/gorm"
)

type TodoController struct {
	DB *gorm.DB
}

func NewTodoController(DB *gorm.DB) TodoController {
	return TodoController{DB}
}

func (tc *TodoController) CreateTodo(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateTodoRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newTodo := models.Todo{
		Title:     payload.Title,
		User:      currentUser.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := tc.DB.Create(&newTodo)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": newTodo})
}

func (tc *TodoController) UpdateTodo(ctx *gin.Context) {
	todoId := ctx.Param("todoId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdateTodo
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}

	var updatedTodo models.Todo
	result := tc.DB.First(&updatedTodo, "id = ?", todoId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	now := time.Now()
	todoToUpdate := models.Todo{
		Title:     payload.Title,
		User:      currentUser.ID,
		CreatedAt: updatedTodo.CreatedAt,
		UpdatedAt: now,
	}

	tc.DB.Model(&updatedTodo).Updates(todoToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedTodo})
}

func (tc *TodoController) FindTodoById(ctx *gin.Context) {
	todoId := ctx.Param("todoId")

	var todo models.Todo
	result := tc.DB.First(&todo, "id = ?", todoId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No todo found with that id"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": todo})
}
