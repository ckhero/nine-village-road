/**
 *@Description
 *@ClassName user_red_packet
 *@Date 2021/5/14 下午4:26
 *@Author ckhero
 */

package domain

import (
	"context"
	"github.com/ckhero/go-common/db/mysql"
	"gorm.io/gorm"
)

// UserRedPacket [...]
type UserRedPacket struct {
	mysql.BaseEntity
	UserRedPacketId uint64 `gorm:"primaryKey;column:user_red_packet_id;type:bigint(11);not null"`
	UserId          uint64 `gorm:"index:idx_user_id;column:user_id;type:bigint(11);not null"`
	OpenId          string `gorm:"index:idx_open_id;column:open_id;type:varchar(128);not null"`
	TradeNo         uint64 `gorm:"index:idx_trade_no;column:trade_no;type:bigint(20);not null"`
	Amount          uint64 `gorm:"column:amount;type:bigint(20);not null"`
	Status          string `gorm:"column:status;type:varchar(16);not null"` // INIT FAIL SUCC
	Remark          string `gorm:"column:remark;type:varchar(1000);not null"`
	Package         string `gorm:"column:package;type:varchar(1000);not null"`
}

func (*UserRedPacket) TableName() string {
	return "user_red_packet"
}

type UserRedPacketRepo interface {
	CreateTx(ctx context.Context, packet *UserRedPacket) func(tx *gorm.DB) error
	ConfirmTx(ctx context.Context, packet *UserRedPacket) func(tx *gorm.DB) error
	WaitRecvTx(ctx context.Context, packet *UserRedPacket) func(tx *gorm.DB) error
	CancelTx(ctx context.Context, packet *UserRedPacket) func(tx *gorm.DB) error
	HandleRedPacket(ctx context.Context, txList ...func(tx *gorm.DB) error) error
	GetRedPacketByStatus(ctx context.Context, userId uint64, status string) (*UserRedPacket, error)
}
