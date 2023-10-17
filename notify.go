package unionpay

import (
	"fmt"
	"net/http"
	"net/url"
)

// DecodeNotification 解析通知。
//
// 各通知结构体会尽量包含前台通知(frontURL)和后台通知(backURL)已知字段，所以本方法可用于解析前台通知(frontURL)和后台通知(backURL)。
//
// 返回以下类型中的一种：
//
// *PaymentNotification
//
// *RevokeNotification
//
// *RefundNotification
func (c *Client) DecodeNotification(values url.Values) (interface{}, error) {
	if err := c.VerifySign(values); err != nil {
		return nil, err
	}

	var txnType = values.Get("txnType")
	switch txnType {
	case "01":
		return DecodePaymentNotification(values)
	case "04":
		return DecodeRevokeNotification(values)
	case "31":
		return DecodeRefundNotification(values)
	}

	return nil, fmt.Errorf("unknown txnType %s", txnType)
}

func (c *Client) ACKNotification(w http.ResponseWriter) {
	ACKNotification(w)
}

// ACKNotification 返回异步通知成功处理的消息给银联。
//
// 后台通知以标准的HTTP协议的POST方法向商户的后台通知URL发送，超时时间为10秒。
//
// 由于网络等原因，商户可能会收到重复的后台通知，商户应能正确识别并处理。
//
// 商户返回码为200时，银联判定为通知成功，其他返回码为通知失败。
//
// 如10秒内未收到应答，银联判定为通知失败。
//
// 第一次通知失败后，银联会重发，最多发送五次（间隔1、2、4、5分钟）。
func ACKNotification(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func DecodePaymentNotification(values url.Values) (*PaymentNotification, error) {
	var notification *PaymentNotification
	if err := DecodeValues(values, &notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func DecodeRevokeNotification(values url.Values) (*RevokeNotification, error) {
	var notification *RevokeNotification
	if err := DecodeValues(values, &notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func DecodeRefundNotification(values url.Values) (*RefundNotification, error) {
	var notification *RefundNotification
	if err := DecodeValues(values, &notification); err != nil {
		return nil, err
	}
	return notification, nil
}
