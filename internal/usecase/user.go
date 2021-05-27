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
	"google.golang.org/grpc/codes"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
)

type userUsecase struct {
	userRepo      domain.UserRepo
	redPacketRepo domain.UserRedPacketRepo
	wxRepo        domain.WeixinRepo
}

func NewUserUsecase(userRepo domain.UserRepo, wxRepo domain.WeixinRepo, redPacketRepo domain.UserRedPacketRepo) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo, wxRepo: wxRepo, redPacketRepo: redPacketRepo}
}

func (u *userUsecase) Login(ctx context.Context, code string) (*domain.User, error) {
	code2Session, err := u.wxRepo.Code2Session(ctx, code)
	if err != nil {

		return nil, errors.Newf(codes.Unauthenticated, "user", "code失效", "%v", err)
	}
	user, err := u.userRepo.FirstOrCreate(ctx, &domain.User{
		OpenId:     code2Session.OpenId,
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

func (u *userUsecase) CheckUserIllegal(ctx context.Context, openId string) (*domain.User, error) {
	user, err := u.userRepo.GetByOpenId(ctx, openId)

	if err != nil {
		return nil, err
	}

	if user.IsRecved() {
		return nil, errors.Newf(codes.Unknown, "user", "已经领取成功，无法重复领取", "%d", user.UserId)
	}
	return user, nil
}

func (u *userUsecase) GetByOpenId(ctx context.Context, openId string) (*domain.User, error) {
	return u.userRepo.GetByOpenId(ctx, openId)
}

func (u *userUsecase) ListRedPacket(ctx context.Context, userId uint64) ([]*domain.UserRedPacket, error) {
	res, err := u.redPacketRepo.GetRedPacketByStatus(ctx, userId, constant.UserRedPacketSucc)
	if errors.IsNotFound(err) {
		return []*domain.UserRedPacket{}, nil
	}
	if err != nil {
		return nil, err
	}
	return []*domain.UserRedPacket{res}, nil
}
