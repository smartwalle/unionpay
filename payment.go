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

// CreateWebPayment 消费接口-创建网页支付，返回值为 HTML 代码，需要在浏览器中执行该代码。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=754&apiservId=448&version=V2.2&bussType=0
//
// orderId：商户订单号。
//
// txnTime：订单发送时间，格式为 YYYYMMDDhhmmss，商户的 orderId 和 txnTime 组成唯一订单信息，交易状态查询接口(GetPayment) 需要用到。
//
// amount：交易金额，单位分，不要带小数点。
//
// frontURL：前台通知地址。
//
// backURL：后台通知地址。
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

// CreateAppPayment 消费接口-创建 App 支付，返回值为为客户端调用银联 SDK 需要的 tn。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?apiservId=3021&acpAPIId=961&bussType=0
//
// orderId：商户订单号。
//
// txnTime：订单发送时间，格式为 YYYYMMDDhhmmss，商户的 orderId 和 txnTime 组成唯一订单信息，交易状态查询接口(GetPayment) 需要用到。
//
// amount：交易金额，单位分，不要带小数点。
//
// backURL：后台通知地址。
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

// GetPayment 交易状态查询接口
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=757&apiservId=448&version=V2.2&bussType=0
//
// orderId：商户订单号。
//
// txnTime：订单发送时间，格式为 YYYYMMDDhhmmss，商户的 orderId 和 txnTime 组成唯一订单信息。
//
// 注：
// 应答报文中，“应答码”即respCode字段，表示的是查询交易本身的应答，即查询这个动作是否成功，不代表被查询交易的状态；
// 若查询动作成功，即应答码为“00“，则根据“原交易应答码”即origRespCode来判断被查询交易是否成功。此时若origRespCode为00，则表示被查询交易成功。
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

// Revoke 消费撤销。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=755&apiservId=448&version=V2.2&bussType=0
//
// queryId：原消费交易返回的的queryId，可以从消费交易后台通知接口中或者交易状态查询接口(GetPayment)中获取。
//
// orderId：商户订单号，和要消费撤销的订单号没有关系。后续可用 orderId 和 txnTime 通过交易状态查询接口(GetPayment) 查询消费撤销信息。
//
// txnTime：订单发送时间，格式为 YYYYMMDDhhmmss。
//
// amount：退货金额，单位分，不要带小数点。
//
// backURL：后台通知地址。
func (this *Client) Revoke(queryId, orderId, txnTime, amount, backURL string, opts ...CallOption) (*Revoke, error) {
	var values = url.Values{}
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "07")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}
	values.Set("bizType", "000201") // 业务类型，000201 - B2C网关支付和手机wap支付
	values.Set("txnType", "31")
	values.Set("txnSubType", "00")
	values.Set("origQryId", queryId)
	values.Set("orderId", orderId)
	values.Set("txnTime", txnTime)
	values.Set("txnAmt", amount)
	values.Set("backUrl", backURL)

	var rValues, err = this.Request(kBackTrans, values)
	if err != nil {
		return nil, err
	}

	var revoke *Revoke
	if err = internal.DecodeValues(rValues, &revoke); err != nil {
		return nil, err
	}
	return revoke, nil
}

// Refund 退货接口。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=756&apiservId=448&version=V2.2&bussType=0
//
// queryId：原消费交易返回的的queryId，可以从消费交易后台通知接口中或者交易状态查询接口(GetPayment)中获取。
//
// orderId：商户订单号，和要退款的订单号没有关系。后续可用 orderId 和 txnTime 通过交易状态查询接口(GetPayment) 查询退货信息。
//
// txnTime：订单发送时间，格式为 YYYYMMDDhhmmss。
//
// amount：退货金额，单位分，不要带小数点。
//
// backURL：后台通知地址。
func (this *Client) Refund(queryId, orderId, txnTime, amount, backURL string, opts ...CallOption) (*Refund, error) {
	var values = url.Values{}
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "07")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}
	values.Set("bizType", "000201") // 业务类型，000201 - B2C网关支付和手机wap支付
	values.Set("txnType", "04")
	values.Set("txnSubType", "00")
	values.Set("origQryId", queryId)
	values.Set("orderId", orderId)
	values.Set("txnTime", txnTime)
	values.Set("txnAmt", amount)
	values.Set("backUrl", backURL)

	var rValues, err = this.Request(kBackTrans, values)
	if err != nil {
		return nil, err
	}

	var refund *Refund
	if err = internal.DecodeValues(rValues, &refund); err != nil {
		return nil, err
	}
	return refund, nil
}
