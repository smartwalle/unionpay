package unionpay

const (
	kQueryTrans = "/gateway/api/queryTrans.do"
)

// Query 交易状态查询接口 https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=757&apiservId=448&version=V2.2&bussType=0
func (this *Client) Query(orderId string) (interface{}, error) {
	var payload = NewPayload(kQueryTrans)
	payload.Set("bizType", "000000")
	payload.Set("txnType", "00")
	payload.Set("txnSubType", "00")
	payload.Set("accessType", "0")
	payload.Set("orderId", orderId)
	return this.Request(payload)
}
