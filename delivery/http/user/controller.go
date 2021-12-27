package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	webConsts "sql/const"
	model "sql/model/user"
	"sql/pkg/service/user"
)

type UserController struct {
	Us     *user.UserService
	Engine *gin.Engine
}

const (
	userPath = "user/"
)

func (uc *UserController) init() {
	uc.
		Engine.
		PUT(
			webConsts.ApiPath+userPath+"save",
			func(ctx *gin.Context) {
				var request *model.UserSaveRequest
				err := ctx.BindJSON(request)
				if err != nil {
					panic(err)
				}
				save, err := (*uc.Us).Save(request)
				if err != nil {
					panic(err)
				}
				fmt.Println(save)
			},
		)
}
