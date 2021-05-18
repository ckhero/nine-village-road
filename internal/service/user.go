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
	"sync"
)

type UserService struct {
	uc                domain.UserUsecase
	weixinUsecase     domain.WeixinUsecase
	userScenicUsecase domain.UserScenicUsecase
}

func NewUserService(uc domain.UserUsecase, weixinUsecase domain.WeixinUsecase, userScenicUsecase domain.UserScenicUsecase) *UserService {
	return &UserService{uc: uc, weixinUsecase: weixinUsecase, userScenicUsecase: userScenicUsecase}
}

func (u *UserService) Login(c *gin.Context) {
	ctx, _ := context.ContextWithSpan(c)
	code := c.Query("code")
	user, err := u.uc.Login(ctx, code)
	if err != nil {
		format.Fail(c, err)
		return
	}
	format.Success(c, api.LogigRsp{Token: user.Token})
}

func (u *UserService) SendAppletRed(c *gin.Context) {
	openId := gin2.GetOpenId(c)
	ctx, _ := context.ContextWithSpan(c)
	// 校验用户
	user, err := u.uc.CheckUserIllegal(ctx, openId)
	if err != nil {
		format.Fail(c, err)
		return
	}
	// TODO 校验是否可领取
	// TODO 领取
	pasySign, err := u.weixinUsecase.SendAppletRed(ctx, user)
	if err != nil {
		format.Fail(c, err)
		return
	}
	// TODO confirm or cancel
	rsp := api.SendAppletRedRsp{}
	_ = json.DeepCopyPHP(pasySign, &rsp)
	format.Success(c, rsp)
}

var mu sync.Mutex
func (u *UserService) WalletTransfer(c *gin.Context) {
	mu.Lock()
	openId := gin2.GetOpenId(c)
	ctx, _ := context.ContextWithSpan(c)
	// 校验用户
	user, err := u.uc.CheckUserIllegal(ctx, openId)
	if err != nil {
		format.Fail(c, err)
		return
	}
	// TODO 校验是否可领取
	if err := u.userScenicUsecase.CheckAllScenicScaned(ctx, user.UserId); err != nil {
		format.Fail(c, err)
		return
	}
	// TODO 领取
	err = u.weixinUsecase.WalletTransfer(ctx, user)
	if err != nil {
		format.Fail(c, err)
		return
	}
	// TODO confirm or cancel
	format.Success(c, api.SendAppletRedRsp{})
}

// 扫码
func (u *UserService) Scan(c *gin.Context) {
	openId := gin2.GetOpenId(c)
	ctx, _ := context.ContextWithSpan(c)
	scenic := c.Query("scenic")
	user, err := u.uc.CheckUserIllegal(ctx, openId)
	if err != nil {
		format.Fail(c, err)
		return
	}
	_, err = u.userScenicUsecase.Scan(ctx, user, scenic)
	if err != nil {
		format.Fail(c, err)
		return
	}

	format.Success(c, struct{}{})
}

// 查看打卡列表
func (u *UserService) ListScenic(c *gin.Context) {
	userId := gin2.GetUserId(c)
	ctx, _ := context.ContextWithSpan(c)
	list, err := u.userScenicUsecase.ListUserScenic(ctx, userId)
	if err != nil {
		format.Fail(c, err)
		return
	}
	rspList := []*api.UserScenic{}
	_ = json.DeepCopyPHP(list, &rspList)
	format.Success(c, api.ListUserScenicRsp{List: rspList})
}
