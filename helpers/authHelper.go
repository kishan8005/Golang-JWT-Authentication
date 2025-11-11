package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(ctx *gin.Context, role string) error {
	userType := ctx.GetString("user_type")
	var err error
	if userType != role {
		err = errors.New("User is not authorized")
	}
	return err
}

func MatchUserTypeToUid(ctx *gin.Context, userId string) error {
	userType := ctx.GetString("user_type")
	uid := ctx.GetString("uid")
	var err error
	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	return err
}
