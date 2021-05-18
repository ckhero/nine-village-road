/**
 *@Description
 *@ClassName wechat
 *@Date 2021/5/18 上午9:35
 *@Author ckhero
 */

package service

import (
	"github.com/ckhero/go-common/format"
	"github.com/ckhero/go-common/util/context"
	"github.com/gin-gonic/gin"
	"nine-village-road/api"
	"nine-village-road/internal/domain"
)

type WechatService struct {
	uc domain.WeixinUsecase
}

func NewWechatService(uc domain.WeixinUsecase) *WechatService {
	return &WechatService{uc:uc}
}

// https://www.toolnb.com/tools/base64ToImages.html
func(s *WechatService) QRCode(c *gin.Context) {
	scenic := c.Query("scenic")
	ctx, _ := context.ContextWithSpan(c)
	qrcode, err := s.uc.QRCode(ctx, scenic)
	if err != nil {
		format.Fail(c, err)
		return
	}
	format.Success(c, api.QRCodeRsp{QRCode: qrcode})
}
