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
	"nine-village-road/internal/domain"
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
