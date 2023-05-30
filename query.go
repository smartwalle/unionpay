package unionpay

import "net/url"

const (
	kQueryTrans = "/gateway/api/queryTrans.do"
)

// Query 交易状态查询接口 https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=757&apiservId=448&version=V2.2&bussType=0
func (this *Client) Query(orderId string, opts ...CallOption) (map[string]string, error) {
	var values = url.Values{}
	values.Set("accessType", "0")
	for _, opt := range opts {
		if opt != nil {
			opt(values)
		}
	}
	values.Set("bizType", "000000")
	values.Set("txnType", "00")
	values.Set("txnSubType", "00")
	values.Set("orderId", orderId)

	var rValues, err = this.Request(kQueryTrans, values)
	if err != nil {
		return nil, err
	}

	var result = make(map[string]string)
	for key := range rValues {
		result[key] = rValues.Get(key)
	}
	return result, nil
}
