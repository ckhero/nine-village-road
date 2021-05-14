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
	"github.com/ckhero/go-common/logger"
	"github.com/pkg/errors"
	"nine-village-road/internal/domain"
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

		logger.GetLoggerWithErr(ctx, errors.Cause(err)).Info("fsdfasdfasf")
		return nil, err
	}
	user, err := u.userRepo.FirstOrCreate(ctx, &domain.User{
		OpenId: code2Session.OpenId,
	})
	if err != nil {
		return nil, err
	}
	user.Token, _, err = auth.NewUserJwtToken(user.UserId, map[string]interface{}{
		"userId":  user.UserId,
		"opendId": user.OpenId,
	}, config.GetAuthCfg().SecretKey)

	return user, err
}
