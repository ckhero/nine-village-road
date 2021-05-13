package router

import (
	"github.com/gin-gonic/gin"
)

func registerWechatRouter(engine *gin.Engine) {
	handler, _, _ := newUserService()
	group := engine.Group("/api/v1/user")
	{
		// denglu
		group.GET("/login", handler.Login)
	}
}

