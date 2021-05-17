/**
 *@Description
 *@ClassName user
 *@Date 2021/5/12 下午5:02
 *@Author ckhero
 */

package usecase

import (
	"context"
	"github.com/ckhero/go-common/auth"
	"github.com/ckhero/go-common/config"
	"github.com/ckhero/go-common/errors"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
)

type userUsecase struct {
	userRepo domain.UserRepo
	wxRepo   domain.WeixinRepo
}

func NewUserUsecase(userRepo domain.UserRepo, wxRepo domain.WeixinRepo) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo, wxRepo: wxRepo}
}

func (u *userUsecase) Login(ctx context.Context, code string) (*domain.User, error) {
	code2Session, err := u.wxRepo.Code2Session(ctx, code)
	if err != nil {

		return nil, err
	}
	user, err := u.userRepo.FirstOrCreate(ctx, &domain.User{
		OpenId: code2Session.OpenId,
		RecvStatus: constant.UserRecvStatusInit,
	})
	if err != nil {
		return nil, err
	}
	user.Token, _, err = auth.NewUserJwtToken(user.UserId, map[string]interface{}{
		"userId": user.UserId,
		"openId": user.OpenId,
	}, config.GetAuthCfg().SecretKey)

	return user, err
}

var whiteOpenIdList = map[string]struct{}{
	"om-Po5PJsl3_gkeX-KfL3nPFqOuE" : {},
	"om-Po5B0EtZ1Io6ouz6i2ZbZsnaQ" : {},
	"om-Po5FpkMGCD2EHGPgrzya5Rhyk" : {},
	"om-Po5PV5GjUGH-5Mn40YrGyhWzE" : {},
	"om-Po5Km1HK4vunZCYrE2yPhbRM4" : {},
	"om-Po5F11dYyhNLSZaS4k-UJ2teo" : {},
	"om-Po5O7syEkF-ncaN4FEs72FhwY" : {},
	"om-Po5FTg7RKP0YR-4eoy1eSwT2c" : {},
}


func (u *userUsecase) CheckUserIllegal(ctx context.Context, openId string) (*domain.User, error) {
	user, err := u.userRepo.GetByOpenId(ctx, openId)

	if err != nil {
		return nil, err
	}

	if _, ok := whiteOpenIdList[user.OpenId]; !ok {
		return nil, errors.NotFound("user", "非法用户", "")
	}

	return user, nil
}


func (u *userUsecase) GetByOpenId(ctx context.Context, openId string) (*domain.User, error) {
	return u.userRepo.GetByOpenId(ctx, openId)
}