package unionpay

// FrontTransParam https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=754&apiservId=448&version=V2.2&bussType=0
type FrontTransParam struct {
	BackURL      string // 后台通知地址
	FrontURL     string // 前台通知地址
	OrderId      string // 商户订单号
	CurrencyCode string // 交易币种
	TxnAmount    string // 交易金额 单位为分
	TxnSubType   string // 交易子类 01：自助消费，通过地址的方式区分前台消费和后台消费（含无跳转支付） 03：分期付款
	OrderDesc    string // 订单描述
}
