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
	"nine-village-road/internal/domain"
)

type userUsecase struct {
	userRepo domain.UserRepo
	wxRepo   domain.WeixinRepo
}

func NewUserUsecase(userRepo domain.UserRepo, wxRepo domain.WeixinRepo) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo, wxRepo: wxRepo}
}

func (u *userUsecase) Login(ctx context.Context, code string) (string, error) {
	code2Session, err := u.wxRepo.Code2Session(ctx, code)
	if err != nil {
		return "", err
	}
	user, err := u.userRepo.FirstOrCreate(ctx, &domain.User{
		OpenId: code2Session.OpenId,
	})
	if err != nil {
		return "", err
	}
	token, _, err := auth.NewUserJwtToken(user.UserId, map[string]interface{}{
		"userId":  user.UserId,
		"opendId": user.OpenId,
	}, config.GetAuthCfg().SecretKey)

	return token, err
}
