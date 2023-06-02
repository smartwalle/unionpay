package main

import (
	"encoding/json"
	"fmt"
	"github.com/smartwalle/unionpay"
	"github.com/smartwalle/xid"
	"log"
	"net/http"
)

// TODO 设置回调地址域名
const kServerDomain = "http://127.0.0.1:9988"

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	var client, err = unionpay.NewWithPFXFile("./acp_test_sign.pfx", "000000", "777290058165621", false)
	if err != nil {
		log.Println("初始化银联支付失败, 错误信息为", err)
		return
	}

	if err = client.LoadRootCertFromFile("./acp_test_root.cer"); err != nil {
		log.Println("加载证书发生错误", err)
		return
	}
	if err = client.LoadIntermediateCertFromFile("./acp_test_middle.cer"); err != nil {
		log.Println("加载证书发生错误", err)
		return
	}

	// 加载敏感信息加密密钥
	if err = client.LoadEncryptKey(); err != nil {
		log.Println("加载敏感信息加密密钥发生错误", err)
		return
	}

	http.HandleFunc("/unionpay", func(writer http.ResponseWriter, request *http.Request) {
		var html = `
<html>
<head>
</head>
<body>
<a href="/unionpay/web">web</a>
<a href="/unionpay/app">app</a>
</body>
</html>`
		writer.Write([]byte(html))
	})

	http.HandleFunc("/unionpay/web", func(writer http.ResponseWriter, request *http.Request) {
		var payment, err = client.CreateWebPayment(fmt.Sprintf("%d", xid.Next()), "100", kServerDomain+"/unionpay/front", kServerDomain+"/unionpay/back")
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		writer.Write([]byte(payment.HTML))
		log.Printf("查询交易状态：%s/unionpay/query?order_id=%s&txn_time=%s \n", kServerDomain, payment.OrderId, payment.TxnTime)
	})

	http.HandleFunc("/unionpay/app", func(writer http.ResponseWriter, request *http.Request) {
		var payment, err = client.CreateAppPayment(fmt.Sprintf("%d", xid.Next()), "100", kServerDomain+"/union/back")
		if err != nil {
			writer.Write([]byte(err.Error()))
			return
		}

		var data, _ = json.Marshal(payment)
		writer.Write(data)
		log.Printf("查询交易状态：%s/unionpay/query?order_id=%s&txn_time=%s \n", kServerDomain, payment.OrderId, payment.TxnTime)
	})

	http.HandleFunc("/unionpay/query", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		var orderId = request.Form.Get("order_id")
		var txnTime = request.Form.Get("txn_time")

		var transaction, err = client.GetTransaction(orderId, txnTime)
		if err != nil {
			log.Println("查询错误:", err)
			return
		}

		var data, _ = json.Marshal(transaction)
		writer.Write(data)
	})

	http.HandleFunc("/unionpay/front", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		var notification, err = client.DecodeNotification(request.Form)
		if err != nil {
			log.Println("解析前台通知发生错误：", err)
			fmt.Fprintln(writer, "解析前台通知发生错误：", err)
			return
		}

		switch nType := notification.(type) {
		case *unionpay.PaymentNotification:
			// 测试对敏感信息解密
			accNo, err := client.Decrypt(nType.AccNo)
			if err == nil {
				nType.AccNo = accNo
			}
		}

		writer.WriteHeader(http.StatusOK)
		var data, _ = json.Marshal(notification)
		writer.Write(data)
	})

	http.HandleFunc("/unionpay/back", func(writer http.ResponseWriter, request *http.Request) {
		request.ParseForm()

		var notification, err = client.DecodeNotification(request.Form)
		if err != nil {
			log.Println("解析后台通知失败：", err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println("解析后台通知成功")
		log.Println(notification)

		client.ACKNotification(writer)
	})

	log.Println(kServerDomain + "/unionpay")
	http.ListenAndServe(":9988", nil)
}
