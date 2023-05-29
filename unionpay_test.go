package unionpay_test

import (
	"fmt"
	"github.com/smartwalle/unionpay"
	"os"
)

var client *unionpay.Client

func init() {
	var err error
	client, err = unionpay.NewWithPFXFile("./acp_test_sign.pfx", "000000", "777290058165621", false)

	if err != nil {
		fmt.Println("初始化银联支付失败, 错误信息为", err)
		os.Exit(-1)
	}

	fmt.Println(client.LoadRootCertFromFile("./acp_test_root.cer"))
	fmt.Println(client.LoadIntermediateCertFromFile("./acp_test_middle.cer"))
}
