/**
 *@Description
 *@ClassName user_scenic
 *@Date 2021/5/15 下午7:25
 *@Author ckhero
 */

package repo

import (
	"context"
	"github.com/ckhero/go-common/db"
	"github.com/ckhero/go-common/errors"
	"google.golang.org/grpc/codes"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
)

type userScenicRepo struct {
	database *db.Database
}

func NewUserScenicRepo(database *db.Database) domain.UserScenicRepo {
	return &userScenicRepo{database: database}
}

func (s *userScenicRepo) GetSpecialScenicByUserId(ctx context.Context, userId uint64, scenic string) (*domain.UserScenic, error) {
	return nil, nil
}

func (s *userScenicRepo) CreateUserScenic(ctx context.Context, data *domain.UserScenic) (*domain.UserScenic, error) {
	conn := s.database.RDB(ctx).Omit("user_scenic_id").Save(data)
	if conn.Error != nil {
		return nil, errors.Errorf(codes.Unknown, "user_senic", "扫码保存失败", "%+v", conn.Error)
	}
	return data, nil
}

func (s *userScenicRepo) ListUserScenic(ctx context.Context, userId uint64) ([]*domain.UserScenic, error) {
	list := []*domain.UserScenic{}
	conn := s.database.RDB(ctx).
		Where("user_id = ?", userId).
		Where("status = ?", constant.UserScenicStatusValid).
		Find(&list)
	if conn.Error != nil {
		return nil, errors.Errorf(codes.Unknown, "user_senic", "扫码记录保存失败", "%+v", conn.Error)
	}
	return list, nil
}
