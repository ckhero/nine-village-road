/**
 *@Description
 *@ClassName user_red_packet
 *@Date 2021/5/14 下午4:31
 *@Author ckhero
 */

package repo

import (
	"context"
	"github.com/ckhero/go-common/db"
	"gorm.io/gorm"
	"nine-village-road/internal/domain"
)

type userRedPacketRepo struct {
	database *db.Database
}

func CreateTx(ctx context.Context, packet *domain.UserRedPacket) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		txConn := tx.Create(packet)
		return txConn.Error
	}
}