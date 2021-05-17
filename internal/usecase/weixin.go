/**
 *@Description
 *@ClassName weixin_usecase
 *@Date 2021/5/13 上午10:24
 *@Author ckhero
 */

package usecase

import (
	"context"
	"github.com/ckhero/go-common/errors"
	"github.com/ckhero/go-common/util/uuid"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
)

type weixinUsecase struct {
	repo domain.WeixinRepo
	redPacketRepo domain.UserRedPacketRepo
	uerRepo domain.UserRepo
}

func NewWeixinUsecase(repo domain.WeixinRepo, redPacketRepo domain.UserRedPacketRepo, uerRepo domain.UserRepo) domain.WeixinUsecase {
	return &weixinUsecase{
		repo: repo,
		redPacketRepo: redPacketRepo,
		uerRepo: uerRepo,
	}
}

func(w *weixinUsecase) SendAppletRed(ctx context.Context, user *domain.User) (*domain.AppletRedPaySign, error) {

	var redPacket *domain.UserRedPacket
	var err error
	if user.IsRecving() {
		redPacket, err = w.redPacketRepo.GetRedPacketByStatus(ctx, user.UserId, constant.UserRedPacketWaitRecv)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
	}
	// 不存在记录
	if redPacket == nil {
		redPacket = &domain.UserRedPacket{
			UserRedPacketId: uuid.GenUUID(),
			UserId:          user.UserId,
			OpenId:          user.OpenId,
			TradeNo:         uuid.GenUUID(),
			Amount:          30,
		}
		// 创建数据
		if err := w.redPacketRepo.HandleRedPacket(ctx,
			w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId, user.RecvStatus, constant.UserRecvStatusRecving),
			w.redPacketRepo.CreateTx(ctx, redPacket),
		); err != nil {
			return nil, err
		}
		var err error
		redPacket.Package, redPacket.Remark, err = w.repo.SendAppletRed(ctx, &domain.AppletRed{
			MchBillno:   redPacket.TradeNo,
			MchName:    "测试",
			OpenId:      redPacket.OpenId,
			TotalAmount: redPacket.Amount,
			TotalNum:    1,
			Wishing:     "测试祝福语",
			ActName:     "测试活动名字",
			Remark:      "测试备注",
		})
		// 领取失败
		if err != nil {
			if err := w.redPacketRepo.HandleRedPacket(ctx,
				w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId,  constant.UserRecvStatusRecving, constant.UserRecvStatusInit),
				w.redPacketRepo.CancelTx(ctx, redPacket),
			); err != nil {
				return nil, err
			}
			return nil, err
		}
	}
	// 领取成功
	res := w.repo.AppletRedPaySign(ctx, redPacket.Package)

	if err := w.redPacketRepo.HandleRedPacket(ctx,
		w.redPacketRepo.WaitRecvTx(ctx, redPacket),
	); err != nil {
		return nil, err
	}
	return res, nil
}

func(w *weixinUsecase) WalletTransfer(ctx context.Context, user *domain.User) (error) {
	var redPacket *domain.UserRedPacket
	var err error
	if user.IsRecving() {
		redPacket, err = w.redPacketRepo.GetRedPacketByStatus(ctx, user.UserId, constant.UserRecvStatusRecved)
		if err !=nil && !errors.IsNotFound(err) {
			return err
		}
	}
	// 不存在记录
	if redPacket == nil {
		redPacket = &domain.UserRedPacket{
			UserRedPacketId: uuid.GenUUID(),
			UserId:          user.UserId,
			OpenId:          user.OpenId,
			TradeNo:         uuid.GenUUID(),
			Amount:          30,
		}
		// 创建数据
		if err := w.redPacketRepo.HandleRedPacket(ctx,
			w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId, user.RecvStatus, constant.UserRecvStatusRecving),
			w.redPacketRepo.CreateTx(ctx, redPacket),
		); err != nil {
			return err
		}
		var err error
		redPacket.Remark, err = w.repo.WalletTransfer(ctx, &domain.WalletTransfer{
			TradeNo:   redPacket.TradeNo,
			OpenId:    redPacket.OpenId,
			CheckName: "NO_CHECK",
			Amount:    redPacket.Amount,
			Desc:      "测试",
		})
		// 领取失败
		if err != nil {
			if err := w.redPacketRepo.HandleRedPacket(ctx,
				w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId,  constant.UserRecvStatusRecving, constant.UserRecvStatusInit),
				w.redPacketRepo.CancelTx(ctx, redPacket),
			); err != nil {
				return err
			}
			return err
		}
	}
	// 领取成功

	if err := w.redPacketRepo.HandleRedPacket(ctx,
		w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId,  constant.UserRecvStatusRecving, constant.UserRecvStatusRecved),
		w.redPacketRepo.ConfirmTx(ctx, redPacket),
	); err != nil {
		return err
	}
	return nil
}

