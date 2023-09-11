package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/salassep/golang-todolist-api/controllers"
	"github.com/salassep/golang-todolist-api/initializers"
	"github.com/salassep/golang-todolist-api/routes"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	TodoController      controllers.TodoController
	TodoRouteController routes.TodoRouteController

	SubTodoController      controllers.SubTodoController
	SubTodoRouteController routes.SubTodoRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	TodoController = controllers.NewTodoController(initializers.DB)
	TodoRouteController = routes.NewTodoRouteController(TodoController)

	SubTodoController = controllers.NewSubTodoController(initializers.DB)
	SubTodoRouteController = routes.NewSubTodoRouteController(SubTodoController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Todolist API"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	TodoRouteController.TodoRoute(router)
	SubTodoRouteController.SubTodoRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
