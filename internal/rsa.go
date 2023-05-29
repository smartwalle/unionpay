package internal

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
)

type RSAMethod struct {
	h          crypto.Hash
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewRSAMethod(h crypto.Hash, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *RSAMethod {
	var nRSA = &RSAMethod{}
	nRSA.h = h
	nRSA.privateKey = privateKey
	nRSA.publicKey = publicKey
	return nRSA
}

func (this *RSAMethod) hash(data []byte) ([]byte, error) {
	var h = this.h.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	var hashed = h.Sum(nil)
	return hashed, nil
}

func (this *RSAMethod) Sign(data []byte) ([]byte, error) {
	hashed, err := this.hash(data)
	if err != nil {
		return nil, err
	}

	hashed, err = this.hash([]byte(hex.EncodeToString(hashed)))
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, this.privateKey, this.h, hashed)
}

func (this *RSAMethod) Verify(data []byte, signature []byte) error {
	hashed, err := this.hash(data)
	if err != nil {
		return err
	}

	hashed, err = this.hash([]byte(fmt.Sprintf("%x", hashed)))
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(this.publicKey, this.h, hashed, signature)
}
