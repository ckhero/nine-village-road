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

func(w *weixinUsecase) SendAppletRed(ctx context.Context) error {
	return w.repo.SendAppletRed(ctx, &domain.AppletRed{
		MchBillno:   uuid.GenUUID(),
		MchName:    "测试",
		OpenId:      "o8M5t5Qrsg7yEsihXvOFZTIBiwSU",
		TotalAmount: 1,
		TotalNum:    1,
		Wishing:     "测试祝福语",
		ActName:     "测试活动名字",
		Remark:      "测试备注",
	})
}