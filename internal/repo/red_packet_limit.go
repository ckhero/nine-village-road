/**
 *@Description
 *@ClassName red_packet_limit
 *@Date 2021/5/18 下午1:18
 *@Author ckhero
 */

package repo

import (
	"context"
	"github.com/ckhero/go-common/config"
	"github.com/ckhero/go-common/db"
	"github.com/ckhero/go-common/errors"
	"github.com/ckhero/go-common/util/uuid"
	"github.com/jinzhu/now"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
)

type redPacketLimitRepo struct {
	database *db.Database
}

func NewRedPacketLimitRepo(database *db.Database) domain.RedPacketLimitRepo {
	return &redPacketLimitRepo{database: database}
}

func (rl *redPacketLimitRepo) UpdateTx(ctx context.Context, amount, num int) func(tx *gorm.DB) error {
	return func(tx *gorm.DB) error {
		// 更新全局
		conn := tx.Model(domain.RedPacketLimit{}).
			Where("limit_type = ?", constant.LimitTypeGlobal).
			Updates(map[string]interface{}{
				"left_amount":   gorm.Expr("left_amount - ?", amount),
				"left_recv_num": gorm.Expr("left_recv_num - ?", num),
			})
		if conn.Error != nil || conn.RowsAffected == 0 {
			return errors.Newf(codes.Unknown, "red packet limit", "红包已经被领取完了", "")
		}
		// 更新当天
		currDay := now.BeginningOfDay().Format("2006-01-02")
		conn = tx.Model(domain.RedPacketLimit{}).
			Where("limit_type = ?", constant.LimitTypeDay).
			Where("start_date >= ?", currDay).
			Where("end_date >= ?", currDay).
			Updates(map[string]interface{}{
				"left_amount":   gorm.Expr("left_amount - ?", amount),
				"left_recv_num": gorm.Expr("left_recv_num - ?", num),
			})

		if conn.RowsAffected == 0 {
			// 创建数据
			cfg := config.GetRedPacketLimit()
			if cfg.Amount <= uint64(amount) || cfg.RecvNum <= uint64(num) || amount < 0 || num < 0 {
				return errors.Newf(codes.Unknown, "red packet limit", "红包配置有误",
					"cfg amount [%d] cfg num [%d] amount [%d] num [%d]", cfg.Amount, cfg.RecvNum, amount, num)
			}
			conn = tx.Create(domain.RedPacketLimit{
				RedPacketLimitId: uuid.GenUUID(),
				LimitType:        constant.LimitTypeDay,
				Amount:           cfg.Amount,
				LeftAmount:       cfg.Amount - uint64(amount),
				RecvNum:          cfg.RecvNum,
				LeftRecvNum:      cfg.RecvNum - uint64(num),
				StartDate:        currDay,
				EndDate:          currDay,
			})
			if conn.Error != nil {
				return errors.Newf(codes.Unknown, "red packet limit", "红包配置创建失败", "")
			}
		}
		if conn.Error != nil {
			return errors.Newf(codes.Unknown, "red packet limit", "今天的红包已经被领取完了，请明日再来", "")
		}
		return nil
	}
}
