package unionpay

import (
	"net/url"
	"time"
)

// CreateAccountPayment 无跳转支付-消费接口。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=814&apiservId=449&version=V2.2&bussType=0
//
// 文档地址：https://open.unionpay.com/upload/download/%E6%97%A0%E8%B7%B3%E8%BD%AC%E6%94%AF%E4%BB%98%E4%BA%A7%E5%93%81%E6%8E%A5%E5%8F%A3%E8%A7%84%E8%8C%83V2.0.pdf
//
// orderId：商户消费订单号。
//
// amount：交易金额，单位分，不要带小数点。
//
// backURL：后台通知地址。
//
// accNo：账号、卡号。
func (c *Client) CreateAccountPayment(orderId, amount, backURL, accNo string, customer *Customer, opts ...CallOption) (*AccountPayment, error) {
	var values = url.Values{}
	// 此处的参数可被 WithPayload() 替换
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "07")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	values.Set("bizType", "000301")   // 业务类型，000301 - 认证支付2.0
	values.Set("txnType", "01")
	values.Set("txnSubType", "01") // 01：自助消费，通过地址的方式区分前台消费和后台消费（含无跳转支付） 03：分期付款
	values.Set("txnTime", time.Now().Format("20060102150405"))
	values.Set("accType", "01") // 账号类型 后台类交易且卡号上送； 跨行收单且收单机构收集银行卡信息时上送 01：银行卡 02：存折 03：IC卡 默认取值：01 取值“03”表示以IC终端发起的IC卡交易，IC作为普通银行卡进行支付时，此域填写为“01”
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}

	values.Set("orderId", orderId)
	values.Set("txnAmt", amount)
	values.Set("backUrl", backURL)

	values.Set("encryptCertId", c.EncryptCertId())
	acc, err := c.Encrypt(accNo)
	if err != nil {
		return nil, err
	}
	values.Set("accNo", acc)

	customerInfo, err := c.EncryptCustomer(customer, accNo)
	if err != nil {
		return nil, err
	}
	values.Set("customerInfo", customerInfo)

	rValues, err := c.Request(kBackTrans, values)
	if err != nil {
		return nil, err
	}

	var payment *AccountPayment
	if err = DecodeValues(rValues, &payment); err != nil {
		return nil, err
	}
	return payment, nil
}

// ReverseAccountPayment 无跳转支付-冲正（退货）。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=815&apiservId=449&version=V2.2&bussType=0
//
// orderId：商户订单号。
//
// txnTime：订单发送时间，格式为 YYYYMMDDhhmmss，orderId 和 txnTime 组成唯一订单信息。
//
// 冲正必须与原始消费在同一天（准确讲是昨日23:00至本日23:00之间）。 冲正交易，仅用于超时无应答等异常场景，只有发生支付系统超时或者支付结果未知时可调用冲正，其他正常支付的订单如果需要实现相通功能，请调用消费撤销或者退货。
func (c *Client) ReverseAccountPayment(orderId, txnTime string, opts ...CallOption) (*Reverse, error) {
	var values = url.Values{}
	// 此处的参数可被 WithPayload() 替换
	values.Set("accessType", "0")
	values.Set("currencyCode", "156") // 交易币种 156 - 人民币
	values.Set("channelType", "07")   // 渠道类型，这个字段区分B2C网关支付和手机wap支付；07 - PC,平板  08 - 手机
	values.Set("bizType", "000000")   // 业务类型
	values.Set("txnType", "99")
	values.Set("txnSubType", "01")
	//values.Set("txnTime", time.Now().Format("20060102150405"))
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}

	values.Set("orderId", orderId)
	values.Set("txnTime", txnTime)

	var rValues, err = c.Request(kQueryTrans, values)
	if err != nil {
		return nil, err
	}

	var reverse *Reverse
	if err = DecodeValues(rValues, &reverse); err != nil {
		return nil, err
	}
	return reverse, nil
}
