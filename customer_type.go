package unionpay

type Customer struct {
	CertType string // 证件类型 01：身份证  02：军官证 03：护照 04：回乡证 05：台胞证 06：警官证 07：士兵证 99：其它证件
	CertId   string // 证件号码
	Name     string // 姓名
	SMSCode  string // 短信验证码
	PIN      string // 持卡人密码
	CVN2     string // CVN2
	Expired  string // 有效期
	PhoneNo  string // 手机号
}
