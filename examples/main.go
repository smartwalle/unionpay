package main

import (
	"fmt"
	"github.com/smartwalle/unionpay"
	"github.com/smartwalle/xid"
	"net/http"
)

// TODO 设置回调地址域名
const kServerDomain = "http://127.0.0.1:9988"

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
		var orderId = fmt.Sprintf("%d", xid.Next())
		var html, txnTime, _ = client.CreateWebPayment(orderId, "100", kServerDomain+"/unionpay/front", kServerDomain+"/unionpay/back")
		writer.Write([]byte(html))
		fmt.Println(orderId, txnTime)
		fmt.Printf("%s/unionpay/query?order_id=%s&txn_time=%s \n", kServerDomain, orderId, txnTime)
	})

	http.HandleFunc("/unionpay/app", func(writer http.ResponseWriter, request *http.Request) {
		var orderId = fmt.Sprintf("%d", xid.Next())
		var tn, txnTime, _ = client.CreateAppPayment(orderId, "100", kServerDomain+"/union/back")
		writer.Write([]byte(tn))
		fmt.Println(orderId, txnTime)
		fmt.Printf("%s/unionpay/query?order_id=%s&txn_time=%s \n", kServerDomain, orderId, txnTime)
	})

	http.HandleFunc("/unionpay/query", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		var orderId = request.Form.Get("order_id")
		var txnTime = request.Form.Get("txn_time")

		var payment, err = client.GetTransaction(orderId, txnTime)
		if err != nil {
			fmt.Println("查询错误:", err)
			return
		}
		fmt.Printf("%v \n", payment)
	})

	http.HandleFunc("/unionpay/front", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		if err = client.VerifySign(request.Form); err != nil {
			writer.Write([]byte(err.Error()))
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Good"))

		fmt.Println(request.Form)
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
