/**
 *@Description
 *@ClassName user
 *@Date 2021/5/12 下午5:04
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
	token, err := u.uc.Login(ctx, code)
	if err != nil {
		format.Fail(c, err)
		return
	}
	format.Success(c, api.LogigRsp{Token: token})
}
