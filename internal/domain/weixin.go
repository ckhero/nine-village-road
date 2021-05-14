/**
 *@Description
 *@ClassName weixin
 *@Date 2021/5/13 上午9:38
 *@Author ckhero
 */

package domain

import "context"

type AppletRed struct {
	MchBillno uint64
	MchName string
	OpenId string
	TotalAmount uint64
	TotalNum int32
	Wishing string
	ActName string
	Remark string
}

type AppletRedPaySign struct {
	Timestamp string
	NonceStr string
	Package string
	SignType string
	PaySign string
}

type WalletTransfer struct {
	TradeNo uint64
	OpenId string
	CheckName string
	Amount uint64
	Desc string
}

type Code2Session struct {
	OpenId string
	SessionKey string
	UnionId string
}

type WeixinRepo interface {
	SendAppletRed(ctx context.Context, data *AppletRed) (*AppletRedPaySign, error)
	WalletTransfer(ctx context.Context, data *WalletTransfer) error
	Code2Session(ctx context.Context, code string) (*Code2Session, error)
}

type WeixinUsecase interface {
	SendAppletRed(ctx context.Context, openId string) (*AppletRedPaySign, error)
}