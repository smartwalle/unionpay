package unionpay

import (
	"github.com/smartwalle/nhttp"
	"github.com/smartwalle/unionpay/internal"
	"net/url"
)

var mapper = nhttp.NewMapper("query")

func DecodeValues(values url.Values, dst interface{}) error {
	return mapper.Bind(values, dst)
}

func EncodeValues(values url.Values) string {
	return internal.EncodeValues(values)
}
