package controllers

import (
	"net/http"
	"strconv"
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
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": newTodo})
}

func (tc *TodoController) UpdateTodo(ctx *gin.Context) {
	todoId := ctx.Param("todoId")
	currentUser := ctx.MustGet("currentUser").(models.User)

	var payload *models.UpdateTodo
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updatedTodo models.Todo
	result := tc.DB.First(&updatedTodo, "id = ?", todoId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Todo with that title exists"})
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

func (tc *TodoController) FindTodos(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offSet := (intPage - 1) * intLimit

	var todos []models.Todo
	results := tc.DB.Limit(intLimit).Offset(offSet).Find(&todos)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": results.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(todos), "data": todos})
}

func (tc *TodoController) DeleteTodo(ctx *gin.Context) {
	todoId := ctx.Param("todoId")

	result := tc.DB.Delete(&models.Todo{}, "id = ?", todoId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
