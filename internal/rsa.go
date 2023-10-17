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

func (m *RSAMethod) hash(data []byte) ([]byte, error) {
	var h = m.h.New()
	if _, err := h.Write(data); err != nil {
		return nil, err
	}
	var hashed = h.Sum(nil)
	return hashed, nil
}

func (m *RSAMethod) Sign(data []byte) ([]byte, error) {
	hashed, err := m.hash(data)
	if err != nil {
		return nil, err
	}

	hashed, err = m.hash([]byte(hex.EncodeToString(hashed)))
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(rand.Reader, m.privateKey, m.h, hashed)
}

func (m *RSAMethod) Verify(data []byte, signature []byte) error {
	hashed, err := m.hash(data)
	if err != nil {
		return err
	}

	hashed, err = m.hash([]byte(fmt.Sprintf("%x", hashed)))
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(m.publicKey, m.h, hashed, signature)
}
