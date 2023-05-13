package routes

import (
  "github.com/gin-gonic/gin"
  controllers "jwt-project/golang-gingonic/controllers"
)

func AuthRoutes(router *gin.RouterGroup) {
  router.POST("/user/register", controllers.Register)
  router.POST("/user/login", controllers.Login)
}
