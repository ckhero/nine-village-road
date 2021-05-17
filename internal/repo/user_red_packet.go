/**
 *@Description
 *@ClassName user_red_packet
 *@Date 2021/5/14 下午4:31
 *@Author ckhero
 */

package repo

import (
	"context"
	"fmt"
	"github.com/ckhero/go-common/db"
	"github.com/ckhero/go-common/errors"
	"gorm.io/gorm"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
)

type userRedPacketRepo struct {
	database *db.Database
}

func NewUserRedPacketRepo(database *db.Database) domain.UserRedPacketRepo {
	return &userRedPacketRepo{database: database}
}

func (u *userRedPacketRepo) CreateTx(ctx context.Context, packet *domain.UserRedPacket) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		txConn := tx.Create(packet)
		return txConn.Error
	}
}

func (u *userRedPacketRepo) HandleRedPacket(ctx context.Context, txList ...func(tx *gorm.DB) error) error {
	tx := u.database.RDB(ctx).Begin()
	for _, f := range txList {
		if err := f(tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func (u *userRedPacketRepo) ConfirmTx(ctx context.Context, packet *domain.UserRedPacket) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		txConn := tx.Model(domain.UserRedPacket{}).
			Where("user_red_packet_id = ?", packet.UserRedPacketId).
			Updates(domain.UserRedPacket{
				Status:  constant.UserRedPacketSucc,
			})
		return txConn.Error
	}
}

func (u *userRedPacketRepo) WaitRecvTx(ctx context.Context, packet *domain.UserRedPacket) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		txConn := tx.Model(domain.UserRedPacket{}).
			Where("user_red_packet_id = ?", packet.UserRedPacketId).
			Updates(domain.UserRedPacket{
				Status:  constant.UserRedPacketWaitRecv,
				Package: packet.Package,
				Remark:  packet.Remark,
			})
		return txConn.Error
	}
}

func (u *userRedPacketRepo) CancelTx(ctx context.Context, packet *domain.UserRedPacket) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		txConn := tx.Model(domain.UserRedPacket{}).
			Where("user_red_packet_id = ?", packet.UserRedPacketId).
			Updates(domain.UserRedPacket{
				Status: constant.UserRedPacketFail,
				Remark: packet.Remark,
			})
		return txConn.Error
	}
}

func (u *userRedPacketRepo) GetRedPacketByStatus(ctx context.Context, userId uint64, status string) (*domain.UserRedPacket, error) {
	data := domain.UserRedPacket{}
	conn := u.database.RDB(ctx).Where("user_id = ?", userId).
		Where("status = ?", status).
		First(&data)
	if conn.RowsAffected == 0 {
		return nil, errors.NotFound("red packet", "未找到红包领取记录", fmt.Sprintf("%d", userId))
	}
	return &data, nil
}
