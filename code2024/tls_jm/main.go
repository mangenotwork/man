package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

func decodeWithPrivateKey(data []byte, privateKeyPEM []byte) error {
	// 解析 PEM 格式的私钥
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}

	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		privateKeyEC, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %v", err)
		}

		log.Println(privateKeyEC)

	} else {
		// 进行解码
		plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
		if err != nil {
			return fmt.Errorf("decryption error: %v", err)
		}
		log.Println("解码数据: ", string(plaintext))

	}

	//// 根据私钥类型进行解码
	//switch key := privateKey.(type) {
	//case *rsa.PrivateKey:
	//	// 使用 RSA 私钥进行解码操作
	//	// 此处添加具体的 RSA 解码逻辑
	//	fmt.Println("RSA Private Key")
	//case *ecdsa.PrivateKey:
	//	// 使用 ECDSA 私钥进行解码操作
	//	// 此处添加具体的 ECDSA 解码逻辑
	//	fmt.Println("ECDSA Private Key")
	//default:
	//	return fmt.Errorf("unsupported private key type")
	//}

	return nil
}

func main() {
	// 读取私钥文件
	privateKeyPEM, err := ioutil.ReadFile("privkey.pem")
	if err != nil {
		fmt.Println("Error reading private key file:", err)
		return
	}

	// 待解码的数据
	//data, _ := ioutil.ReadFile("resp.pem")

	str := `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`

	data, _ := base64.StdEncoding.DecodeString(str)

	//log.Println(data)

	err = decodeWithPrivateKey(data, privateKeyPEM)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
}

//// 生成 ECDSA 签名
//func generateECDSASignature(privateKey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
//	r, s, err := ecdsa.Sign(rand.Reader, privateKey, message)
//	if err != nil {
//		return nil, err
//	}
//
//	signature := append(r.Bytes(), s.Bytes()...)
//	return signature, nil
//}
//
//// 验证 ECDSA 签名
//func verifyECDSASignature(publicKey *ecdsa.PublicKey, message []byte, signature []byte) bool {
//	r := new(ecdsa.BigInt)
//	s := new(ecdsa.BigInt)
//
//	r.SetBytes(signature[:len(signature)/2])
//	s.SetBytes(signature[len(signature)/2:])
//
//	return ecdsa.Verify(publicKey, message, r, s)
//}