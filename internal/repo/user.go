/**
 *@Description
 *@ClassName user
 *@Date 2021/5/12 下午4:57
 *@Author ckhero
 */

package repo

import (
	"context"
	"github.com/ckhero/go-common/db"
	"github.com/ckhero/go-common/errors"
	xerrors "github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
)

type userRepo struct {
	database *db.Database
}

func NewUserRepo(database *db.Database) domain.UserRepo {
	return &userRepo{database: database}
}

func (u *userRepo) GetByOpenId(ctx context.Context, openId string) (*domain.User, error) {
	userInfo := domain.User{}
	conn := u.database.RDB(ctx).Model(&domain.User{}).Where("open_id = ?", openId)
	conn.First(&userInfo)
	if conn.RowsAffected == 0 {
		return nil, errors.NotFound("user", "用户不存在", openId)
	}
	return &userInfo, xerrors.Wrap(conn.Error, "find user by open id fail")
}

func (u *userRepo) FirstOrCreate(ctx context.Context, user *domain.User) (*domain.User, error) {
	conn := u.database.RDB(ctx).FirstOrCreate(user, domain.User{OpenId: user.OpenId})
	return user, xerrors.Wrap(conn.Error, "first or create user by open id fail")
}

func (u *userRepo) UpdateRecvStatusTx(ctx context.Context, userId uint64, oldRecvStatus, recvStatus string)  func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		if oldRecvStatus == constant.UserRecvStatusRecved {
			return errors.Newf(codes.Unknown, "user", "更新用户领取状态失败", "[userId] %d [oldRecvStatus] %s [newRecvStatus] %s", userId, oldRecvStatus, recvStatus)
		}

		//if oldRecvStatus == constant.UserRecvStatusRecving {
		//	return errors.Newf(codes.Unknown, "user", "请勿重复领取", "[userId] %d [oldRecvStatus] %s [newRecvStatus] %s", userId, oldRecvStatus, recvStatus)
		//}

		conn := u.database.RDB(ctx).Model(domain.User{}).
			Where("user_id = ? and recv_status = ?", userId, oldRecvStatus).
			Update("recv_status", recvStatus)

		return conn.Error
	}
}