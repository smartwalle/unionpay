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
		var html, _ = client.CreateWebPayment(fmt.Sprintf("%d", xid.Next()), "100", "http://127.0.0.1:9091/pay/front", "http://127.0.0.1:9091/pay/back")
		writer.Write([]byte(html))
	})

	http.HandleFunc("/pay/front", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		if err = client.VerifySign(request.Form); err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Write([]byte("Good"))
	})

	http.ListenAndServe(":9091", nil)
}
