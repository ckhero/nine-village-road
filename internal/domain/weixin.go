/**
 *@Description
 *@ClassName weixin
 *@Date 2021/5/13 上午9:38
 *@Author ckhero
 */

package domain

import (
	"context"
)

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
	Timestamp string `json:"timestamp"`
	NonceStr string `json:"nonce_str"`
	Package string `json:"package"`
	SignType string `json:"sign_type"`
	PaySign string `json:"pay_sign"`
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
	SendAppletRed(ctx context.Context, data *AppletRed) (string, string, error)
	AppletRedPaySign(ctx context.Context, pack string) (*AppletRedPaySign)
	WalletTransfer(ctx context.Context, data *WalletTransfer) (string, error)
	Code2Session(ctx context.Context, code string) (*Code2Session, error)
	QRCode(ctx context.Context, scenic string) ([]byte, error)
}

type WeixinUsecase interface {
	SendAppletRed(ctx context.Context, user *User) (*AppletRedPaySign, error)
	WalletTransfer(ctx context.Context, user *User) error
	QRCode(ctx context.Context, scenic string) ([]byte, error)
}