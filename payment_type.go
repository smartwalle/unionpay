package unionpay

// txnType
// 00：查询交易
// 01：消费
// 02：预授权
// 03：预授权完成
// 04：退货
// 05：圈存
// 11：代收
// 12：贷记
// 13：账单支付
// 14：转账（保留
// 21：批量交易
// 22：批量查询
// 31：消费撤销
// 32：预授权撤销
// 33：预授权完成撤销
// 72：实名认证 建立绑定关系
// 73：账单查询
// 74：解除绑定关系
// 75：查询绑定关系
// 76：文件传输
// 77：发送短信验证码交易
// 78：开通查询交易
// 79：开通交易
// 94：IC卡脚本通知
// 95：查询更新加密公钥证书

// bizType
// 000201：B2C网关支付
// 000301：认证支付 2.0
// 000302：评级支付
// 000401：贷记
// 000501：代收
// 000601：账单支付
// 000801：跨行收单
// 000901：绑定支付
// 001001：订购
// 000202：B2B

// channelType
// 05：语音
// 07：互联网
// 08：移动
// 16：数字机顶盒

type AppPayment struct {
	Error
	AcqInsCode  string `query:"acqInsCode"`  // 收单机构代码
	TN          string `query:"tn"`          // 银联受理订单号, 客户端调用银联 SDK 需要的银联订单号(tn)；
	Version     string `query:"version"`     // 版本号
	BizType     string `query:"bizType"`     // 产品类型
	TxnTime     string `query:"txnTime"`     // 订单发送时间
	TxnType     string `query:"txnType"`     // 交易类型
	TxnSubType  string `query:"txnSubType"`  // 交易子类
	AccessType  string `query:"accessType"`  // 接入类型
	ReqReserved string `query:"reqReserved"` // 请求方保留域
	MerId       string `query:"merId"`       // 商户代码
	OrderId     string `query:"orderId"`     // 商户订单号
	Reserved    string `query:"reserved"`    // 保留域
}

type Transaction struct {
	Error
	QueryId            string `query:"queryId"`            // 查询流水号
	TraceTime          string `query:"traceTime"`          // 交易传输时间
	TxnType            string `query:"txnType"`            // 交易类型
	TxnSubType         string `query:"txnSubType"`         // 交易子类
	SettleCurrencyCode string `query:"settleCurrencyCode"` // 清算币种
	SettleAmt          string `query:"settleAmt"`          // 清算金额
	SettleDate         string `query:"settleDate"`         // 清算日期
	TraceNo            string `query:"traceNo"`            // 系统跟踪号
	BindId             string `query:"bindId"`             // 绑定标识号
	ExchangeDate       string `query:"exchangeDate"`       // 兑换日期
	IssuerIdentifyMode string `query:"issuerIdentifyMode"` // 发卡机构识别模式
	CurrencyCode       string `query:"currencyCode"`       // 交易币种
	TxnAmt             string `query:"txnAmt"`             // 交易金额
	ExchangeRate       string `query:"exchangeRate"`       // 清算汇率
	CardTransData      string `query:"cardTransData"`      // 有卡交易信息域
	OrigRespCode       string `query:"origRespCode"`       // 原交易应答码
	OrigRespMsg        string `query:"origRespMsg"`        // 原交易应答信息
	AccNo              string `query:"accNo"`              // 账号
	PayType            string `query:"payType"`            // 支付方式
	PayCardNo          string `query:"payCardNo"`          // 支付卡标识
	PayCardType        string `query:"payCardType"`        // 支付卡类型
	PayCardIssueName   string `query:"payCardIssueName"`   // 支付卡名称
	Version            string `query:"version"`            // 版本号
	BizType            string `query:"bizType"`            // 产品类型
	TxnTime            string `query:"txnTime"`            // 订单发送时间
	AccessType         string `query:"accessType"`         // 接入类型
	MerId              string `query:"merId"`              // 商户代码
	OrderId            string `query:"orderId"`            // 商户订单号
	Reserved           string `query:"reserved"`           // 保留域
	ReqReserved        string `query:"reqReserved"`        // 请求方保留域
	AcqInsCode         string `query:"acqInsCode"`         // 收单机构代码
	PreAuthId          string `query:"preAuthId"`          // 预授权号
	InstalTransInfo    string `query:"instalTransInfo"`    // 分期付款信息域
}

type Revoke struct {
	Error
	TxnType     string `query:"txnType"`     // 交易类型
	TxnSubType  string `query:"txnSubType"`  // 交易子类
	BizType     string `query:"bizType"`     // 产品类型
	AccessType  string `query:"accessType"`  // 接入类型
	AcqInsCode  string `query:"acqInsCode"`  // 收单机构代码
	MerId       string `query:"merId"`       // 商户代码
	OrderId     string `query:"orderId"`     // 商户消费撤销订单号
	OrgQryId    string `query:"origQryId"`   // 原始交易流水号
	TxnTime     string `query:"txnTime"`     // 订单发送时间
	TxnAmt      string `query:"txnAmt"`      // 交易金额
	ReqReserved string `query:"reqReserved"` // 请求方保留域
	Reserved    string `query:"reserved"`    // 保留域
	QueryId     string `query:"queryId"`     // 银联交易流水号
}

type Refund struct {
	Error
	TxnType     string `query:"txnType"`     // 交易类型
	TxnSubType  string `query:"txnSubType"`  // 交易子类
	BizType     string `query:"bizType"`     // 产品类型
	AccessType  string `query:"accessType"`  // 接入类型
	AcqInsCode  string `query:"acqInsCode"`  // 收单机构代码
	MerId       string `query:"merId"`       // 商户代码
	OrderId     string `query:"orderId"`     // 商户退货订单号
	OrgQryId    string `query:"origQryId"`   // 原始交易流水号
	TxnTime     string `query:"txnTime"`     // 订单发送时间
	TxnAmt      string `query:"txnAmt"`      // 交易金额
	ReqReserved string `query:"reqReserved"` // 请求方保留域
	Reserved    string `query:"reserved"`    // 保留域
	QueryId     string `query:"queryId"`     // 银联交易流水号
}
