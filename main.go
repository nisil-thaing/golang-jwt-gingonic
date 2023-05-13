package main

import (
  "os"
  "net/http"
  "github.com/gin-gonic/gin"
  routes "jwt-project/golang-gingonic/routes"
)

var DEFAULT_PORT = "8080"

func welcomeController(c *gin.Context) {
  c.JSON(http.StatusOK, gin.H{"message": "Welcome to Golang JWT API with Gin Gonic!"})
}

func main() {
  port := os.Getenv("PORT")
  
  if port == "" {
    port = DEFAULT_PORT
  }

  router := gin.New()
  router.Use(gin.Logger())

  router.GET("/", welcomeController) 

  publicRouter := router.Group("/api")

  publicRouter.GET("/", welcomeController) 

  routes.AuthRoutes(publicRouter)
  routes.UserRoutes(publicRouter)

  router.Run(":" + port)
}
