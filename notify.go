package unionpay

import (
	"errors"
	"fmt"
	"github.com/smartwalle/unionpay/internal"
	"net/http"
)

// DecodeNotification 解析通知
//
// 返回以下类型中的一种：
//
// *PaymentNotification
//
// *RevokeNotification
//
// *RefundNotification
func (this *Client) DecodeNotification(req *http.Request) (interface{}, error) {
	if req == nil {
		return nil, errors.New("request 参数不能为空")
	}
	if err := req.ParseForm(); err != nil {
		return nil, err
	}

	if err := this.VerifySign(req.Form); err != nil {
		return nil, err
	}

	fmt.Println(req.Form)

	var txnType = req.Form.Get("txnType")
	switch txnType {
	case "01":
		var notification *PaymentNotification
		if err := internal.DecodeValues(req.Form, &notification); err != nil {
			return nil, err
		}
		return notification, nil
	case "04":
		var notification *RevokeNotification
		if err := internal.DecodeValues(req.Form, &notification); err != nil {
			return nil, err
		}
		return notification, nil
	case "31":
		var notification *RefundNotification
		if err := internal.DecodeValues(req.Form, &notification); err != nil {
			return nil, err
		}
		return notification, nil
	}

	return nil, nil
}
