package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kishan8005/Golang-JWT-Authentication/database"
	helper "github.com/kishan8005/Golang-JWT-Authentication/helpers"
	"github.com/kishan8005/Golang-JWT-Authentication/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword() {}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
	var valid = true
	var msg string

	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		valid = false
		msg = fmt.Sprintln("email or password is incorrect")
	}

	return valid, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while fetching email from database"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this phone number already exists"})
			return
		}

		user.Created_at = time.Now()
		user.Upadted_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.First_name, user.Last_name, user.Phone, user.User_id, user.User_type)
		user.Token = token
		user.Refresh_token = refreshToken

		resultInsertNumber, InsertErr := userCollection.InsertOne(ctx, user)
		if InsertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}
		c.JSON(http.StatusOK, resultInsertNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user, foundUser models.User
		defer cancel()

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordValid, msg := VerifyPassword(user.Password, foundUser.Password)
		if !passwordValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

	}
}

func GetUsers() {}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("user_id")

		if err := helper.MatchUserTypeToUid(ctx, userId); err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Access denied": err.Error()})
			return
		}

		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := userCollection.FindOne(c, bson.M{"user_id": userId}).Decode(&user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
		return
	}
}
