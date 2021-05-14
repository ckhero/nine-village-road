/**
 *@Description
 *@ClassName weixin_usecase
 *@Date 2021/5/13 上午10:24
 *@Author ckhero
 */

package usecase

import (
	"context"
	"github.com/ckhero/go-common/util/uuid"
	"nine-village-road/internal/domain"
)

type weixinUsecase struct {
	repo domain.WeixinRepo
}

func NewWeixinUsecase(repo domain.WeixinRepo) domain.WeixinUsecase {
	return &weixinUsecase{
		repo: repo,
	}
}

func(w *weixinUsecase) SendAppletRed(ctx context.Context, openId string) (*domain.AppletRedPaySign, error) {
	return w.repo.SendAppletRed(ctx, &domain.AppletRed{
		MchBillno:   uuid.GenUUID(),
		MchName:    "测试",
		OpenId:      openId,
		TotalAmount: 30,
		TotalNum:    1,
		Wishing:     "测试祝福语",
		ActName:     "测试活动名字",
		Remark:      "测试备注",
	})

	//return w.repo.WalletTransfer(ctx, &domain.WalletTransfer{
	//	TradeNo:   uuid.GenUUID(),
	//	OpenId:    openId,
	//	CheckName: "NO_CHECK",
	//	Amount:    30,
	//	Desc:      "测试",
	//})
}

