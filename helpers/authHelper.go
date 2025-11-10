package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(ctx *gin.Context, role string) (err error) {
	return err
}

func MatchUserTypeToUid(ctx *gin.Context, userId string) (err error) {
	userType := ctx.GetString("user_type")
	uid := ctx.GetString("uid")

	if userType == "USER" && uid != userId {
		err := errors.New("Unauthorised to access this resource")
		return err
	}
	err = CheckUserType(ctx, userType)
	return err
}
