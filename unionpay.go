package unionpay

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"github.com/smartwalle/ncrypto"
	"github.com/smartwalle/ncrypto/pkcs12"
	"github.com/smartwalle/ngx"
	"github.com/smartwalle/nsign"
	"github.com/smartwalle/unionpay/internal"
	"io"
	"net/http"
	"net/url"
	"os"
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

	rootCert  *x509.Certificate
	interCert *x509.Certificate

	// 签名和验签
	signer    Signer
	verifiers map[string]Verifier
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
	nClient.Client = http.DefaultClient
	if isProduction {
		nClient.host = kProductionGateway
	} else {
		nClient.host = kSandboxGateway
	}
	nClient.merchantId = merchantId
	nClient.certId = certificate.SerialNumber.String()

	nClient.version = kVersion
	nClient.signMethod = kSignMehtod

	nClient.signer = nsign.New(nsign.WithMethod(internal.NewRSAMethod(crypto.SHA256, privateKey, nil)))
	nClient.verifiers = make(map[string]Verifier)

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

func (this *Client) loadRootCert(b []byte) error {
	cert, err := ncrypto.DecodeCertificate(b)
	if err != nil {
		return nil
	}
	this.rootCert = cert
	return nil
}

// LoadRootCert 加载银联根证书
func (this *Client) LoadRootCert(s string) error {
	return this.loadRootCert([]byte(s))
}

// LoadRootCertFromFile 从文件加载银联根证书
func (this *Client) LoadRootCertFromFile(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return this.loadRootCert(b)
}

func (this *Client) loadIntermediateCert(b []byte) error {
	cert, err := ncrypto.DecodeCertificate(b)
	if err != nil {
		return nil
	}
	this.interCert = cert
	return nil
}

// LoadIntermediateCert 加载银联中间证书
func (this *Client) LoadIntermediateCert(s string) error {
	return this.loadIntermediateCert([]byte(s))
}

// LoadIntermediateCertFromFile 从文件加载银联中间证书
func (this *Client) LoadIntermediateCertFromFile(filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return this.loadIntermediateCert(b)
}

func (this *Client) URLValues(values url.Values) (url.Values, error) {
	if values == nil {
		values = url.Values{}
	}

	values.Set("version", this.version)
	values.Set("encoding", "UTF-8")
	values.Set("txnTime", time.Now().Format("20060102150405"))
	values.Set("merId", this.merchantId)
	values.Set("certId", this.certId)
	values.Set("signMethod", this.signMethod)

	signature, err := this.signer.SignValues(values)
	if err != nil {
		return nil, err
	}
	values.Set("signature", base64.StdEncoding.EncodeToString(signature))
	return values, nil
}

func (this *Client) request(method, api string, values url.Values) (url.Values, error) {
	values, err := this.URLValues(values)
	if err != nil {
		return nil, err
	}

	var req = ngx.NewRequest(method, this.host+api, ngx.WithClient(this.Client))
	req.SetParams(values)

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

	var cert = rValues.Get("signPubKeyCert")

	// 解析证书
	certificate, err := ncrypto.DecodeCertificate([]byte(cert))
	if err != nil {
		return nil, err
	}

	if err = internal.VerifyCert(this.rootCert, this.interCert, certificate); err != nil {
		return nil, err
	}

	// 验证签名
	var verifier = nsign.New(nsign.WithMethod(internal.NewRSAMethod(crypto.SHA256, nil, certificate.PublicKey.(*rsa.PublicKey))))

	signature, err := base64.StdEncoding.DecodeString(rValues.Get("signature"))
	if err != nil {
		return nil, err
	}

	if err = verifier.VerifyValues(rValues, signature, nsign.WithIgnore("signature")); err != nil {
		return nil, err
	}

	return rValues, nil
}
