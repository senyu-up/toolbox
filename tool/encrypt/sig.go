package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

type Sig struct {
	//私钥
	privateKey *rsa.PrivateKey
}

func NewSig(privateKey string) (*Sig, error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &Sig{privateKey: key}, nil
}

func (s *Sig) Signature(t *Tag) (string, error) {
	hash256 := sha256.New()
	_, err := hash256.Write([]byte(wrap(t.Seed, t.Ts, t.Data).String()))
	if err != nil {
		return "", err
	}
	hash := hash256.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA256, hash)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}
