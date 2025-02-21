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
	appId := "2021004159667146"
	notifyUrl := "https://www.ecosmos.cc/cn/"
	returnUrl := "https://www.ecosmos.cc/cn/"
	subject := "Test Product"
	outTradeNo := fmt.Sprintf("trade_%d", time.Now().Unix())
	totalAmount := "0.01" // 注意：这里仅为测试，实际金额请按需设置
	privateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCts44WBFiJEefu
6p3bZsQtLaJDMpYOoeDGLrjJ0X3pOIqPlsG3j1wwPcIWdTO4UucZHuNmlukorCtI
hX0Jlei7aW/i8Ol74JcT8QA3twBhsEXDJXjAxAlUVRcUf6dP6FveHVeYhi6l68MC
hq9OrbxntVN5UhqVDwVB1CvmZlyAHNEnb9njz0wq5jkP6MBuZ6pW53J8ejog5e3u
I8dDsrncrZgpzY8ztNgQ3LHKpjTAnYpWRojS79xVs8EwlgPwtuKjaa6Wh2KjKYfA
mtSEW2lbOdhg12smU/d68ml4WAV/BHsuuo/il0DZnqoQrv5J5dNgXtL1Nq87Cxgg
s4N5UrsrAgMBAAECggEBAIkcCVThu0z/AFe7hD1SIhoTQljOjlogdz+YU66imUPF
qMHs2x5coAVISnLVsqyVa+uNUSyChKrhNA07qVYuqZV9hZ7aUULCJh7MhkJ0Rm3V
6Us/wdBPLZoOzHgWx2ew3ws1mBZCHIJF1hmhXLG7O9OU8r36DBeK0riClOB5/hv0
1kViAKHhNsa0t01oDr9dDps44nBFXraaNObWYt3fqOe/OxuKYqsBl6hdap08mNAk
j/AkxgOvz7ANr+p/aXVrrnzAWLJa//MuIoc1kluZjQipF+s3HyAh/JbX56JvsSjx
QvaseJglOVumNOjFsueCVbs4kfQqpBRN9OtFFKcKYSkCgYEA1i1xEhltxeMcNADG
w5RcHYNxIsX8mms1RBPL+EnRr4YS0geipG1xudZAnsTQgm+jaYejZAWaiIy1W5J6
tPQ67+WhZ6H8oIjEEhAA8TW9unydu6//i7v+0JI2SJGAE+txPsF76X7l2foq4Tcf
upIaBkX7zmbKROmrK/TL75WU2MUCgYEAz56+x3RiNEa7jgNc74INvNEg/3RWdw51
uOWLnh8Yb90D0ftSss7U/eyX6Z5yKd9odBwV41BgCMxqlwtIYuF6gjOKt+Z+hgkB
mhDcvya11FmMldAl/GN7rIJnWG3ckppv1q/tDgqWhy2u4P/eiLYcw+pjWTbz45mt
/J/uaEGsIy8CgYA97r/+ktnaWjUCmKLhVVpZsnOZsZS89nldqTfXIUmALw3sLAcM
8xTqvxjKkHEW9r9TOcS2nKQ2DjI3O6E+CE2up0FIHWBW75V6/6O2HGszrOtTpa4I
syEZIN6Pl3toxzFlC0AQogBHSv7xRyZmpe7el4gcBD9DNCqqOExsiF2VXQKBgCcl
kmk/K4kZ0SFcxvgt+HMip2sjP25hXpcHSQT+bfghnyfHkHdAgm6CXr5g7ruwcRx7
czESJZljGbHzIanrQ9Mq7rvwDOku54tqJIUyQlSQse5JefAVverwB5Zn2JAX6IB9
WWAtZOaGGZQ5CneShuf12NeogeHnRyP779LelxtnAoGAXvd1TVCHb3X9kECXF3De
qKZSpMaC1cuCrCu4Gn0vvAoLwo5YfYSuqfoYsq7V4V17w9bc24v21kjN+dPucxzG
9reqYkQyRxoUOBCxVJSiLpDU2VfPf57+iIJKT89NjS0EJeOB7GcJrTe6KQ0EjhFd
7SWStwEaZz2bGCPGXAEwLbE=
-----END RSA PRIVATE KEY-----`

	if err := sendPaymentRequest(appId, notifyUrl, returnUrl, subject, outTradeNo, totalAmount, privateKey); err != nil {
		fmt.Println("Error:", err)
	}
}
