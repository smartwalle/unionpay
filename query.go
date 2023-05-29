package unionpay

import (
	"net/http"
	"net/url"
)

const (
	kQueryTrans = "/gateway/api/queryTrans.do"
)

// Query https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=757&apiservId=448&version=V2.2&bussType=0
func (this *Client) Query(orderId string) (interface{}, error) {
	var values = url.Values{}
	values.Set("bizType", "000000")
	values.Set("txnType", "00")
	values.Set("txnSubType", "00")
	values.Set("accessType", "0")
	values.Set("orderId", orderId)
	return this.request(http.MethodPost, kQueryTrans, values)
}
