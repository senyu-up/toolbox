package index

import (
	"github.com/gin-gonic/gin"
	"github.com/senyu-up/toolbox/example/internal/service"
)

func RegisterRouterG(app *gin.Engine) {
	userGroup := app.Group("/user")
	{
		userGroup.POST("/login", service.UserController.UserLogin)
	}
}
