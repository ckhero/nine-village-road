/**
 *@Description
 *@ClassName user_scenic
 *@Date 2021/5/14 下午4:26
 *@Author ckhero
 */

package domain

import (
	"context"
	"github.com/ckhero/go-common/db/mysql"
	"nine-village-road/pkg/constant"
)

// UserScenic [...]
type UserScenic struct {
	mysql.BaseEntity
	UserScenicId uint64  `gorm:"column:user_scenic_id;type:bigint(11);not null"`
	Scenic       string `gorm:"primaryKey;column:scenic;type:varchar(32);not null"` // 景点
	UserId       uint64  `gorm:"primaryKey;column:user_id;type:bigint(20);not null"`
	OpenId       string `gorm:"index:idx_openid;column:open_id;type:varchar(128);not null"`
	Status       string `gorm:"column:status;type:varchar(16);not null"` // [VALID;INVALID]
}

func (*UserScenic) TableName() string {
	return "user_scenic"
}

func (u *UserScenic) IsValid() bool {
	return u.Status == constant.UserScenicStatusValid
}

type UserScenicRepo interface {
	GetSpecialScenicByUserId(ctx context.Context, userId uint64, scenic string) (*UserScenic, error)
	CreateUserScenic(ctx context.Context, data *UserScenic) (*UserScenic, error)
	ListUserScenic(ctx context.Context, userId uint64) ([]*UserScenic, error)
}

type UserScenicUsecase interface {
	Scan(ctx context.Context, user *User, scenic string) (*UserScenic, error)
	ListUserScenic(ctx context.Context, userId uint64) ([]*UserScenic, error)
	CheckAllScenicScaned(ctx context.Context, userId uint64) error
}
