package user

import (
	"github.com/GoncharovMikhail/go-sql/pkg/service/user"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Us     *user.UserService
	Engine *gin.Engine
}

const (
	userPath = "user/"
)
