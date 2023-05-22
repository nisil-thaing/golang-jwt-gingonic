package routes

import (
  "github.com/gin-gonic/gin"

  controllers "jwt-project/golang-gingonic/controllers"
  middleware "jwt-project/golang-gingonic/middleware"
)

func UserRoutes(router *gin.RouterGroup) {
  usersRouter := router.Group("/users", middleware.Authenticate)

  usersRouter.GET("/", controllers.FetchUsers)
  usersRouter.GET("/:user_id", controllers.FetchUserById)
}
