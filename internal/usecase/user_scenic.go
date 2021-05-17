/**
 *@Description
 *@ClassName user_scenic
 *@Date 2021/5/15 下午7:37
 *@Author ckhero
 */

package usecase

import (
	"context"
	"github.com/ckhero/go-common/errors"
	"github.com/ckhero/go-common/logger"
	"github.com/ckhero/go-common/util/uuid"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
)

type userScenicUsecase struct {
	userScenicRepo domain.UserScenicRepo
}

func NewUserScenicUsecase(userScenicRepo domain.UserScenicRepo) domain.UserScenicUsecase {
	return &userScenicUsecase{
		userScenicRepo: userScenicRepo,
	}
}

func(s *userScenicUsecase) Scan(ctx context.Context, user *domain.User, scenic string) (*domain.UserScenic, error) {
	logger.GetLogger(ctx).Infof("scan userId %d scenic %s", user.UserId, scenic)
	data := &domain.UserScenic{
		UserScenicId: uuid.GenUUID(),
		Scenic:       scenic,
		UserId:       user.UserId,
		OpenId:       user.OpenId,
		Status:       constant.UserScenicStatusValid,
	}

	return s.userScenicRepo.CreateUserScenic(ctx, data)
}

func(s *userScenicUsecase) ListUserScenic(ctx context.Context, userId uint64) ([]*domain.UserScenic, error) {
	return s.userScenicRepo.ListUserScenic(ctx, userId)
}

func(s *userScenicUsecase) CheckAllScenicScaned(ctx context.Context, userId uint64) error {
	list, err := s.userScenicRepo.ListUserScenic(ctx, userId)
	if err != nil {
		return err
	}
	alreadyScanedScenic := []string{}
	for _, userScenic := range list {
		if userScenic.IsValid() {
			alreadyScanedScenic = append(alreadyScanedScenic, userScenic.Scenic)
		}
	}
	diff, _ := funk.DifferenceString(constant.AllScenic, alreadyScanedScenic)
	if len(diff) == 0 {
		return nil
	}
	return errors.Newf(codes.Unknown, "user_scenic", "尚未完成打卡", "尚未完成打卡的列表 %+v", diff)
}
