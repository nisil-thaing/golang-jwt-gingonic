package helpers

import (
  "context"
  "fmt"
  "log"
  "os"
  "time"
  jwt "github.com/golang-jwt/jwt/v5"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  database "jwt-project/golang-gingonic/database"
  models "jwt-project/golang-gingonic/models"
)

type SignedDetails struct {
  Email string
  FirstName string
  LastName string
  Uid string
  UserType string
  jwt.RegisteredClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var secretKey string = os.Getenv("SECRET_KEY")

func GenerateTokens(user models.User) (accessToken string, refreshToken string, err error) {
  accessClaims := SignedDetails{
    Email: user.Email,
    FirstName: user.FirstName,
    LastName: user.LastName,
    Uid: user.UserID,
    UserType: user.Type,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Duration(24) * time.Hour)),
    },
  }
  refreshClaims := SignedDetails{
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(time.Duration(168) * time.Hour)),
    },
  }

  signedAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secretKey))
  
  if err != nil {
    log.Panic(err)
    return "", "", err
  }

  signedRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))
  if err != nil {
    log.Panic(err)
    return "", "", err
  }

  return signedAccessToken, signedRefreshToken, nil
}

func UpdateUserTokens(userId string, accessToken string, refreshToken string) (err error) {
  var updatingData primitive.D
  ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)

  updatingData = append(updatingData, bson.E{"access_token", accessToken})
  updatingData = append(updatingData, bson.E{"refresh_token", refreshToken})

  currentTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
  updatingData = append(updatingData, bson.E{"updated_at", currentTime})

  upsert := false
  filter := bson.M{"user_id": userId}
  opt := options.UpdateOptions {
    Upsert: &upsert,
  }

  _, err = userCollection.UpdateOne(
    ctx,
    filter,
    bson.D{{"$set", updatingData}},
    &opt,
  )
  defer cancel()

  if err != nil {
    log.Panic(err)
    return err
  }

  return nil
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
  token, err := jwt.ParseWithClaims(
    signedToken,
    &SignedDetails{},
    func (token *jwt.Token) (interface{}, error) {
      return []byte(secretKey), nil
    },
  )

  if err != nil {
    return nil, err
  }

  claims := token.Claims.(*SignedDetails)

  if claims.ExpiresAt.Time.Before(time.Now().Local()) {
    err = fmt.Errorf("This token is expired!")
    return nil, err
  }

  return claims, nil
}
