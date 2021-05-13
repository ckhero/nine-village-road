/**
 *@Description
 *@ClassName user
 *@Date 2021/4/29 下午4:45
 *@Author ckhero
 */

package domain

import (
	"context"
	"github.com/ckhero/go-common/db/mysql"
)

// User [...]
type User struct {
	mysql.BaseEntity
	UserId uint64    `gorm:"primaryKey;column:user_id;type:int(11) unsigned;not null"`
	OpenId string `gorm:"unique;column:open_id;type:varchar(128);not null"`
}

func (*User) TableName() string {
	return "user"
}


type UserRepo interface {
	GetByOpenId(ctx context.Context, openId string) (*User, error)
	FirstOrCreate(ctx context.Context, user *User) (*User, error)
}

type UserUsecase interface {
	Login(ctx context.Context, code string) (string, error)
}