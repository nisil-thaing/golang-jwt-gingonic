package helpers

import (
  "fmt"

  models "jwt-project/golang-gingonic/models"
)

func ToPublishableInfo(user models.User) models.PublishableUserInfo {
  displayingName := ""

  if user.FirstName == "" {
    if user.LastName != "" {
      displayingName = user.LastName
    }
  } else {
    if user.LastName == "" {
      displayingName = user.FirstName
    } else {
      displayingName = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
    }
  }
 
  data := models.PublishableUserInfo {
    UserID: user.UserID,
    DisplayingName: displayingName,
    Email: user.Email,
    PhoneNumber: user.PhoneNumber,
    Type: user.Type,
    CreatedAt: user.CreatedAt,
    UpdatedAt: user.UpdatedAt,
  }

  return data
}
