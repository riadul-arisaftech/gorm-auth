package main

import (
	"fmt"
	"log"
	"net/http"

	"codemead.com/gorm_auth/controllers"
	"codemead.com/gorm_auth/initializers"
	"codemead.com/gorm_auth/models"
	"codemead.com/gorm_auth/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("could not load env variables:", err)
	}

	initializers.ConnectDB(&config)

	err = initializers.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("could not migrate:", err)
	}
	fmt.Println("Migration complete")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8080", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))

}
