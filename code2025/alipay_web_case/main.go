package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"
)

// 获取签名
func getSign(params map[string]string, privateKey string) (string, error) {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	signStr := ""
	for _, key := range keys {
		if signStr != "" {
			signStr += "&"
		}
		signStr += fmt.Sprintf("%s=%s", key, params[key])
	}

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return "", fmt.Errorf("failed to decode private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(signStr))
	hashed := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv.(*rsa.PrivateKey), crypto.SHA256, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// 发送支付请求
func sendPaymentRequest(appId, notifyUrl, returnUrl, subject, outTradeNo, totalAmount, privateKey string) error {
	params := map[string]string{
		"app_id":      appId,
		"method":      "alipay.trade.page.pay",
		"format":      "JSON",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  notifyUrl,
		"return_url":  returnUrl,
		"biz_content": fmt.Sprintf(`{"subject":"%s","out_trade_no":"%s","total_amount":"%s","product_code":"FAST_INSTANT_TRADE_PAY"}`, subject, outTradeNo, totalAmount),
	}

	sign, err := getSign(params, privateKey)
	if err != nil {
		return err
	}
	params["sign"] = sign

	// 将参数转换为url.Values形式以便发送表单请求
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	resp, err := http.PostForm("https://openapi.alipay.com/gateway.do", values)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	file, err := os.OpenFile("output.html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 将字符串写入文件
	_, err = file.Write(body)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	appId := ""
	notifyUrl := ""
	returnUrl := ""
	subject := "Test Product"
	outTradeNo := fmt.Sprintf("trade_%d", time.Now().Unix())
	totalAmount := "0.01" // 注意：这里仅为测试，实际金额请按需设置
	privateKey := `-----BEGIN RSA PRIVATE KEY-----
.......
-----END RSA PRIVATE KEY-----`

	if err := sendPaymentRequest(appId, notifyUrl, returnUrl, subject, outTradeNo, totalAmount, privateKey); err != nil {
		fmt.Println("Error:", err)
	}
}
