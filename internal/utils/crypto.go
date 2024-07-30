package utils

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
)

func Md5ByBytes(str string) []byte {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		panic(err)
	}
	return m.Sum(nil)
}

func Md5ByString(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		panic(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}

// 读取私钥文件
func RsaLoadPrivateKeyFile(filePath string) (*rsa.PrivateKey, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return RsaLoadPrivateKey(content)
}

// 读取私钥
func RsaLoadPrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	//获取私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("rsa private key error!")
	}
	//解析PKCS1格式的私钥
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

// 读取公钥文件
func RsaLoadPublicKeyFile(filePath string) (*rsa.PublicKey, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return RsaLoadPublicKey(content)
}

// 读取公钥
func RsaLoadPublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("rsa public key error!")
	}
	// 解析公钥
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

// 公钥加密
func RsaEncrypt(data []byte, key *rsa.PublicKey) ([]byte, error) {
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, key, data)
}

// 私钥解密
func RsaDecrypt(ciphertext []byte, key *rsa.PrivateKey) ([]byte, error) {
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, key, ciphertext)
}

// 签名
func RsaSignWithSha256(data []byte, key *rsa.PrivateKey) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed)
}

// 验证
func RsaVerySignWithSha256(data, signData []byte, key *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signData); err != nil {
		return err
	}
	return nil
}

// Aes加密
func AesEncrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, len(data))
	stream := cipher.NewCTR(block, make([]byte, block.BlockSize()))
	stream.XORKeyStream(ciphertext, data)
	return ciphertext, nil
}

// aes解密
func AesDecrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	stream := cipher.NewCTR(block, make([]byte, block.BlockSize()))
	stream.XORKeyStream(decrypted, data)
	return decrypted, nil
}
