package unionpay

import (
	"bytes"
	"github.com/smartwalle/unionpay/internal"
	"net/url"
)

const (
	kFrontTrans = "/gateway/api/frontTransReq.do"
	kBackTrans  = "/gateway/api/backTransReq.do"
	kQueryTrans = "/gateway/api/queryTrans.do"
	kAppTrans   = "/gateway/api/appTransReq.do"
)

// CreateWebPayment 消费接口-创建网页支付。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=754&apiservId=448&version=V2.2&bussType=0
//
// 返回值为 HTML 代码，需要在浏览器中执行该代码。
//
// 注意：txnTime 的格式为 YYYYMMDDhhmmss，交易状态查询接口(GetPayment) 需要用到。
func (this *Client) CreateWebPayment(orderId, txnTime, amount, frontURL, backURL string, opts ...CallOption) (string, error) {
	var values = url.Values{}
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "07")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	values.Set("txnSubType", "01")    // 01：自助消费，通过地址的方式区分前台消费和后台消费（含无跳转支付） 03：分期付款
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}
	values.Set("bizType", "000201") // 业务类型，000201 - B2C网关支付和手机wap支付
	values.Set("txnType", "01")
	values.Set("orderId", orderId)
	values.Set("txnTime", txnTime)
	values.Set("txnAmt", amount)
	values.Set("frontUrl", frontURL)
	values.Set("backUrl", backURL)

	values, err := this.URLValues(values)
	if err != nil {
		return "", err
	}

	var buff = bytes.NewBufferString("")
	this.frontTransTpl.Execute(buff, map[string]interface{}{
		"Values": values,
		"Action": this.host + kFrontTrans,
	})
	return buff.String(), nil
}

// CreateAppPayment 消费接口-创建 App 支付。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?apiservId=3021&acpAPIId=961&bussType=0
//
// 返回值为为客户端调用银联 SDK 需要的 tn。
//
// 注意：txnTime 的格式为 YYYYMMDDhhmmss，交易状态查询接口(GetPayment) 需要用到。
func (this *Client) CreateAppPayment(orderId, txnTime, amount, backURL string, opts ...CallOption) (string, error) {
	var values = url.Values{}
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "08")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	values.Set("txnSubType", "01")    // 01：自助消费，通过地址的方式区分前台消费和后台消费（含无跳转支付） 03：分期付款
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}
	values.Set("bizType", "000201") // 业务类型，000201 - B2C网关支付和手机wap支付
	values.Set("txnType", "01")
	values.Set("orderId", orderId)
	values.Set("txnTime", txnTime)
	values.Set("txnAmt", amount)
	values.Set("backUrl", backURL)

	var rValues, err = this.Request(kAppTrans, values)
	if err != nil {
		return "", err
	}
	return rValues.Get("tn"), nil
}

// GetPayment 交易状态查询接口 https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=757&apiservId=448&version=V2.2&bussType=0
//
// 注意：txnTime 为创建支付提供的 txnTime, 商户的 orderId 和 txnTime 组成唯一订单信息。
func (this *Client) GetPayment(orderId, txnTime string, opts ...CallOption) (*Payment, error) {
	var values = url.Values{}
	values.Set("accessType", "0")
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}
	values.Set("bizType", "000000")
	values.Set("txnType", "00")
	values.Set("txnSubType", "00")
	values.Set("orderId", orderId)
	values.Set("txnTime", txnTime)

	var rValues, err = this.Request(kQueryTrans, values)
	if err != nil {
		return nil, err
	}

	var payment *Payment
	if err = internal.DecodeValues(rValues, &payment); err != nil {
		return nil, err
	}

	return payment, nil
}

// RevokePayment 消费撤销
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=755&apiservId=448&version=V2.2&bussType=0
func (this *Client) RevokePayment() (map[string]string, error) {
	return nil, nil
}

// Refund 退货接口
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=756&apiservId=448&version=V2.2&bussType=0
func (this *Client) Refund() (map[string]string, error) {
	return nil, nil
}
