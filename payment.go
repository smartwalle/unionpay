package unionpay

import (
	"bytes"
	"github.com/smartwalle/unionpay/internal"
	"net/url"
	"time"
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
// orderId：商户消费订单号。
//
// amount：交易金额，单位分，不要带小数点。
//
// frontURL：前台通知地址。
//
// backURL：后台通知地址。
func (this *Client) CreateWebPayment(orderId, amount, frontURL, backURL string, opts ...CallOption) (*WebPayment, error) {
	var values = url.Values{}
	// 此处的参数可被 WithPayload() 替换
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "07")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	values.Set("txnSubType", "01")    // 01：自助消费，通过地址的方式区分前台消费和后台消费（含无跳转支付） 03：分期付款
	values.Set("txnTime", time.Now().Format("20060102150405"))
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}

	// 固定参数
	values.Set("bizType", "000201") // 业务类型，000201 - B2C网关支付和手机wap支付
	values.Set("txnType", "01")

	values.Set("orderId", orderId)
	values.Set("txnAmt", amount)
	values.Set("frontUrl", frontURL)
	values.Set("backUrl", backURL)

	values, err := this.URLValues(values)
	if err != nil {
		return nil, err
	}

	var buff = bytes.NewBufferString("")
	if err := this.frontTransTpl.Execute(buff, map[string]interface{}{"Values": values, "Action": this.host + kFrontTrans}); err != nil {
		return nil, err
	}

	var payment = &WebPayment{}
	payment.Code = CodeSuccess
	payment.HTML = buff.String()
	payment.Version = values.Get("version")
	payment.BizType = values.Get("bizType")
	payment.TxnTime = values.Get("txnTime")
	payment.TxnType = values.Get("txnType")
	payment.TxnSubType = values.Get("txnSubType")
	payment.AccessType = values.Get("accessType")
	payment.MerId = values.Get("merId")
	payment.OrderId = values.Get("orderId")
	return payment, nil
}

// CreateAppPayment 消费接口-创建 App 支付。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?apiservId=3021&acpAPIId=961&bussType=0
//
// orderId：商户消费订单号。
//
// amount：交易金额，单位分，不要带小数点。
//
// backURL：后台通知地址。
func (this *Client) CreateAppPayment(orderId, amount, backURL string, opts ...CallOption) (*AppPayment, error) {
	var values = url.Values{}
	// 此处的参数可被 WithPayload() 替换
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "08")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	values.Set("txnSubType", "01")    // 01：自助消费，通过地址的方式区分前台消费和后台消费（含无跳转支付） 03：分期付款
	values.Set("txnTime", time.Now().Format("20060102150405"))
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}

	// 固定参数
	values.Set("bizType", "000201") // 业务类型，000201 - B2C网关支付和手机wap支付
	values.Set("txnType", "01")

	values.Set("orderId", orderId)
	values.Set("txnAmt", amount)
	values.Set("backUrl", backURL)

	var rValues, err = this.Request(kAppTrans, values)
	if err != nil {
		return nil, err
	}

	var payment *AppPayment
	if err = internal.DecodeValues(rValues, &payment); err != nil {
		return nil, err
	}
	return payment, nil
}

// GetTransaction 交易状态查询接口
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=757&apiservId=448&version=V2.2&bussType=0
//
// orderId：商户订单号。
//
// txnTime：订单发送时间，格式为 YYYYMMDDhhmmss，orderId 和 txnTime 组成唯一订单信息。
//
// 注：
// 应答报文中，“应答码”即respCode字段，表示的是查询交易本身的应答，即查询这个动作是否成功，不代表被查询交易的状态；
// 若查询动作成功，即应答码为“00“，则根据“原交易应答码”即origRespCode来判断被查询交易是否成功。此时若origRespCode为00，则表示被查询交易成功。
func (this *Client) GetTransaction(orderId, txnTime string, opts ...CallOption) (*Transaction, error) {
	var values = url.Values{}
	// 此处的参数可被 WithPayload() 替换
	values.Set("accessType", "0")
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}

	// 固定参数
	values.Set("bizType", "000000")
	values.Set("txnType", "00")
	values.Set("txnSubType", "00")

	values.Set("orderId", orderId)
	values.Set("txnTime", txnTime)

	var rValues, err = this.Request(kQueryTrans, values)
	if err != nil {
		return nil, err
	}

	var transaction *Transaction
	if err = internal.DecodeValues(rValues, &transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}

// Revoke 消费撤销。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=755&apiservId=448&version=V2.2&bussType=0
//
// queryId：原消费交易返回的的queryId，可以从消费交易后台通知接口中或者交易状态查询接口(GetTransaction)中获取。
//
// orderId：商户撤销订单号，和要消费撤销的订单号没有关系。后续可用本 orderId 和返回结构体中的 TxnTime 通过交易状态查询接口(GetTransaction) 查询消费撤销信息。
//
// amount：退货金额，单位分，不要带小数点。
//
// backURL：后台通知地址。
//
// 消费撤销和退货的区别：
//
// 消费撤销仅能对当天的消费做，必须为全额，一般当日或第二日到账，可能存在极少数银行不支持。
//
// 退货（除二维码产品外）能对90天内（见注1、2）的消费做（包括当天），支持部分退货或全额退货，到账时间较长，一般1-10天（多数发卡行5天内，但工行可能会10天），所有银行都支持。
//
// 二维码产品退货支持30天，30天以上的退货可能可以发成功（失败应该会同步报错“原交易不存在或状态不正确[2011000]”之类的信息），但不保证一定可以成功。
//
// 注1：以上的天均指清算日，一般前一日23点至当天23点为一个清算日。
//
// 注2：系统实际支持330天的退货，但银联对发卡行的退货支持要求仅为90天，超过90天的退货发卡行虽然也会承兑，但可能为人工处理，到账速度较慢。330天以上的退货也可能成功，但不保证一定可以成功（失败应该会同步报错4040007之类的应答码），建议直接给用户转账来退款。
func (this *Client) Revoke(queryId, orderId, amount, backURL string, opts ...CallOption) (*Revoke, error) {
	var values = url.Values{}
	// 此处的参数可被 WithPayload() 替换
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "07")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	values.Set("txnTime", time.Now().Format("20060102150405"))
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}

	// 固定参数
	values.Set("bizType", "000201") // 业务类型，000201 - B2C网关支付和手机wap支付
	values.Set("txnType", "31")
	values.Set("txnSubType", "00")

	values.Set("origQryId", queryId)
	values.Set("orderId", orderId)
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
// queryId：原消费交易返回的的queryId，可以从消费交易后台通知接口中或者交易状态查询接口(GetTransaction)中获取。
//
// orderId：商户退货订单号，和要退款的订单号没有关系。后续可用本 orderId 和返回结构体中的 TxnTime 通过交易状态查询接口(GetTransaction) 查询退货信息。
//
// amount：退货金额，单位分，不要带小数点。
//
// backURL：后台通知地址。
//
// 消费撤销和退货的区别：
//
// 消费撤销仅能对当天的消费做，必须为全额，一般当日或第二日到账，可能存在极少数银行不支持。
//
// 退货（除二维码产品外）能对90天内（见注1、2）的消费做（包括当天），支持部分退货或全额退货，到账时间较长，一般1-10天（多数发卡行5天内，但工行可能会10天），所有银行都支持。
//
// 二维码产品退货支持30天，30天以上的退货可能可以发成功（失败应该会同步报错“原交易不存在或状态不正确[2011000]”之类的信息），但不保证一定可以成功。
//
// 注1：以上的天均指清算日，一般前一日23点至当天23点为一个清算日。
//
// 注2：系统实际支持330天的退货，但银联对发卡行的退货支持要求仅为90天，超过90天的退货发卡行虽然也会承兑，但可能为人工处理，到账速度较慢。330天以上的退货也可能成功，但不保证一定可以成功（失败应该会同步报错4040007之类的应答码），建议直接给用户转账来退款。
func (this *Client) Refund(queryId, orderId, amount, backURL string, opts ...CallOption) (*Refund, error) {
	var values = url.Values{}

	// 此处的参数可被 WithPayload() 替换
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "07")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	values.Set("txnTime", time.Now().Format("20060102150405"))
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}

	// 固定参数
	values.Set("bizType", "000201") // 业务类型，000201 - B2C网关支付和手机wap支付
	values.Set("txnType", "04")
	values.Set("txnSubType", "00")

	values.Set("origQryId", queryId)
	values.Set("orderId", orderId)
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
