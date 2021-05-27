/**
 *@Description
 *@ClassName red_packet_limit
 *@Date 2021/5/18 上午11:02
 *@Author ckhero
 */

package domain

import (
	"context"
	"github.com/ckhero/go-common/db/mysql"
	"gorm.io/gorm"
)

// RedPacketLimit [...]
type RedPacketLimit struct {
	mysql.BaseEntity
	RedPacketLimitId uint64 `gorm:"primaryKey;column:red_packet_limit_id;type:int(11) unsigned;not null"`
	LimitType        string `gorm:"column:limit_type;type:varchar(32);not null"`
	Amount           uint64 `gorm:"column:amount;type:bigint(20);not null"`
	LeftAmount       uint64 `gorm:"column:left_amount;type:bigint(20);not null"`
	RecvNum          uint64 `gorm:"column:recv_num;type:bigint(20);not null"`
	LeftRecvNum      uint64 `gorm:"column:left_recv_num;type:bigint(20);not null"`
	StartDate        string `gorm:"column:start_date;type:varchar(32);not null"`
	EndDate          string `gorm:"column:end_date;type:varchar(32);not null"`
}

func (*RedPacketLimit) TableName() string {
	return "red_packet_limit"
}

type RedPacketLimitRepo interface {
	UpdateTx(ctx context.Context, amount, num int) func(tx *gorm.DB) error
}
