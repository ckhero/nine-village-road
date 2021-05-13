/**
 *@Description
 *@ClassName wx_pay
 *@Date 2021/5/13 上午9:36
 *@Author ckhero
 */

package repo

import (
	"context"
	"github.com/ckhero/go-common/config"
	"github.com/ckhero/go-common/logger"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/pkg/util"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/miniprogram"
	"nine-village-road/internal/domain"
)

type weixinRepo struct {
	miniPayClient *wechat.Client
	miniClient    *miniprogram.MiniProgram
}

func NewWeixinRepo(miniPayClient *wechat.Client, miniClient *miniprogram.MiniProgram) domain.WeixinRepo {
	return &weixinRepo{miniPayClient: miniPayClient, miniClient: miniClient}
}

// https://pay.weixin.qq.com/wiki/doc/api/tools/miniprogram_hb.php?chapter=18_2&index=3
func (w *weixinRepo) SendAppletRed(ctx context.Context, data *domain.AppletRed) error {
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
	rsp, err := w.miniPayClient.SendAppletRed(bm)

	logger.GetLoggerWithBody(ctx, rsp).Info("发送红包反馈")

	if err != nil || rsp.ResultCode != "SUCCESS" {
		logger.GetLogger(ctx).Errorf("发送红包失败 %v", err)
		return errors.New(rsp.ErrCodeDes)
	}

	return err
}

func (w *weixinRepo) Code2Session(ctx context.Context, code string) (*domain.Code2Session, error) {
	res, err := w.miniClient.GetAuth().Code2Session(code)
	if err != nil {
		return nil, errors.Wrapf(err, "登录失败 code [%s]", code)
	}

	return &domain.Code2Session{
		OpenId:     res.OpenID,
		SessionKey: res.SessionKey,
		UnionId:    res.UnionID,
	}, nil
}
