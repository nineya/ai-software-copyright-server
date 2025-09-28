package main

import (
	"ai-software-copyright-server/internal/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GenerateRSAKey() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	utils.PanicErr(err)
	// 将私钥序列化为PEM格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	fmt.Println(string(privateKeyPem))
	utils.WriteToFile("rsa_private.pem", privateKeyPem)
	//生成公钥
	publicKey := privateKey.PublicKey
	// 将私钥序列化为PEM格式
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&publicKey)
	publicKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	fmt.Println(string(publicKeyPem))
	utils.WriteToFile("rsa_public.pem", publicKeyPem)
}
