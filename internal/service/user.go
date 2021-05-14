/**
 *@Description
 *@ClassName user
 *@Date 2021/5/12 下午5:04
 *@Author ckhero
 */

package service

import (
	"github.com/ckhero/go-common/format"
	gin2 "github.com/ckhero/go-common/gin"
	"github.com/ckhero/go-common/util/context"
	"github.com/ckhero/go-common/util/json"
	"github.com/gin-gonic/gin"
	"nine-village-road/api"
	"nine-village-road/internal/domain"
)

type UserService struct {
	uc domain.UserUsecase
	weixinUsecase domain.WeixinUsecase
}

func NewUserService(uc domain.UserUsecase, weixinUsecase domain.WeixinUsecase) *UserService {
	return &UserService{uc:uc, weixinUsecase: weixinUsecase}
}

func(u *UserService) Login(c *gin.Context) {
	ctx, _ := context.ContextWithSpan(c)
	code := c.Query("code")
	user, err := u.uc.Login(ctx, code)
	if err != nil {
		format.Fail(c, err)
		return
	}
	format.Success(c, api.LogigRsp{Token: user.Token})
}


func(u *UserService) SendAppletRed(c *gin.Context) {
	openId := gin2.GetOpenId(c)
	ctx, _ := context.ContextWithSpan(c)
	if err := u.uc.CheckUserIllegal(ctx, openId); err != nil {
		format.Fail(c, err)
		return
	}
	pasySign, err := u.weixinUsecase.SendAppletRed(ctx, openId)
	if err != nil {
		format.Fail(c, err)
		return
	}
	rsp := api.SendAppletRedRsp{}
	_ = json.DeepCopyPHP(pasySign, &rsp)
	format.Success(c, rsp)
}
