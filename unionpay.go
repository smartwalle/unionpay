package unionpay

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/smartwalle/ncrypto"
	"github.com/smartwalle/ncrypto/pkcs12"
	"github.com/smartwalle/ngx"
	"github.com/smartwalle/nsign"
	"github.com/smartwalle/unionpay/internal"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"text/template"
	"time"
)

type Signer interface {
	SignValues(values url.Values, opts ...nsign.SignOption) ([]byte, error)

	SignBytes(data []byte, opts ...nsign.SignOption) ([]byte, error)
}

type Verifier interface {
	VerifyValues(values url.Values, signature []byte, opts ...nsign.SignOption) error

	VerifyBytes(data []byte, signature []byte, opts ...nsign.SignOption) error
}

type OptionFunc func(c *Client)

func WithHTTPClient(client *http.Client) OptionFunc {
	return func(c *Client) {
		if client != nil {
			c.Client = client
		}
	}
}

type Client struct {
	Client     *http.Client
	host       string
	merchantId string
	certId     string

	version    string
	signMethod string

	webPaymentTpl *template.Template

	rootCert  *x509.Certificate
	interCert *x509.Certificate

	// 签名和验签
	mu        sync.Mutex
	signer    Signer
	verifiers map[string]Verifier

	// 敏感信息加密&解密
	decryptPrivateKey *rsa.PrivateKey
	encryptPublicKey  *rsa.PublicKey
	encryptCertId     string
}

// New 初始银联客户端
//
// pfx - 商户私钥证书
//
// password - 商户私钥证书密码
//
// merchantId - 商户号
//
// isProduction - 是否为生产环境，传 false 的时候为沙箱环境，用于开发测试，正式上线的时候需要改为 true
func New(pfx []byte, password, merchantId string, isProduction bool, opts ...OptionFunc) (*Client, error) {
	rawKey, certificate, _, err := pkcs12.Decode(pfx, password)
	if err != nil {
		return nil, err
	}

	privateKey, _ := rawKey.(*rsa.PrivateKey)
	if privateKey == nil {
		return nil, errors.New("key is not a valid *rsa.PrivateKey")
	}

	var nClient = &Client{}
	if err = nClient.LoadWebPaymentTemplate(kWebPaymentTemplate); err != nil {
		return nil, err
	}

	nClient.Client = http.DefaultClient
	if isProduction {
		nClient.host = kProductionGateway
	} else {
		nClient.host = kSandboxGateway
	}
	nClient.merchantId = merchantId
	nClient.certId = certificate.SerialNumber.String()

	nClient.version = kVersion
	nClient.signMethod = kSignMethod

	nClient.signer = nsign.New(nsign.WithMethod(internal.NewRSAMethod(crypto.SHA256, privateKey, nil)))
	nClient.verifiers = make(map[string]Verifier)

	nClient.decryptPrivateKey = privateKey

	for _, opt := range opts {
		if opt != nil {
			opt(nClient)
		}
	}

	return nClient, nil
}

// NewWithPFXFile 初始银联客户端
//
// filename - 商户私钥证书文件
//
// password - 商户私钥证书密码
//
// merchantId - 商户号
//
// isProduction - 是否为生产环境，传 false 的时候为沙箱环境，用于开发测试，正式上线的时候需要改为 true
func NewWithPFXFile(filename, password, merchantId string, isProduction bool) (*Client, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return New(data, password, merchantId, isProduction)
}

// LoadWebPaymentTemplate 用于加载跳转银联支付页面的网页模版。
//
// 网页支付需要先在浏览器中打开业务方(商户)提供的网页，通过该网页跳转到银联的支付页面。
//
// CreateWebPayment 方法中会构建相应的参数，然后把本方法加载的模版渲染成 HTML 代码。
//
// 模版参考 unionpay_type.go 文件中的 kWebPaymentTemplate 常量，该常量也是本库默认使用的模版。
func (c *Client) LoadWebPaymentTemplate(tpl string) error {
	nTemplate, err := template.New("").Parse(tpl)
	if err != nil {
		return err
	}
	c.webPaymentTpl = nTemplate
	return nil
}

func (c *Client) loadRootCert(b []byte) error {
	cert, err := ncrypto.DecodeCertificate(b)
	if err != nil {
		return err
	}
	c.rootCert = cert
	return nil
}

// LoadRootCert 加载银联根证书
func (c *Client) LoadRootCert(s string) error {
	return c.loadRootCert([]byte(s))
}

// LoadRootCertFromFile 从文件加载银联根证书
func (c *Client) LoadRootCertFromFile(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return c.loadRootCert(b)
}

func (c *Client) loadIntermediateCert(b []byte) error {
	cert, err := ncrypto.DecodeCertificate(b)
	if err != nil {
		return err
	}
	c.interCert = cert
	return nil
}

// LoadIntermediateCert 加载银联中间证书
func (c *Client) LoadIntermediateCert(s string) error {
	return c.loadIntermediateCert([]byte(s))
}

// LoadIntermediateCertFromFile 从文件加载银联中间证书
func (c *Client) LoadIntermediateCertFromFile(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return c.loadIntermediateCert(b)
}

// LoadEncryptKey 银联加密公钥更新查询接口（敏感加密证书）。
//
// 商户定期（1天1次）向银联全渠道系统发起获取加密公钥信息交易。在加密公钥证书更新期间，全渠道系统支持新老证书的共同使用，新老证书并行期为1个月。全渠道系统向商户返回最新的加密公钥证书，由商户服务器替换本地证书。
//
// 文档地址：https://open.unionpay.com/tjweb/acproduct/APIList?acpAPIId=758&apiservId=448&version=V2.2&bussType=0
func (c *Client) LoadEncryptKey() error {
	var values = url.Values{}
	values.Set("accessType", "0")
	values.Set("channelType", "07") // 渠道类型
	values.Set("txnType", "95")     // 交易类型 95-银联加密公钥更新查询
	values.Set("txnSubType", "00")  // 交易子类型 默认00
	values.Set("bizType", "000000") // 业务类型  默认
	values.Set("certType", "01")    // 01：敏感信息加密公钥(只有01可用)
	values.Set("orderId", time.Now().Format("20060102150405"))
	values.Set("txnTime", time.Now().Format("20060102150405"))

	var rValues, err = c.Request(kBackTrans, values)
	if err != nil {
		return err
	}
	var cert = strings.ReplaceAll(rValues.Get("encryptPubKeyCert"), "\r", "\n")

	certificate, err := c.decodeCertificate([]byte(cert))
	if err != nil {
		return err
	}
	c.encryptPublicKey, _ = certificate.PublicKey.(*rsa.PublicKey)
	c.encryptCertId = certificate.SerialNumber.String()
	return nil
}

// LoadEncryptKeyFromFile 从文件加载银联敏感加密证书。
func (c *Client) LoadEncryptKeyFromFile(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	certificate, err := c.decodeCertificate(b)
	if err != nil {
		return err
	}
	c.encryptPublicKey, _ = certificate.PublicKey.(*rsa.PublicKey)
	c.encryptCertId = certificate.SerialNumber.String()
	return nil
}

func (c *Client) decodeCertificate(b []byte) (*x509.Certificate, error) {
	certificate, err := ncrypto.DecodeCertificate(b)
	if err != nil {
		return nil, err
	}
	return certificate, nil
}

func (c *Client) URLValues(values url.Values) (url.Values, error) {
	if values == nil {
		values = url.Values{}
	}

	values.Set("version", c.version)
	values.Set("encoding", "UTF-8")
	values.Set("merId", c.merchantId)
	values.Set("certId", c.certId)
	values.Set("signMethod", c.signMethod)

	signature, err := c.signer.SignValues(values)
	if err != nil {
		return nil, err
	}
	values.Set("signature", base64.StdEncoding.EncodeToString(signature))
	return values, nil
}

func (c *Client) Request(api string, values url.Values) (url.Values, error) {
	values, err := c.URLValues(values)
	if err != nil {
		return nil, err
	}

	var req = ngx.NewRequest(http.MethodPost, c.host+api, ngx.WithClient(c.Client))
	req.SetForm(values)

	rsp, err := req.Do(context.Background())
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	// 解析返回数据
	rValues, err := internal.ParseQuery(string(data))
	if err != nil {
		return nil, err
	}

	// 验证签名
	if err = c.VerifySign(rValues); err != nil {
		return nil, err
	}

	return rValues, nil
}

func (c *Client) VerifySign(values url.Values) error {
	verifier, err := c.getVerifier(values.Get("signPubKeyCert"))
	if err != nil {
		return err
	}

	signature, err := base64.StdEncoding.DecodeString(values.Get("signature"))
	if err != nil {
		return err
	}

	return verifier.VerifyValues(values, signature, nsign.WithIgnore("signature"))
}

func (c *Client) getVerifier(cert string) (Verifier, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var verifier = c.verifiers[cert]
	if verifier == nil {
		certificate, err := ncrypto.DecodeCertificate([]byte(cert))
		if err != nil {
			return nil, err
		}

		if err = internal.VerifyCert(c.rootCert, c.interCert, certificate); err != nil {
			return nil, err
		}

		verifier = nsign.New(nsign.WithMethod(internal.NewRSAMethod(crypto.SHA256, nil, certificate.PublicKey.(*rsa.PublicKey))))
		c.verifiers[cert] = verifier
	}
	return verifier, nil
}

// Decrypt 用于解密从银联获取到的敏感信息。
//
// 如果商户号开通了【商户对敏感信息加密】的权限，那么需要对获取到的 accNo、pin、phoneNo、cvn2、expired 进行解密。
//
// 如果商户号未开通【商户对敏感信息加密】权限，那么不需要对敏感信息进行解密。
//
// https://open.unionpay.com/tjweb/support/faq/mchlist?id=537
func (c *Client) Decrypt(s string) (string, error) {
	var ciphertext, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", nil
	}

	ciphertext, err = ncrypto.RSADecrypt(ciphertext, c.decryptPrivateKey)
	if err != nil {
		return "", nil
	}
	return string(ciphertext), nil
}

// EncryptCertId 获取敏感信息加密证书 id。
//
// 用于各接口中的 encryptCertId 字段。
func (c *Client) EncryptCertId() string {
	return c.encryptCertId
}

// Encrypt 对数据进行加密，并对加密的结果使用 base64 进行编码，（用于加密敏感信息）。
//
// 如果商户号开通了【商户对敏感信息加密】的权限，那么需要对提交的 accNo、pin、phoneNo、cvn2、expired 进行加密。
//
// 如果商户号未开通【商户对敏感信息加密】权限，那么不需要对敏感信息进行加密。
//
// https://open.unionpay.com/tjweb/support/faq/mchlist?id=537
func (c *Client) Encrypt(s string) (string, error) {
	return c.EncryptBytes([]byte(s))
}

func (c *Client) EncryptBytes(b []byte) (string, error) {
	if c.encryptPublicKey == nil || c.encryptCertId == "" {
		return "", errors.New("public key not found, you need to call LoadEncryptKey() first")
	}

	var ciphertext, err = ncrypto.RSAEncrypt(b, c.encryptPublicKey)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// EncryptPIN 对 PIN 进行加密，并对加密的结果使用 base64 进行编码。
func (c *Client) EncryptPIN(pan, pin string) (string, error) {
	return c.EncryptBytes(PINBlock(pan, pin))
}

// PINBlock https://paymentcardtools.com/pin-block-calculators/iso9564-format-0
func PINBlock(pan, pin string) []byte {
	pan = "0000" + string(pan[len(pan)-13:len(pan)-1])
	pin = fmt.Sprintf("0%d%s", len(pin), pin)

	var blockLen = 8

	var pinBytes, _ = hex.DecodeString(pin)
	for i := blockLen / 2; i < blockLen; i++ {
		pinBytes = append(pinBytes, 0xff)
	}

	var panBytes, _ = hex.DecodeString(pan)

	var block = make([]byte, 8)
	for i := 0; i < blockLen; i++ {
		block[i] = pinBytes[i] ^ panBytes[i]
	}
	return block
}
