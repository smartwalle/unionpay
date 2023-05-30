package unionpay

type PaymentNotification struct {
	Error
	TxnType            string `query:"txnType"`            // 交易类型
	TxnSubType         string `query:"txnSubType"`         // 交易子类
	BizType            string `query:"bizType"`            // 产品类型
	AccessType         string `query:"accessType"`         // 接入类型
	AcqInsCode         string `query:"acqInsCode"`         // 收单机构代码
	MerId              string `query:"merId"`              // 商户代码
	OrderId            string `query:"orderId"`            // 商户订单号
	TxnTime            string `query:"txnTime"`            // 订单发送时间
	TxnAmt             string `query:"txnAmt"`             // 交易金额
	CurrencyCode       string `query:"currencyCode"`       // 交易币种
	ReqReserved        string `query:"reqReserved"`        // 请求方保留域
	Reserved           string `query:"reserved"`           // 保留域
	QueryId            string `query:"queryId"`            // 查询流水号
	SettleAmt          string `query:"settleAmt"`          // 清算金额
	SettleCurrencyCode string `query:"settleCurrencyCode"` // 清算币种
	SettleDate         string `query:"settleDate"`         // 清算日期
	TraceNo            string `query:"traceNo"`            // 系统跟踪号
	TraceTime          string `query:"traceTime"`          // 交易传输时间
	ExchangeDate       string `query:"exchangeDate"`       // 兑换日期
	ExchangeRate       string `query:"exchangeRate"`       // 清算汇率
	AccNo              string `query:"accNo"`              // 账号
	PayCardType        string `query:"payCardType"`        // 支付卡类型
	PayType            string `query:"payType"`            // 支付方式
	PayCardNo          string `query:"payCardNo"`          // 支付卡标识
	PayCardIssueName   string `query:"payCardIssueName"`   // 支付卡名称
	BindId             string `query:"bindId"`             // 绑定标识号
	InstalTransInfo    string `query:"instalTransInfo"`    // 分期付款信息域
	Version            string `query:"version"`            // 版本号
}

type RevokeNotification struct {
	Refund
	CurrencyCode       string `query:"currencyCode"`       // 交易币种
	SettleAmt          string `query:"settleAmt"`          // 清算金额
	SettleCurrencyCode string `query:"settleCurrencyCode"` // 清算币种
	SettleDate         string `query:"settleDate"`         // 清算日期
	TraceNo            string `query:"traceNo"`            // 系统跟踪号
	TraceTime          string `query:"traceTime"`          // 交易传输时间
	ExchangeDate       string `query:"exchangeDate"`       // 兑换日期
	ExchangeRate       string `query:"exchangeRate"`       // 清算汇率
	AccNo              string `query:"accNo"`              // 账号
}

type RefundNotification struct {
	Refund
	CurrencyCode       string `query:"currencyCode"`       // 交易币种
	SettleAmt          string `query:"settleAmt"`          // 清算金额
	SettleCurrencyCode string `query:"settleCurrencyCode"` // 清算币种
	SettleDate         string `query:"settleDate"`         // 清算日期
	TraceNo            string `query:"traceNo"`            // 系统跟踪号
	TraceTime          string `query:"traceTime"`          // 交易传输时间
	ExchangeDate       string `query:"exchangeDate"`       // 兑换日期
	ExchangeRate       string `query:"exchangeRate"`       // 清算汇率
	AccNo              string `query:"accNo"`              // 账号
}
