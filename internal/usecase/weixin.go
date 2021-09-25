/**
 *@Description
 *@ClassName weixin_usecase
 *@Date 2021/5/13 上午10:24
 *@Author ckhero
 */

package usecase

import (
	"context"
	"encoding/base64"
	"github.com/ckhero/go-common/errors"
	"github.com/ckhero/go-common/util/uuid"
	"google.golang.org/grpc/codes"
	"io/ioutil"
	"nine-village-road/internal/domain"
	"nine-village-road/pkg/constant"
	"nine-village-road/pkg/rand_amount"
	"os"
)

type weixinUsecase struct {
	repo           domain.WeixinRepo
	redPacketRepo  domain.UserRedPacketRepo
	uerRepo        domain.UserRepo
	redPacketLimitRepo domain.RedPacketLimitRepo
}

func NewWeixinUsecase(repo domain.WeixinRepo, redPacketRepo domain.UserRedPacketRepo, uerRepo domain.UserRepo, redPacketLimitRepo domain.RedPacketLimitRepo) domain.WeixinUsecase {
	return &weixinUsecase{
		repo:           repo,
		redPacketRepo:  redPacketRepo,
		uerRepo:        uerRepo,
		redPacketLimitRepo: redPacketLimitRepo,
	}
}

func (w *weixinUsecase) SendAppletRed(ctx context.Context, user *domain.User) (*domain.AppletRedPaySign, error) {

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
			MchName:     "测试",
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
				w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId, constant.UserRecvStatusRecving, constant.UserRecvStatusInit),
				w.redPacketRepo.CancelTx(ctx, redPacket),
				w.redPacketLimitRepo.UpdateTx(ctx, -int(redPacket.Amount), -constant.LimitDefaultNum),
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

func (w *weixinUsecase) WalletTransfer(ctx context.Context, user *domain.User) (*domain.UserRedPacket, error) {
	var redPacket *domain.UserRedPacket
	var err error
	if user.IsRecving() {
		redPacket, err = w.redPacketRepo.GetRedPacketByStatus(ctx, user.UserId, constant.UserRecvStatusRecved)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
	}
	//weights := map[uint64]uint64{
	//	110 : 46,
	//	220 : 30,
	//	330 : 30,
	//	550 : 30,
	//	660 : 30,
	//	880 : 20,
	//	990 : 10,
	//	8800 : 1,
	//}
	weights := map[uint64]uint64{
		30 : 50,
		40 : 25,
		50 : 25,
	}
	// 不存在记录
	if redPacket == nil {
		redPacket = &domain.UserRedPacket{
			UserRedPacketId: uuid.GenUUID(),
			UserId:          user.UserId,
			OpenId:          user.OpenId,
			TradeNo:         uuid.GenUUID(),
			Amount:          rand_amount.GetRandAmount(weights),
		}
		// 创建数据
		if err := w.redPacketRepo.HandleRedPacket(ctx,
			w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId, user.RecvStatus, constant.UserRecvStatusRecving),
			w.redPacketRepo.CreateTx(ctx, redPacket),
			w.redPacketLimitRepo.UpdateTx(ctx, int(redPacket.Amount), constant.LimitDefaultNum),
		); err != nil {
			return nil, err
		}
		var err error
		redPacket.Remark, err = w.repo.WalletTransfer(ctx, &domain.WalletTransfer{
			TradeNo:   redPacket.TradeNo,
			OpenId:    redPacket.OpenId,
			CheckName: "NO_CHECK",
			Amount:    redPacket.Amount,
			Desc:      "点亮红包",
		})
		// 领取失败
		if err != nil {
			if err := w.redPacketRepo.HandleRedPacket(ctx,
				w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId, constant.UserRecvStatusRecving, constant.UserRecvStatusInit),
				w.redPacketRepo.CancelTx(ctx, redPacket),
				w.redPacketLimitRepo.UpdateTx(ctx, -int(redPacket.Amount), -constant.LimitDefaultNum),
			); err != nil {
				return nil, err
			}
			return nil, errors.Newf(codes.Unknown, "weixin", "红包领取失败", "userid [%d] openId [%s]", user.UserId, user.OpenId)
		}
	}
	// 领取成功

	if err := w.redPacketRepo.HandleRedPacket(ctx,
		w.uerRepo.UpdateRecvStatusTx(ctx, user.UserId, constant.UserRecvStatusRecving, constant.UserRecvStatusRecved),
		w.redPacketRepo.ConfirmTx(ctx, redPacket),
	); err != nil {
		return nil, err
	}
	return redPacket, nil
}

func (w *weixinUsecase) QRCode(ctx context.Context, scenic string) ([]byte, error) {
	qrCode, err := w.repo.QRCode(ctx, scenic)
	if err != nil {
		return nil, err
	}
	// base64 > png
	sourcestring := base64.StdEncoding.EncodeToString(qrCode)
	tmpPath := "./qr_code/a.png.txt"
	//写入临时文件
	_ = ioutil.WriteFile(tmpPath, []byte(sourcestring), 0667)
	//读取临时文件
	cc, _ := ioutil.ReadFile(tmpPath)

	//解压
	dist, _ := base64.StdEncoding.DecodeString(string(cc))
	//写入新文件
	f, _ := os.OpenFile("./qr_code/" + scenic + ".png", os.O_RDWR|os.O_CREATE, os.ModePerm)
	defer func() {
		_ = f.Close()
	}()
	_, _ = f.Write(dist)
	return qrCode, nil
}
