package unionpay

import (
	"github.com/smartwalle/nhttp"
	"net/url"
)

var mapper = nhttp.NewMapper("query")

func DecodeValues(values url.Values, dst interface{}) error {
	return mapper.Bind(values, dst)
}
