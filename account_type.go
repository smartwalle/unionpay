package unionpay

type AccountPayment struct {
	Error
	QueryId      string `query:"queryId"`      // 查询流水号
	AcqInsCode   string `query:"acqInsCode"`   // 收单机构代码
	TN           string `query:"tn"`           // 银联受理订单号, 客户端调用银联 SDK 需要的银联订单号(tn)
	AccNo        string `query:"accNo"`        // 账号
	PayType      string `query:"payType"`      // 支付方式
	PayCardType  string `query:"payCardType"`  // 支付卡类型
	BizType      string `query:"bizType"`      // 产品类型
	TxnTime      string `query:"txnTime"`      // 订单发送时间
	CurrencyCode string `query:"currencyCode"` // 交易币种
	TxnAmt       string `query:"txnAmt"`       // 交易金额
	TxnType      string `query:"txnType"`      // 交易类型
	TxnSubType   string `query:"txnSubType"`   // 交易子类
	AccessType   string `query:"accessType"`   // 接入类型
	ReqReserved  string `query:"reqReserved"`  // 请求方保留域
	MerId        string `query:"merId"`        // 商户代码
	OrderId      string `query:"orderId"`      // 商户订单号
	Reserved     string `query:"reserved"`     // 保留域
	Version      string `query:"version"`      // 版本号
}

type Reverse struct {
	Error
	BizType     string `query:"bizType"`     // 产品类型
	TxnTime     string `query:"txnTime"`     // 订单发送时间
	TxnType     string `query:"txnType"`     // 交易类型
	TxnSubType  string `query:"txnSubType"`  // 交易子类
	AccessType  string `query:"accessType"`  // 接入类型
	ReqReserved string `query:"reqReserved"` // 请求方保留域
	MerId       string `query:"merId"`       // 商户代码
	OrderId     string `query:"orderId"`     // 商户订单号
	Reserved    string `query:"reserved"`    // 保留域
	Version     string `query:"version"`     // 版本号
}
