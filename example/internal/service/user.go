package service

import (
	"github.com/gin-gonic/gin"
	"github.com/senyu-up/toolbox/example/internal/logic"
	"github.com/senyu-up/toolbox/tool/http/gin_server/controller"
)

var UserController = new(user)

type user struct {
	controller.BaseController
}

// UserLogin
// @tags 用户管理
// @summary 用户登录
// @param req body model.UserLoginParams true "json入参"
// @success 200 {object} model.UserLoginResp
// @router /user/login [post]
func (u *user) UserLogin(ctx *gin.Context) {
	u.Call_(ctx, nil, logic.UserLogic.UserLogin)
}
