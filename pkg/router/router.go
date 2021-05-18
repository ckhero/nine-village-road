package router

import (
	"github.com/ckhero/go-common/web"
	"github.com/gin-gonic/gin"
)

func RegisterRouter() web.RegisterRouter {

	return func(engine *gin.Engine) {

		registerWechatRouter(engine)
		registerUserRouter(engine)

	}
}
