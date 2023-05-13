package helpers

import (
  "errors"
  "github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
  userType := c.GetString("user_type")

  if userType != role {
    err := errors.New("You are unauthorized to access this resource!")
    return err
  }

  return nil
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
  userType := c.GetString("user_type")
  uid := c.GetString("uid")

  if userType != "ADMIN" && uid != userId {
    err := errors.New("You are unauthorized to access this resource!")
    return err
  }

  err = CheckUserType(c, userType)
  return err
}
