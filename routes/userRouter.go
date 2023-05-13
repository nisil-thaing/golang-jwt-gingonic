package routes

import (
  "github.com/gin-gonic/gin"

  controllers "jwt-project/golang-gingonic/controllers"
  middleware "jwt-project/golang-gingonic/middleware"
)

func UserRoutes(router *gin.RouterGroup) {
  router.Use(middleware.Authenticate)
  router.GET("/users", controllers.FetchUsers)
  router.GET("/users/:user_id", controllers.FetchUserById)
}
