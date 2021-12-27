package user

import (
	"fmt"
	webConsts "github.com/GoncharovMikhail/go-sql/const"
	model "github.com/GoncharovMikhail/go-sql/model/user"
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
