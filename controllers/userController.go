package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	database "jwt-project/golang-gingonic/database"
	helpers "jwt-project/golang-gingonic/helpers"
	models "jwt-project/golang-gingonic/models"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var validate = validator.New()

func Register(c *gin.Context) {
	var registeringUser models.RegisteringUser
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)

	if err := c.BindJSON(&registeringUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    defer cancel()
    return
	}

	if err := validate.Struct(registeringUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
    return
	}

	numOfExistingEmails, err := userCollection.CountDocuments(ctx, bson.M{"email": registeringUser.Email})
	defer cancel()

	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred when we trying to validate your information!"})

		return
	}

	numOfExistingPhoneNumbers, err := userCollection.CountDocuments(ctx, bson.M{"phone_number": registeringUser.PhoneNumber})
	defer cancel()

	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred when we trying to validate your information!"})

		return
	}

	if numOfExistingEmails > 0 || numOfExistingPhoneNumbers > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Your email or phone number has been used before by someone, please double-check your information carefully!"})
		return
	}

	currentTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
  userId := primitive.NewObjectID()
  hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(registeringUser.Password), 14)

  user := models.User {
    Type: "USER",
    Email: registeringUser.Email,
	  PhoneNumber: registeringUser.PhoneNumber,
	  FirstName: registeringUser.FirstName,
	  LastName: registeringUser.LastName,
	  CreatedAt: currentTime,
	  UpdatedAt: currentTime,
    ID: userId,
    UserID: userId.Hex(),
    HashedPassword: string(hashedPassword),
  }

  accessToken, refreshToken, _ := helpers.GenerateTokens(user)

  user.AccessToken = accessToken
  user.RefreshToken = refreshToken

	_, err = userCollection.InsertOne(ctx, user)
	defer cancel()

	if err != nil {
		errorMessage := "Oops! We couldn't register new account by your provided information! Please try again later!"
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMessage})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}

func Login(c *gin.Context) {
	var user models.LoginUser
	var matchingUser models.User
	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    defer cancel()
		return
	}

	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&matchingUser)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Your email or password is invalid! Please double-check carefully!"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(matchingUser.HashedPassword), []byte(user.Password))
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Your email or password is invalid! Please double-check carefully!"})
		return
	}

	accessToken, refreshToken, _ := helpers.GenerateTokens(matchingUser)

	helpers.UpdateUserTokens(matchingUser.UserID, accessToken, refreshToken)

	err = userCollection.FindOne(ctx, bson.M{"user_id": matchingUser.UserID}).Decode(&matchingUser)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}

func FetchUsers(c *gin.Context) {
	if err := helpers.CheckUserType(c, "ADMIN"); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)

	recordsPerPage, err := strconv.Atoi(c.Query("recordsPerPage"))

	if err != nil || recordsPerPage < 1 {
		recordsPerPage = 10
	}

	page, err := strconv.Atoi(c.Query("page"))

	if err != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordsPerPage

  matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
  groupStage := bson.D{{Key: "$group", Value: bson.D{
    {Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
    {Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
    {Key: "data", Value: bson.D{
      {Key: "$push", Value: bson.D{
        {Key: "user_id", Value: "$user_id"},
        {Key: "type", Value: "$type"},
        {Key: "displaying_name", Value: bson.D{{Key: "$concat", Value: []interface{}{"$first_name", " ", "$last_name"}}}},
        {Key: "email", Value: "$email"},
        {Key: "phone_number", Value: "$phone_number"},
        {Key: "created_at", Value: "$created_at"},
        {Key: "updated_at", Value: "$updated_at"},
      }},
    }},
	}}}
  projectStage := bson.D{{Key: "$project", Value: bson.D{
    {Key: "_id", Value: 0},
    {Key: "total_count", Value: 1},
    {Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordsPerPage}}}},
	}}}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})

	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
		return
	}

	var allUsers []bson.M

	if err = result.All(ctx, &allUsers); err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Oops! Something went wrong!"})
		return
	}

	c.JSON(http.StatusOK, allUsers[0])
}

func FetchUserById(c *gin.Context) {
	userId := c.Param("user_id")

	if err := helpers.MatchUserTypeToUid(c, userId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

	defer cancel()

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Hmm maybe there are no users mapped with that ID!"})
		return
	}

  responseData := helpers.ToPublishableInfo(user)

	c.JSON(http.StatusOK, responseData)
}
