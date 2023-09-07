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

func (stc *SubTodoController) CreateSubTodo(ctx *gin.Context) {
	var payload *models.CreateSubTodoRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var todo models.Todo
	result := stc.DB.First(&todo, "id = ?", payload.Todo)
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

	result = stc.DB.Create(&newSubTodo)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": newSubTodo})
}

func (stc *SubTodoController) UpdateSubTodo(ctx *gin.Context) {
	subTodoId := ctx.Param("subTodoId")

	var payload *models.UpdateSubTodo
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var updatedSubTodo models.SubTodo
	result := stc.DB.First(&updatedSubTodo, "id = ?", subTodoId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No subtodo found with that id"})
		return
	}

	now := time.Now()
	subTodoToUpdate := models.UpdateSubTodo{
		Content:   payload.Content,
		Todo:      updatedSubTodo.Todo,
		CreatedAt: updatedSubTodo.CreatedAt,
		UpdatedAt: now,
	}

	stc.DB.Model(&updatedSubTodo).Updates(subTodoToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": subTodoToUpdate})

}

func (stc *SubTodoController) DeleteSubTodo(ctx *gin.Context) {
	subTodoId := ctx.Param("subTodoId")

	result := stc.DB.Delete(&models.SubTodo{}, "id = ?", subTodoId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
