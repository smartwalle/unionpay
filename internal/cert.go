package internal

import (
	"crypto/x509"
	"errors"
	"strings"
)

func VerifyCert(rootCert, intermediateCert, cert *x509.Certificate) error {
	var roots = x509.NewCertPool()
	roots.AddCert(rootCert)

	var intermediates = x509.NewCertPool()
	intermediates.AddCert(intermediateCert)
	intermediates.AddCert(rootCert)

	var opts = x509.VerifyOptions{
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		Intermediates: intermediates,
		Roots:         roots,
	}
	if _, err := cert.Verify(opts); err != nil {
		return err
	}

	var commons = strings.Split(cert.Subject.CommonName, "@")
	if len(commons) < 2 || commons[2] != "中国银联股份有限公司" {
		return errors.New("invalid certificate")
	}
	return nil
}
