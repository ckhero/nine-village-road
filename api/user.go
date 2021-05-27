/**
 *@Description
 *@ClassName user
 *@Date 2021/5/13 下午2:20
 *@Author ckhero
 */

package api

type LogigRsp struct {
	Token string `json:"token"`
}

type SendAppletRedRsp struct {
	Timestamp string `json:"timestamp"`
	NonceStr string `json:"nonce_str"`
	Package string `json:"package"`
	SignType string `json:"sign_type"`
	PaySign string `json:"pay_sign"`
}

type WalletTransferRsp struct {
	Amount uint64 `json:"amount"`
}

type UserScenic struct {
	Scenic string `json:"scenic"`
}

type ListUserScenicRsp struct {
	List []*UserScenic `json:"list"`
}