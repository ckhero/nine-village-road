package router

import (
	"github.com/ckhero/go-common/config"
	gin2 "github.com/ckhero/go-common/gin"
	"github.com/gin-gonic/gin"
)

func registerWechatRouter(engine *gin.Engine) {
	handler, _, _ := newUserService()
	group := engine.Group("/api/v1/user")
	{
		// denglu
		group.GET("/login", handler.Login)
		group.Use(
			gin2.UserJwtAuthMiddleware(config.GetAuthCfg().SecretKey),
		)
		group.GET("/sendAppletRed", handler.SendAppletRed)
		group.GET("/walletTransfer", handler.WalletTransfer)
		group.GET("/scan", handler.Scan)
		group.GET("/listScenic", handler.ListScenic)
	}
}

