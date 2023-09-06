package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/salassep/golang-todolist-api/models"
	"gorm.io/gorm"
)

type SubTodoController struct {
	DB *gorm.DB
}

func NewSubTodoController(DB *gorm.DB) SubTodoController {
	return SubTodoController{DB}
}

func (tc *SubTodoController) CreateSubTodo(ctx *gin.Context) {
	var payload *models.CreateSubTodoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var todo models.Todo
	result := tc.DB.First(&todo, "id = ?", payload.Todo)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No todo found with that id"})
		return
	}

	now := time.Now()
	newSubTodo := models.SubTodo{
		Content:   payload.Content,
		Todo:      payload.Todo,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result = tc.DB.Create(&newSubTodo)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": newSubTodo})
}
