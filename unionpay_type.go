package unionpay

import "fmt"

const (
	kSandboxGateway    = "https://gateway.test.95516.com"
	kProductionGateway = "https://gateway.95516.com"

	kVersion    = "5.1.0"
	kSignMehtod = "01"
)

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
	Code Code
	Msg  string
}

func (this Error) Error() string {
	return fmt.Sprintf("%s - %s", this.Code, this.Msg)
}

func (this Error) IsSuccess() bool {
	return this.Code.IsSuccess()
}

func (this Error) IsFailure() bool {
	return this.Code.IsFailure()
}
