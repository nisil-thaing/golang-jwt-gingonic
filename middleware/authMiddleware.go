package middleware

import (
  "fmt"
  "strings"
  "net/http"
  "github.com/gin-gonic/gin"

  helpers "jwt-project/golang-gingonic/helpers"
)

func Authenticate(c *gin.Context) {
  bearerToken := c.Request.Header.Get("Authorization")

  if bearerToken == "" || strings.Contains(bearerToken, "Bearer ") == false {
    errMessage := fmt.Sprintf("You have no permission to access this content!")
    c.JSON(http.StatusUnauthorized, gin.H{"error": errMessage})
    c.Abort()
    return
  }

  clientToken := strings.ReplaceAll(bearerToken, "Bearer ", "")
  clientToken = strings.TrimSpace(clientToken)

  claims, err := helpers.ValidateToken(clientToken)

  if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    c.Abort()
    return
  }

  c.Set("email", claims.Email)
  c.Set("first_name", claims.FirstName)
  c.Set("last_name", claims.LastName)
  c.Set("uid", claims.Uid)
  c.Set("user_type", claims.UserType)

  c.Next()
}
