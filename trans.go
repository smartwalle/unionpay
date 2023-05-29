package unionpay

import (
	"bytes"
	"html/template"
	"net/url"
)

const (
	kFrontTrans = "/gateway/api/frontTransReq.do"
)

func (this *Client) FrontTrans(orderId string) (string, error) {
	var values = url.Values{}
	values.Set("orderId", orderId)
	values.Set("currencyCode", "156")
	values.Set("txnAmt", "156")

	values.Set("frontUrl", "156")
	values.Set("backUrl", "156")

	values.Set("bizType", "000201")
	values.Set("txnType", "01")
	values.Set("txnSubType", "01")
	values.Set("accessType", "0")
	values.Set("channelType", "07")

	values, err := this.URLValues(values)
	if err != nil {
		return "", err
	}

	var buff = bytes.NewBufferString("")

	tpl, err := template.New("").Parse(kFrontTransTemplate)
	tpl.Execute(buff, map[string]interface{}{
		"Values": values,
		"Action": this.host + kFrontTrans,
	})
	return buff.String(), nil
}
