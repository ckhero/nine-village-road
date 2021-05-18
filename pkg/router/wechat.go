/**
 *@Description
 *@ClassName wechat
 *@Date 2021/5/18 上午9:34
 *@Author ckhero
 */

package router

import (
	"github.com/gin-gonic/gin"
)

func registerWechatRouter(engine *gin.Engine) {
	handler, _, _ := newWechatService()
	group := engine.Group("/api/v1/wechat")
	{
		// denglu
		group.GET("/qrCode", handler.QRCode)
	}
}
