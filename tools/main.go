package main

import (
	"ai-software-copyright-server/internal/utils"
	"fmt"
)

func main() {
	// 生成RSA秘钥
	//GenerateRSAKey()

	// 生成AES秘钥
	fmt.Println(utils.GenerateAES128Key())
}
