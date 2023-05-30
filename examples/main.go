package main

import (
	"fmt"
	"github.com/smartwalle/unionpay"
	"github.com/smartwalle/xid"
	"net/http"
)

// TODO 设置回调地址域名
const kServerDomain = "https://www.game2me.net"

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

	http.HandleFunc("/unionpay/web", func(writer http.ResponseWriter, request *http.Request) {
		var html, _ = client.CreateWebPayment(fmt.Sprintf("%d", xid.Next()), "100", kServerDomain+"/unionpay/front", kServerDomain+"/unionpay/back")
		writer.Write([]byte(html))
	})

	http.HandleFunc("/unionpay/app", func(writer http.ResponseWriter, request *http.Request) {
		var tn, _ = client.CreateAppPayment(fmt.Sprintf("%d", xid.Next()), "100", kServerDomain+"/union/back")
		writer.Write([]byte(tn))
	})

	http.HandleFunc("/unionpay/front", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		if err = client.VerifySign(request.Form); err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Good"))
	})

	http.HandleFunc("/unionpay/back", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		if err = client.VerifySign(request.Form); err != nil {
			fmt.Println("验证通知签名失败")
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("bad"))
			return
		}
		fmt.Println("验证通知签名成功")
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("ok"))
	})

	http.ListenAndServe(":9988", nil)
}
