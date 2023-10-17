package unionpay

import (
	"fmt"
	"net/url"
)

const (
	kSandboxGateway    = "https://gateway.test.95516.com"
	kProductionGateway = "https://gateway.95516.com"

	kVersion    = "5.1.0"
	kSignMethod = "01"
)

const kWebPaymentTemplate = `
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" /></head>
<body>
<form id="pay_form" action="{{.Action}}" method="POST">
{{range $k, $v := .Values}}
<input type="hidden" name="{{$k}}" id="{{$k}}" value="{{index $v 0}}" />
{{end}}
</form>
<script type="text/javascript">
document.getElementById("pay_form").submit();
</script>
</body>
</html>
`

type Code string

func (c Code) IsSuccess() bool {
	return c == CodeSuccess
}

func (c Code) IsFailure() bool {
	return c != CodeSuccess
}

const (
	CodeSuccess Code = "00" // 接口调用成功
)

type Error struct {
	Code Code   `query:"respCode"`
	Msg  string `query:"respMsg"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%s - %s", e.Code, e.Msg)
}

func (e Error) IsSuccess() bool {
	return e.Code.IsSuccess()
}

func (e Error) IsFailure() bool {
	return e.Code.IsFailure()
}

type Payload struct {
	values url.Values
}

func NewPayload() *Payload {
	var nPayload = &Payload{}
	nPayload.values = url.Values{}
	return nPayload
}

func (p *Payload) AddParam(key, value string) *Payload {
	if key != "" && value != "" {
		p.values.Set(key, value)
	}
	return p
}

type CallOption func(values url.Values)

func WithPayload(payload *Payload) CallOption {
	return func(values url.Values) {
		if payload != nil {
			for key := range payload.values {
				values[key] = payload.values[key]
			}
		}
	}
}
