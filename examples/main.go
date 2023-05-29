package main

import (
	"fmt"
	"github.com/smartwalle/unionpay"
	"github.com/smartwalle/xid"
	"net/http"
)

func main() {
	var client, err = unionpay.NewWithPFXFile("./acp_test_sign.pfx", "000000", "777290058165621", false)
	if err != nil {
		fmt.Println("初始化银联支付失败, 错误信息为", err)
		return
	}

	if err = client.LoadRootCertFromFile("./acp_test_root.cer"); err != nil {
		fmt.Println("加载证书发生错误", err)
		return
	}
	if err = client.LoadIntermediateCertFromFile("./acp_test_middle.cer"); err != nil {
		fmt.Println("加载证书发生错误", err)
		return
	}

	fmt.Println(client.Query("3619268148194181120"))

	http.HandleFunc("/pay/web", func(writer http.ResponseWriter, request *http.Request) {
		var html, _ = client.FrontTrans(fmt.Sprintf("%d", xid.Next()))
		writer.Write([]byte(html))
	})

	http.ListenAndServe(":9091", nil)
}
