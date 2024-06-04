package logic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/senyu-up/toolbox/tool/su_logger"
)

var UserLogic = new(user)

type user struct {
}

func (u *user) UserLogin(ctx *gin.Context) (data interface{}, err error) {

	go func() {
		su_logger.Info(ctx, "111")
	}()
	su_logger.Error(ctx, fmt.Errorf("异常"), "链路追踪测试")
	return nil, nil
}
