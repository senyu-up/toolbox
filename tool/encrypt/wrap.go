package encrypt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func wrap(seed, ts int, values map[string]interface{}) *strings.Builder {
	var result strings.Builder
	if seed > len(values) || seed < 0 {
		return &result
	}
	values["ts"] = ts
	mapper := newMapper(values)
	sort.Sort(mapper)
	result.WriteString("ts")
	result.WriteString("$")
	result.WriteString(strconv.Itoa(ts))
	return splice(seed, mapper, &result)
}

type MapItem struct {
	Key   string
	Value interface{}
}

type Mapper []*MapItem

func (p Mapper) Less(i, j int) bool {
	return p[i].Key < p[j].Key
}

func (p Mapper) Len() int {
	return len(p)
}

func (p Mapper) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func newMapper(values map[string]interface{}) Mapper {
	m := make(Mapper, 0)
	for key, value := range values {
		m = append(m, &MapItem{Key: key, Value: value})
	}
	return m
}

func splice(seed int, mapper []*MapItem, str *strings.Builder) *strings.Builder {
	var seedStr strings.Builder
	var kind reflect.Kind
	var value interface{}
	for i, item := range mapper {
		value = item.Value
		if value == nil {
			value = ""
			kind = reflect.String
		} else {
			kind = reflect.TypeOf(value).Kind()
		}
		if i != seed && item.Key != "ts" {
			str.WriteString("&")
			str.WriteString(item.Key)
			str.WriteString("$")
			if isCompositeStructure(kind) {
				v, _ := json.Marshal(value)
				str.Write(v)
			} else {
				str.WriteString(cast.ToString(value))
			}
		}
		if i == seed {
			seedStr.WriteString(item.Key)
			seedStr.WriteString("$")
			if isCompositeStructure(kind) {
				v, _ := json.Marshal(value)
				seedStr.Write(v)
			} else {
				seedStr.WriteString(cast.ToString(value))
			}
		}
	}
	str.WriteString("&")
	str.WriteString(seedStr.String())
	return str
}

func isCompositeStructure(k reflect.Kind) bool {
	switch k {
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		fallthrough
	case reflect.Slice:
		return true
	default:
		return false
	}
}

func NewMapper(values map[string]interface{}) Mapper {
	return newMapper(values)
}

func Signature(target, priKey string) (string, error) {
	hash256 := sha256.New()
	hash256.Write([]byte(target))
	hash := hash256.Sum(nil)
	block, _ := pem.Decode([]byte(priKey))
	if block == nil {
		return "", errors.New("private key error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS1PrivateKey err", err)
		return "", nil
	}
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		fmt.Println("Sign err:", err)
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}

// 验证签名
func VerifySignatureWithPubKey(target, sign, pubKey string) (bool, error) {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return false, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	publicKey := pubInterface.(*rsa.PublicKey)
	hash256 := sha256.New()
	hash256.Write([]byte(target))
	hash := hash256.Sum(nil)
	bt, _ := base64.StdEncoding.DecodeString(sign)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, bt)
	if err != nil {
		return false, err
	}
	return true, nil
}
