银联支付

## 鸣谢

[![jetbrains.svg](jetbrains.svg)](https://www.jetbrains.com/?from=AliPay%20SDK%20for%20Go)

## 帮助

在集成的过程中有遇到问题，欢迎加 QQ 群 203357977 讨论。

## 安装

```go
go get github.com/smartwalle/unionpay
```

```go
import github.com/smartwalle/unionpay
```

#### 关于错误 x509: certificate signed by unknown authority (possibly because of "x509: cannot verify signature: insecure algorithm SHA1-RSA (temporarily override with GODEBUG=x509sha1=1)" while trying to verify candidate authority certificate "CFCA TEST OCA1")

由于安全问题，从 Go 1.18 开始，其 crypto/x509 包默认将拒绝使用 SHA-1 哈希函数签名的证书(自签发的除外)。

银联测试环境签发的证书目前使用的是 SHA-1 签名（没有生产环境，无法得知生产环境信息是否也是使用 SHA-1 进行签名）。

#### 解决办法
可以通过设置 GODEBUG=x509sha1=1 环境变量暂时恢复。

```shell
GODEBUG=x509sha1=1 go run main.go 
```

## 其它支付

支付宝 [https://github.com/smartwalle/alipay](https://github.com/smartwalle/alipay)

苹果支付 [https://github.com/smartwalle/apple](https://github.com/smartwalle/apple)

PayPal [https://github.com/smartwalle/paypal](https://github.com/smartwalle/paypal)

## 测试环境

可从 [https://open.unionpay.com/tjweb/ij/user/mchTest/param](https://open.unionpay.com/tjweb/ij/user/mchTest/param) 获取测试环境信息。

## 示例

[https://github.com/smartwalle/unionpay/tree/master/examples](https://github.com/smartwalle/unionpay/tree/master/examples) 包含一个完整示例。

运行该示例代码之后，可以在浏览器中访问 [http://127.0.0.1:9988/unionpay](http://127.0.0.1:9988/unionpay) 以打开测试页面。

## 已实现接口

* 消费接口-创建网页支付 - CreateWebPayment()
* 消费接口-创建 App 支付 - CreateAppPayment()
* 交易状态查询接口 - GetTransaction()
* 消费撤销接口 - Revoke()
* 退货接口接口 - Refund()

## 关于交易状态

在银联系统中，发起消费(支付)、消费撤销(退款)和退货(退款)都会产生交易，都可以通过[交易状态查询接口](https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=757&apiservId=448&version=V2.2&bussType=0)查询相关信息。

**退款状态不会在原支付交易中体现。**

## 重要资料

* [如何判断交易成功？怎么确定交易成功？](https://open.unionpay.com/tjweb/support/faq/mchlist?id=116)
* [应答码及交易状态查询机制
  ](https://open.unionpay.com/tjweb/support/faq/mchlist?id=234)
* [交易状态查询交易，应答的respCode和origRespCode有什么区别？](https://open.unionpay.com/tjweb/support/faq/mchlist?id=610)
* [消费撤销/退货成功了，查询接口能查到它们的状态吗？如何查询退款交易？有统一交易查询接口么？
  ](https://open.unionpay.com/tjweb/support/faq/mchlist?id=79)
* [接收通知后我们需要返什么应答？接收异步通知后需要返什么给银联？异步通知需要返回应答是什么？](https://open.unionpay.com/tjweb/support/faq/mchlist?id=72)
* [交易状态查询接口什么时候调用？](https://open.unionpay.com/tjweb/support/faq/mchlist?id=77)
