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
	"gorm.io/gorm"
	"nine-village-road/pkg/constant"
)

// User [...]
type User struct {
	mysql.BaseEntity
	UserId     uint64 `gorm:"primaryKey;column:user_id;type:int(11) unsigned;not null"`
	OpenId     string `gorm:"unique;column:open_id;type:varchar(128);not null"`
	UnionId     string `gorm:"column:union_id;type:varchar(128);not null"`
	RecvStatus string `gorm:"column:recv_status;type:varchar(32);not null"` // [INIT UNRECV RECVING RECVED]
	Token      string `gorm:"-";json:"token"`
}

func (*User) TableName() string {
	return "user"
}

func (u *User) IsRecving() bool {
	return u.RecvStatus == constant.UserRecvStatusRecving
}

func (u *User) IsRecved() bool {
	return u.RecvStatus == constant.UserRecvStatusRecved
}

type UserRepo interface {
	GetByOpenId(ctx context.Context, openId string) (*User, error)
	FirstOrCreate(ctx context.Context, user *User) (*User, error)
	UpdateRecvStatusTx(ctx context.Context, userId uint64, oldRecvStatus, recvStatus string) func(tx *gorm.DB) error
}

type UserUsecase interface {
	Login(ctx context.Context, code string) (*User, error)
	GetByOpenId(ctx context.Context, openId string) (*User, error)
	CheckUserIllegal(ctx context.Context, openId string) (*User, error)
	ListRedPacket(ctx context.Context, userId uint64) ([]*UserRedPacket, error)
}
