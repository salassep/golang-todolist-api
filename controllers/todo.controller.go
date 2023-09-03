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
