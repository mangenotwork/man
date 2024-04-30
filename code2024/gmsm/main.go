/*

	https://emmansun.github.io/gmsm/

	Go语言商用密码软件: Go语言商用密码软件，简称GMSM，一个安全、高性能、易于使用的Go语言商用密码软件库，涵盖商用密码公开算法SM2/SM3/SM4/SM9/ZUC。
*/

package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/emmansun/gmsm/sm2"
	"log"
	"os"
)

func main() {

	toSign := "ShangMi SM2 Sign Standard"
	sig := ExamplePrivateKey_Sign_forceSM2(toSign)

	ExampleVerifyASN1WithSM2(sig)
}

// 您可以直接使用sm2私钥的签名方法Sign
func ExamplePrivateKey_Sign_forceSM2(test string) string {
	toSign := []byte(test)
	// real private key should be from secret storage
	privKey, _ := hex.DecodeString("6c5a0a0b2eed3cbec3e4f1252bfe0e28c504a1c6bf1999eebb0af9ef0f8e6c85")
	testkey, err := sm2.NewPrivateKey(privKey)
	if err != nil {
		log.Fatalf("fail to new private key %v", err)
	}

	// force SM2 sign standard and use default UID
	sig, err := testkey.Sign(rand.Reader, toSign, sm2.DefaultSM2SignerOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from sign: %s\n", err)
		return ""
	}

	// Since sign is a randomized function, signature will be
	// different each time.
	fmt.Printf("%x\n", sig)
	return fmt.Sprintf("%x", sig)
}

// 您可以使用sm2.VerifyASN1WithSM2来校验SM2签名：
func ExampleVerifyASN1WithSM2(sig string) {
	// real public key should be from cert or public key pem file
	keypoints, _ := hex.DecodeString("048356e642a40ebd18d29ba3532fbd9f3bbee8f027c3f6f39a5ba2f870369f9988981f5efe55d1c5cdf6c0ef2b070847a14f7fdf4272a8df09c442f3058af94ba1")
	testkey, err := sm2.NewPublicKey(keypoints)
	if err != nil {
		log.Fatalf("fail to new public key %v", err)
	}

	toSign := []byte("ShangMi SM2 Sign Standard")
	//signature, _ := hex.DecodeString("304402205b3a799bd94c9063120d7286769220af6b0fa127009af3e873c0e8742edc5f890220097968a4c8b040fd548d1456b33f470cabd8456bfea53e8a828f92f6d4bdcd77")
	signature, _ := hex.DecodeString(sig)
	ok := sm2.VerifyASN1WithSM2(testkey, nil, toSign, signature)

	fmt.Printf("%v\n", ok)
	// Output: true
}
