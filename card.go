package unionpay

import (
	"net/url"
	"time"
)

// CreateCardPayment 消费接口-无跳转支付。
//
// orderId：商户消费订单号。
//
// amount：交易金额，单位分，不要带小数点。
//
// backURL：后台通知地址。
//
// accNo：卡号。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=814&apiservId=449&version=V2.2&bussType=0
//
// 文档地址：https://open.unionpay.com/upload/download/%E6%97%A0%E8%B7%B3%E8%BD%AC%E6%94%AF%E4%BB%98%E4%BA%A7%E5%93%81%E6%8E%A5%E5%8F%A3%E8%A7%84%E8%8C%83V2.0.pdf
func (this *Client) CreateCardPayment(orderId, amount, backURL, accNo string, customer *Customer, opts ...CallOption) (*CardPayment, error) {
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

	values.Set("encryptCertId", this.EncryptCertId())
	acc, err := this.Encrypt(accNo)
	if err != nil {
		return nil, err
	}
	values.Set("accNo", acc)

	customerInfo, err := this.EncryptCustomer(customer, accNo)
	if err != nil {
		return nil, err
	}
	values.Set("customerInfo", customerInfo)

	rValues, err := this.Request(kBackTrans, values)
	if err != nil {
		return nil, err
	}

	var payment *CardPayment
	if err = DecodeValues(rValues, &payment); err != nil {
		return nil, err
	}
	return payment, nil
}
