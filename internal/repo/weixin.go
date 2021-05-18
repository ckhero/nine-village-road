/**
 *@Description
 *@ClassName wx_pay
 *@Date 2021/5/13 上午9:36
 *@Author ckhero
 */

package repo

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ckhero/go-common/config"
	"github.com/ckhero/go-common/errors"
	"github.com/ckhero/go-common/logger"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/pkg/util"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/silenceper/wechat/v2/miniprogram"
	"github.com/silenceper/wechat/v2/miniprogram/qrcode"
	"hash"
	"net/url"
	"nine-village-road/internal/domain"
	"strconv"
	"strings"
	"time"
)

type weixinRepo struct {
	miniPayClient *wechat.Client
	miniClient    *miniprogram.MiniProgram
}

func NewWeixinRepo(miniPayClient *wechat.Client, miniClient *miniprogram.MiniProgram) domain.WeixinRepo {
	return &weixinRepo{miniPayClient: miniPayClient, miniClient: miniClient}
}

// https://pay.weixin.qq.com/wiki/doc/api/tools/miniprogram_hb.php?chapter=18_2&index=3
func (w *weixinRepo) SendAppletRed(ctx context.Context, data *domain.AppletRed) (string, string, error) {
	logger.GetLoggerWithBody(ctx, data).Info("红包参数")

	cfg := config.GetWeixinPayCfg()
	p12 := cfg.CertP12
	certPem := cfg.CertPem
	keyPem := cfg.KeyPem
	_ = w.miniPayClient.AddCertPkcs12FilePath(p12)
	_ = w.miniPayClient.AddCertPemFilePath(certPem, keyPem)
	// 初始化参数结构体
	bm := make(gopay.BodyMap)
	// 随机字符串，不长于32位
	bm.Set("nonce_str", util.GetRandomString(32))
	// 商户订单号
	bm.Set("mch_billno", data.MchBillno)
	// 商户号
	bm.Set("mch_id", cfg.MchId)
	// 公众号appid
	bm.Set("wxappid", cfg.AppIdApp)
	// 商户名称
	bm.Set("send_name", data.MchName)
	// 用户openid
	bm.Set("re_openid", data.OpenId)
	// 付款金额 分
	bm.Set("total_amount", data.TotalAmount)
	//  红包发放总人数
	bm.Set("total_num", data.TotalNum)
	// 红包祝福语
	bm.Set("wishing", data.Wishing)
	// 备注
	bm.Set("act_name", data.ActName)
	// 活动名称
	bm.Set("remark", data.Remark)
	// 通知用户形式 通过JSAPI方式领取红包,小程序红包固定传MINI_PROGRAM_JSAPI

	bm.Set("notify_way", "MINI_PROGRAM_JSAPI")
	bm.Set("scene_id", "PRODUCT_1")
	rsp, err := w.miniPayClient.SendAppletRed(bm)
	//fmt.Print(rsp, err)
	logger.GetLoggerWithBody(ctx, rsp).Info("发送红包反馈")

	if err != nil || rsp.ResultCode != "SUCCESS" {
		msg := "红包发送失败"
		if rsp != nil {
			msg = rsp.ReturnMsg
		}
		return "", msg, errors.InternalServer("weixin", msg, "发送红包失败")
	}
	return url.QueryEscape(rsp.Packages), rsp.ReturnMsg, nil
}

func (w *weixinRepo) AppletRedPaySign(ctx context.Context,  pack string) (*domain.AppletRedPaySign) {
	cfg := config.GetWeixinPayCfg()
	paySign := &domain.AppletRedPaySign{
		Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
		NonceStr:  util.GetRandomString(32),
		Package:   pack,
		SignType:  "MD5",
	}
	paySign.PaySign = paySignV1(cfg.AppIdApp, paySign.NonceStr, paySign.Package, wechat.SignType_MD5, paySign.Timestamp, cfg.ApiKey)
	return paySign
}

func paySignV1(appId, nonceStr, packages, signType, timeStamp, apiKey string) string {
	var (
		buffer strings.Builder
		h      hash.Hash
	)
	buffer.WriteString("appId=")
	buffer.WriteString(appId)
	buffer.WriteString("&nonceStr=")
	buffer.WriteString(nonceStr)
	buffer.WriteString("&package=")
	buffer.WriteString(packages)
	buffer.WriteString("&timeStamp=")
	buffer.WriteString(timeStamp)
	buffer.WriteString("&key=")
	buffer.WriteString(apiKey)
	h = md5.New()
	h.Write([]byte(buffer.String()))
	return hex.EncodeToString(h.Sum(nil))
}

func (w *weixinRepo) WalletTransfer(ctx context.Context, data *domain.WalletTransfer) (string, error) {
	logger.GetLoggerWithBody(ctx, data).Info("企业付款")

	cfg := config.GetWeixinPayCfg()
	p12 := cfg.CertP12
	certPem := cfg.CertPem
	keyPem := cfg.KeyPem
	_ = w.miniPayClient.AddCertPkcs12FilePath(p12)
	_ = w.miniPayClient.AddCertPemFilePath(certPem, keyPem)
	// 初始化参数结构体
	bm := make(gopay.BodyMap)
	// 公众号appid
	bm.Set("mch_appid", cfg.AppIdApp)
	// 商户号
	bm.Set("mchid", cfg.MchId)
	// 随机字符串，不长于32位
	bm.Set("nonce_str", util.GetRandomString(32))
	// 商户订单号
	bm.Set("partner_trade_no", data.TradeNo)
	// 用户openid
	bm.Set("openid", data.OpenId)
	// 不校验真实姓名
	bm.Set("check_name", data.CheckName)
	// 付款金额 分
	bm.Set("amount", data.Amount)
	bm.Set("desc", data.Desc)
	bm.Set("spbill_create_ip", "47.100.86.135")

	rsp, err := w.miniPayClient.Transfer(bm)
	fmt.Print(rsp, err)
	logger.GetLoggerWithBody(ctx, rsp).Info("企业付款反馈")

	if err != nil || rsp.ResultCode != "SUCCESS" {
		msg := "红包发送失败"
		if rsp != nil {
			msg = rsp.ReturnMsg
		}
		return msg, errors.InternalServer("weixin", fmt.Sprintf("%v", err), "企业付款失败")
	}
	return rsp.ReturnMsg, nil
}

func (w *weixinRepo) Code2Session(ctx context.Context, code string) (*domain.Code2Session, error) {
	res, err := w.miniClient.GetAuth().Code2Session(code)
	if err != nil {
		//return nil, xerrors.Wrapf(err, "code [%s]", code)
		return nil, errors.Newf(errors.Code(err), "weixin", "Code解密失败", "%+v", err)
	}

	return &domain.Code2Session{
		OpenId:     res.OpenID,
		SessionKey: res.SessionKey,
		UnionId:    res.UnionID,
	}, nil
}

func (w *weixinRepo) QRCode(ctx context.Context, scenic string) ([]byte, error) {
	res, err := w.miniClient.GetQRCode().GetWXACodeUnlimit(qrcode.QRCoder{
		Scene:     scenic,
	})
	if err != nil {
		//return nil, xerrors.Wrapf(err, "code [%s]", code)
		return nil, errors.Newf(errors.Code(err), "weixin", "二维码生成失败", "%+v", err)
	}

	return res, nil
}
