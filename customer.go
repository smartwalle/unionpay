package unionpay

import (
	"encoding/base64"
	"github.com/smartwalle/unionpay/internal"
	"net/url"
)

// EncodeCustomer 对 Customer 进行编码，不会对 phoneNo，cvn2，expired 进行加密。
//func (this *Client) EncodeCustomer(customer *Customer, accNo string) (string, error) {
//	if customer == nil {
//		return "", nil
//	}
//
//	var v = url.Values{}
//	if customer.CertType != "" {
//		v.Set("certifTp", customer.CertType)
//	}
//	if customer.CertId != "" {
//		v.Set("certifId", customer.CertId)
//	}
//	if customer.Name != "" {
//		v.Set("customerNm", customer.Name)
//	}
//	if customer.SMSCode != "" {
//		v.Set("smsCode", customer.SMSCode)
//	}
//	if customer.PIN != "" {
//		var block, err = this.EncryptPIN(accNo, customer.PIN)
//		if err != nil {
//			return "", err
//		}
//		v.Set("pin", block)
//	}
//	if customer.CVN2 != "" {
//		v.Set("cvn2", customer.CVN2)
//	}
//	if customer.Expired != "" {
//		v.Set("expired", customer.Expired)
//	}
//	if customer.PhoneNo != "" {
//		v.Set("phoneNo", customer.PhoneNo)
//	}
//
//	var r = internal.EncodeValues(v)
//	if r != "" {
//		r = "{" + r + "}"
//	}
//
//	return base64.StdEncoding.EncodeToString([]byte(r)), nil
//}

// EncryptCustomer 对 Customer 进行编码，会对 phoneNo，cvn2，expired 进行加密。
func (c *Client) EncryptCustomer(customer *Customer, accNo string) (string, error) {
	if customer == nil {
		return "", nil
	}

	var v = url.Values{}
	if customer.CertType != "" {
		v.Set("certifTp", customer.CertType)
	}
	if customer.CertId != "" {
		v.Set("certifId", customer.CertId)
	}
	if customer.Name != "" {
		v.Set("customerNm", customer.Name)
	}
	if customer.SMSCode != "" {
		v.Set("smsCode", customer.SMSCode)
	}
	if customer.PIN != "" {
		var block, err = c.EncryptPIN(accNo, customer.PIN)
		if err != nil {
			return "", err
		}
		v.Set("pin", block)
	}

	var ev = url.Values{}
	if customer.CVN2 != "" {
		ev.Set("cvn2", customer.CVN2)
	}
	if customer.Expired != "" {
		ev.Set("expired", customer.Expired)
	}
	if customer.PhoneNo != "" {
		ev.Set("phoneNo", customer.PhoneNo)
	}

	var evs = internal.EncodeValues(ev)
	if evs != "" {
		encryptedInfo, err := c.Encrypt(evs)
		if err != nil {
			return "", err
		}
		v.Set("encryptedInfo", encryptedInfo)
	}

	var r = internal.EncodeValues(v)
	if r != "" {
		r = "{" + r + "}"
	}

	return base64.StdEncoding.EncodeToString([]byte(r)), nil
}
