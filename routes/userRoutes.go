package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/kishan8005/Golang-JWT-Authentication/controllers"
	middleware "github.com/kishan8005/Golang-JWT-Authentication/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}