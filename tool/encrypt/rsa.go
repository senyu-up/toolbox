package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/senyu-up/toolbox/tool/logger"
	"os"
	"strings"
)

var prefix = "XHHY-SDK"

func PublicKeyFrom(key []byte) (*rsa.PublicKey, error) {
	pubInterface, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key")
	}
	return pub, nil
}

func PublicKeyFrom64(key string) (*rsa.PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	return PublicKeyFrom(b)
}

// 生成 私钥和公钥文件
// java 需要PKCS8类型的私钥,需要将 PKCS1 ->PKCS8
// openssl pkcs8 -topk8 -inform PEM -in private.pem -outform pem -nocrypt -out pkcs8.pem
func GenRsaKey() error {
	//生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	//生成公钥文件
	publicKey := &privateKey.PublicKey
	defPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: defPkix,
	}
	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

// 公钥加密
func RsaEncrypt(data []byte, pubKey []byte) (string, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(pubKey)
	if block == nil {
		return "", errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	encrypt, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		return "", err
	}
	sign := base64.StdEncoding.EncodeToString(encrypt)
	return sign, nil
}

// 私钥解密
func RsaDecrypt(ciphertext string, priKey []byte) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode(priKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	b, err := base64.StdEncoding.DecodeString(ciphertext)
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, b)
}

// aes加密转base64
func AESEncryptBase64(src []byte, key []byte) (encryptedStr string) {
	if src == nil {
		return
	}
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted := make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}
	encryptedStr = base64.StdEncoding.EncodeToString(encrypted)
	return encryptedStr
}

// base64 aes解密
func AESDecryptBase64(encryptedStr string, key []byte) (decrypted []byte, err error) {
	if encryptedStr == "" {
		return
	}
	encrypted, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return
	}
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim], nil
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func Base64UrlDecode(input string) (string, error) {
	input = strings.Replace(input, "-", "+", -1)
	input = strings.Replace(input, "_", "/", -1)
	encrypted, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(encrypted), nil
}

func HashHmac(input, key string, raw bool) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write([]byte(input))
	if raw {
		return string(m.Sum(nil))
	}
	signature := hex.EncodeToString(m.Sum(nil))
	return signature
}

func GetMD5Hash(seed, text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum([]byte(seed)))
}

func GetMD5HashWithoutSeed(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encryption(msg string, timestamp int64) string {
	//hash
	sha := sha1.New()
	sha.Write([]byte(fmt.Sprintf("%s:%d:%s:end", prefix, timestamp, msg)))
	hash := sha.Sum(nil)
	//加密
	sign, err := RsaEncrypt(hash, []byte(""))
	if err != nil {
		logger.Error("Encrypt fail.")
		return ""
	}
	return sign
}

// AesCBCEncrypt aes cbc模式加密后转base64
func AesCBCEncrypt(data, key []byte) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	blockSize := block.BlockSize()
	padding := blockSize - len(data)%blockSize
	text := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(data, text...)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypt := make([]byte, len(data))
	blockMode.CryptBlocks(crypt, data)
	return base64.StdEncoding.EncodeToString(crypt)
}

// AesCBCDecrypt base64解码后 aes cbc模式解密
func AesCBCDecrypt(data string, key []byte) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(dataByte)%blockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(dataByte))
	blockMode.CryptBlocks(crypted, dataByte)
	unPadding := int(crypted[len(crypted)-1])
	return crypted[:(len(crypted) - unPadding)], err
}
