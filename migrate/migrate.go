package main

import (
	"fmt"
	"log"

	"github.com/salassep/golang-todolist-api/initializers"
	"github.com/salassep/golang-todolist-api/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.Todo{}, &models.SubTodo{}, &models.User{})
	fmt.Println("Migration complete")
}
