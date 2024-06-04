package encrypt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

type Check struct {
	//公钥
	publicKey *rsa.PublicKey
}

func NewCheck(publicKey string) (*Check, error) {
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &Check{
		publicKey: pubInterface.(*rsa.PublicKey),
	}, nil
}

type Tag struct {
	//时间戳
	Ts int
	//随机种子
	Seed int
	//签名
	Signature string
	//数据
	Data map[string]interface{}
}

// Check 验证签名
func (c *Check) Check(t *Tag) (bool, error) {
	bt, err := base64.StdEncoding.DecodeString(t.Signature)
	if err != nil {
		return false, err
	}
	hash256 := sha256.New()
	w := wrap(t.Seed, t.Ts, t.Data).String()
	hash256.Write([]byte(w))
	h := hash256.Sum(nil)
	if err = rsa.VerifyPKCS1v15(c.publicKey, crypto.SHA256, h, bt); err != nil {
		return false, err
	}
	return true, nil
}
