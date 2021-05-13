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

type Code2Session struct {
	OpenId string
	SessionKey string
	UnionId string
}

type WeixinRepo interface {
	SendAppletRed(ctx context.Context, data *AppletRed) error
	Code2Session(ctx context.Context, code string) (*Code2Session, error)
}

type WeixinUsecase interface {
	SendAppletRed(ctx context.Context) error
}