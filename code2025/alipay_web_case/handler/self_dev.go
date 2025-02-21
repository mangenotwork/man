package handler

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"net/url"
	"shop/common/logger"
	"sort"
	"strings"
	"time"
)

func SelfDevAlipayWeb(ctx *gin.Context) {
	appId := ""

	// 2. 用户确认支付后，支付宝通过 get 请求 returnUrl（商户入参传入），返回同步返回参数。
	// 3. 交易成功后，支付宝通过 post 请求 notifyUrl（商户入参传入），返回异步通知参数。

	notifyUrl := ""
	returnUrl := ""
	subject := "测试支付宝网页支付"
	outTradeNo := fmt.Sprintf("trade_%d", time.Now().Unix())
	totalAmount := "0.01" // 注意：这里仅为测试，实际金额请按需设置

	// cc 环境  应用私钥
	privateKey := `-----BEGIN RSA PRIVATE KEY-----
......
-----END RSA PRIVATE KEY-----`

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
		ctx.String(http.StatusOK, err.Error())
		return
	}
	params["sign"] = sign

	// 将参数转换为url.Values形式以便发送表单请求
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	GatewayURL := "https://openapi.alipay.com/gateway.do"

	payURL := fmt.Sprintf("%s?%s", GatewayURL, values.Encode())

	ctx.Redirect(http.StatusFound, payURL)

	return

}

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

func SelfDevAlipayWebReturn(ctx *gin.Context) {
	logger.Info("===== 支付宝回调 Return")
	logger.Info(ctx.Request)
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(string(body))

	ctx.String(http.StatusFound, "ok")
	return
}

func SelfDevAlipayWebNotify(ctx *gin.Context) {
	logger.Info("===== 支付宝通知 Notify")
	logger.Info(ctx.Request)
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(string(body))

	ctx.String(http.StatusFound, "ok")
	return
}

/*

2025-02-20T06:23:54.271208718Z 2025-02-20 14:23:54.270 [INFO]  api/handler/self_dev.go:169      | ===== 支付宝回调 Return
2025-02-20T06:23:54.271240677Z 2025-02-20 14:23:54.271 [INFO]  api/handler/self_dev.go:170      | &{POST /selfDev/alipay/web/return HTTP/1.0 1 0 map[Connection:[close] Content-Length:[1141] Content-Type:[application/x-www-form-urlencoded; charset=utf-8] User-Agent:[Mozilla/4.0] X-Forwarded-For:[203.209.246.124] X-Real-Ip:[203.209.246.124]] 0xc000044300 <nil> 1141 [] true 127.0.0.1:8223 map[] map[] <nil> map[] 172.18.0.1:49630 /selfDev/alipay/web/return <nil> <nil> <nil> 0xc0005a2910 <nil> [] map[]}
2025-02-20T06:23:54.271247797Z 2025-02-20 14:23:54.271 [INFO]  api/handler/self_dev.go:175      | gmt_create=2025-02-20+14%3A23%3A38&charset=utf-8&gmt_payment=2025-02-20+14%3A23%3A51&notify_time=2025-02-20+14%3A23%3A51&subject=%E6%B5%8B%E8%AF%95%E6%94%AF%E4%BB%98%E5%AE%9D%E7%BD%91%E9%A1%B5%E6%94%AF%E4%BB%98&sign=DMNGyp4O5B7v7KO7wQplc8qxUJOFEp11moHFSx%2B3XMAEvuVV471hUbqjbdZ5ISofHOqQVgX484iwyvRdse9tGIfrLlJ5nw5G1QrOlAg0sKrZs0jgLKfhhxwT9mojcLtg%2F8KtNe0FHNFxqSs%2FF6J9VFWMKn3oItETrqHOabwo2zUj7aONhAreuoxPHT0r%2Bnqumt05axZ4OhTZ8rutHYlZgKBuO80OI%2FLxxa4Vh8iWtfasL1J%2BxzzK9thhip10xbaDsUv42xM1ncpm2ILy922joSSe88WAtvV9R5zrphjA1t1iZzNmx8XpAMDfF81zpVj9CiaipDxBsCC5N0C8RDHhWg%3D%3D&merchant_app_id=2021004159667146&buyer_open_id=027AA1FeHLFATyZaXyZAbHW08pjXIxpPqC0NYRP52_TAG8b&invoice_amount=0.01&version=1.0&notify_id=2025022001222142351055271499566012&fund_bill_list=%5B%7B%22amount%22%3A%220.01%22%2C%22fundChannel%22%3A%22ALIPAYACCOUNT%22%7D%5D&notify_type=trade_status_sync&out_trade_no=trade_1740032616&total_amount=0.01&trade_status=TRADE_SUCCESS&trade_no=2025022022001455271437731848&auth_app_id=2021004159667146&receipt_amount=0.01&point_amount=0.00&buyer_pay_amount=0.01&app_id=2021004159667146&sign_type=RSA2&seller_id=2088350716490605
2025-02-20T06:23:54.271265387Z 2025-02-20 14:23:54.271 [INFO]  api/router/base.go:152   |  302 |     194.066µs | 203.209.246.124 | POST | /selfDev/alipay/web/return
2025-02-20T06:24:00.762208946Z 2025-02-20 14:24:00.761 [INFO]  api/handler/self_dev.go:182      | ===== 支付宝通知 Notify
2025-02-20T06:24:00.762246906Z 2025-02-20 14:24:00.762 [INFO]  api/handler/self_dev.go:183      | &{GET /selfDev/alipay/web/notify?charset=utf-8&out_trade_no=trade_1740032616&method=alipay.trade.page.pay.return&total_amount=0.01&sign=aZIX5eNuCp7ihGcXvzkvBxwecE7409poh7aSctZcxmnAKqhJxm98LJ4m4XD64jYj4hbZ0eDLauNKs5NKMZkF%2BFlVfIWz0u2XrxybW4xcsjQ%2FP8%2BEw6%2By0S8ojLHi%2BthC2ViiNWTCqgZYjBvUnhFesZYifSkwmFlyGOtjlRb8wcvwdp5LLw9o9yxieYW5wSlnfCmIxBbS%2FOobRoZGzyQIPpwy61yuclwhkw3NSHtRRJ8ZaG%2BumilQBPUrf0bigD7JmvYHjYrQejII1n4cQRYglZYbQ7J7ILAZhI3%2FeI%2BCxQFa%2FulzANhHZuzzU2PytH3TPF9DGi85pV%2FkuENI0LmX2w%3D%3D&trade_no=2025022022001455271437731848&auth_app_id=2021004159667146&version=1.0&app_id=2021004159667146&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+14%3A23%3A57 HTTP/1.0 1 0 map[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*\/*;q=0.8,application/signed-exchange;v=b3;q=0.7] Accept-Encoding:[gzip, deflate, br, zstd] Accept-Language:[zh-CN,zh;q=0.9] Connection:[close] Cookie:[source=1; brandInfo=%7B%22company_id%22%3A100002%2C%22company_name%22%3A%22%E6%88%90%E9%83%BD%E5%B7%A2%E4%BA%91%E4%BA%92%E5%8A%A8%E7%A7%91%E6%8A%80%E6%9C%89%E9%99%90%E5%85%AC%E5%8F%B8%22%2C%22login_logo%22%3A%22https%3A%2F%2Fecosmos-test-1306984848.cos.ap-guangzhou.myqcloud.com%2Fcompany%2F1876868938924736512.png%3Fq-sign-algorithm%3Dsha1%26q-ak%3DAKID27MnCY2eFsWCnKlItThxt1EubMriDewC%26q-sign-time%3D1739842356%253B1739845956%26q-key-time%3D1739842356%253B1739845956%26q-header-list%3Dhost%26q-url-param-list%3D%26q-signature%3Dc6c3da16d623696efbafd0bb95f1273c353b0351%22%2C%22login_background%22%3A%22https%3A%2F%2Fecosmos-test-1306984848.cos.ap-guangzhou.myqcloud.com%2Fcompany%2F1876869109309947904.png%3Fq-sign-algorithm%3Dsha1%26q-ak%3DAKID27MnCY2eFsWCnKlItThxt1EubMriDewC%26q-sign-time%3D1739842642%253B1739846242%26q-key-time%3D1739842642%253B1739846242%26q-header-list%3Dhost%26q-url-param-list%3D%26q-signature%3D920d4edcbf64e7f2147a17221460188b31db367f%22%2C%22slogan_master%22%3A%22%22%2C%22slogan_slave%22%3A%22%22%7D] Referer:[https://unitradeprod.alipay.com/] Sec-Ch-Ua:["Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"] Sec-Ch-Ua-Mobile:[?0] Sec-Ch-Ua-Platform:["Windows"] Sec-Fetch-Dest:[document] Sec-Fetch-Mode:[navigate] Sec-Fetch-Site:[cross-site] Upgrade-Insecure-Requests:[1] User-Agent:[Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36] X-Forwarded-For:[171.213.205.80] X-Real-Ip:[171.213.205.80]] {} <nil> 0 [] true 127.0.0.1:8223 map[] map[] <nil> map[] 172.18.0.1:49646 /selfDev/alipay/web/notify?charset=utf-8&out_trade_no=trade_1740032616&method=alipay.trade.page.pay.return&total_amount=0.01&sign=aZIX5eNuCp7ihGcXvzkvBxwecE7409poh7aSctZcxmnAKqhJxm98LJ4m4XD64jYj4hbZ0eDLauNKs5NKMZkF%2BFlVfIWz0u2XrxybW4xcsjQ%2FP8%2BEw6%2By0S8ojLHi%2BthC2ViiNWTCqgZYjBvUnhFesZYifSkwmFlyGOtjlRb8wcvwdp5LLw9o9yxieYW5wSlnfCmIxBbS%2FOobRoZGzyQIPpwy61yuclwhkw3NSHtRRJ8ZaG%2BumilQBPUrf0bigD7JmvYHjYrQejII1n4cQRYglZYbQ7J7ILAZhI3%2FeI%2BCxQFa%2FulzANhHZuzzU2PytH3TPF9DGi85pV%2FkuENI0LmX2w%3D%3D&trade_no=2025022022001455271437731848&auth_app_id=2021004159667146&version=1.0&app_id=2021004159667146&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+14%3A23%3A57 <nil> <nil> <nil> 0xc0005a2b40 <nil> [] map[]}
2025-02-20T06:24:00.762395513Z 2025-02-20 14:24:00.762 [INFO]  api/handler/self_dev.go:188      |
2025-02-20T06:24:00.762409743Z 2025-02-20 14:24:00.762 [INFO]  api/router/base.go:152   |  302 |     284.934µs |  171.213.205.80 | GET | /selfDev/alipay/web/notify?charset=utf-8&out_trade_no=trade_1740032616&method=alipay.trade.page.pay.return&total_amount=0.01&sign=aZIX5eNuCp7ihGcXvzkvBxwecE7409poh7aSctZcxmnAKqhJxm98LJ4m4XD64jYj4hbZ0eDLauNKs5NKMZkF%2BFlVfIWz0u2XrxybW4xcsjQ%2FP8%2BEw6%2By0S8ojLHi%2BthC2ViiNWTCqgZYjBvUnhFesZYifSkwmFlyGOtjlRb8wcvwdp5LLw9o9yxieYW5wSlnfCmIxBbS%2FOobRoZGzyQIPpwy61yuclwhkw3NSHtRRJ8ZaG%2BumilQBPUrf0bigD7JmvYHjYrQejII1n4cQRYglZYbQ7J7ILAZhI3%2FeI%2BCxQFa%2FulzANhHZuzzU2PytH3TPF9DGi85pV%2FkuENI0LmX2w%3D%3D&trade_no=2025022022001455271437731848&auth_app_id=2021004159667146&version=1.0&app_id=2021004159667146&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+14%3A23%3A57


gmt_create=2025-02-20+14%3A23%3A38
&charset=utf-8
&gmt_payment=2025-02-20+14%3A23%3A51
&notify_time=2025-02-20+14%3A23%3A51
&subject=%E6%B5%8B%E8%AF%95%E6%94%AF%E4%BB%98%E5%AE%9D%E7%BD%91%E9%A1%B5%E6%94%AF%E4%BB%98
&sign=DMNGyp4O5B7v7KO7wQplc8qxUJOFEp11moHFSx%2B3XMAEvuVV471hUbqjbdZ5ISofHOqQVgX484iwyvRdse9tGIfrLlJ5nw5G1QrOlAg0sKrZs0jgLKfhhxwT9mojcLtg%2F8KtNe0FHNFxqSs%2FF6J9VFWMKn3oItETrqHOabwo2zUj7aONhAreuoxPHT0r%2Bnqumt05axZ4OhTZ8rutHYlZgKBuO80OI%2FLxxa4Vh8iWtfasL1J%2BxzzK9thhip10xbaDsUv42xM1ncpm2ILy922joSSe88WAtvV9R5zrphjA1t1iZzNmx8XpAMDfF81zpVj9CiaipDxBsCC5N0C8RDHhWg%3D%3D&merchant_app_id=2021004159667146
&buyer_open_id=027AA1FeHLFATyZaXyZAbHW08pjXIxpPqC0NYRP52_TAG8b
&invoice_amount=0.01
&version=1.0
&notify_id=2025022001222142351055271499566012
&fund_bill_list=%5B%7B%22amount%22%3A%220.01%22%2C%22fundChannel%22%3A%22ALIPAYACCOUNT%22%7D%5D
&notify_type=trade_status_sync
&out_trade_no=trade_1740032616       商家订单号
&total_amount=0.01
&trade_status=TRADE_SUCCESS
&trade_no=2025022022001455271437731848    支付宝订单号
&auth_app_id=2021004159667146
&receipt_amount=0.01
&point_amount=0.00
&buyer_pay_amount=0.01
&app_id=2021004159667146
&sign_type=RSA2
&seller_id=2088350716490605



*/

// SelfDevAlipayWap
// 文档: https://opendocs.alipay.com/open/29ae8cb6_alipay.trade.wap.pay?pathHash=1ef587fd&ref=api&scene=21
func SelfDevAlipayWap(ctx *gin.Context) {
	// cc环境
	appId := ""

	// 2. 用户确认支付后，支付宝通过 get 请求 returnUrl（商户入参传入），返回同步返回参数。
	// 3. 交易成功后，支付宝通过 post 请求 notifyUrl（商户入参传入），返回异步通知参数。

	notifyUrl := ""
	returnUrl := ""
	subject := "测试支付宝H5支付"
	outTradeNo := fmt.Sprintf("trade_%d", time.Now().Unix())
	totalAmount := "0.01" // 注意：这里仅为测试，实际金额请按需设置

	// cc 环境  应用私钥
	PrivateKey := `-----BEGIN RSA PRIVATE KEY-----
.....
-----END RSA PRIVATE KEY-----`

	privateKey, err := loadPrivateKey(PrivateKey)
	if err != nil {
		logger.Error(err)
		ctx.String(http.StatusOK, err.Error())
		return
	}

	params := map[string]string{
		"app_id":      appId,
		"method":      "alipay.trade.wap.pay",
		"format":      "JSON",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  notifyUrl,
		"return_url":  returnUrl,
		"biz_content": fmt.Sprintf(`{"subject":"%s","out_trade_no":"%s","total_amount":"%s","product_code":"QUICK_WAP_WAY"}`, subject, outTradeNo, totalAmount),
	}

	sign, err := generateSign(params, privateKey)
	if err != nil {
		logger.Error(err)
		ctx.String(http.StatusOK, err.Error())
		return
	}
	params["sign"] = sign

	// 将参数转换为url.Values形式以便发送表单请求
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	GatewayURL := "https://openapi.alipay.com/gateway.do"

	payURL := fmt.Sprintf("%s?%s", GatewayURL, values.Encode())

	ctx.Redirect(http.StatusFound, payURL)

	return
}

// 加载私钥
func loadPrivateKey(privateKeyStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return nil, errors.New("failed to decode private key")
	}
	// 这里需要判断是 PKCS8 格式  还是 PKCS1 格式
	//  PKCS8 格式 （通常以 -----BEGIN PRIVATE KEY----- 开头）  用  x509.ParsePKCS8PrivateKey
	// PKCS1 格式 （即以 -----BEGIN RSA PRIVATE KEY----- 开头） 用  x509.ParsePKCS1PrivateKey
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// 生成签名
func generateSign(params map[string]string, privateKey *rsa.PrivateKey) (string, error) {
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr string
	for _, k := range keys {
		signStr += fmt.Sprintf("%s=%s&", k, params[k])
	}
	signStr = strings.TrimSuffix(signStr, "&")

	hashed := crypto.SHA256.New()
	hashed.Write([]byte(signStr))
	signature, err := rsa.SignPKCS1v15(nil, privateKey, crypto.SHA256, hashed.Sum(nil))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func SelfDevAlipayWapReturn(ctx *gin.Context) {
	logger.Info("===== H5支付 支付宝回调 Return")
	logger.Info(ctx.Request)
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(string(body))

	ctx.String(http.StatusFound, "ok")
	return
}

func SelfDevAlipayWapNotify(ctx *gin.Context) {
	logger.Info("=====H5支付 支付宝通知 Notify")
	logger.Info(ctx.Request)
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Error(err)
	}
	logger.Info(string(body))

	params := make(map[string]string)
	ctx.Request.ParseForm()
	for k, v := range ctx.Request.PostForm {
		params[k] = v[0]
	}

	// 公钥
	AlipayPublicKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAiurrndPpCoZJKcoL3psK
OIx9gu2UP77+d6kIp0OCZMGX3YYWKF+hXkgLl9YOoe3rPzksTDr8LXeDayYxim1O
AT8xTLRsDomRnzabcIanpkYCiX40nJ8UJKJbsEpBxVPUkf9bsGzIZYEcdJvajIuF
BmTE+7AJWdWhWJMKtb3PGIqShZSLCUhzxi/Wf+JIVWzIMqWHLy/X7tBi/tBwamow
Mm15lMY3EXxLrHWTxPNdfPgdbMWXIr3+T6gD3c0biEA5C25RYSY1lLKogusTHnKq
f/eaBZOJZeduK/xnCF8PjPLOjk4tg/xw9TYUPLIO+/+A5+kgqrXH6rsBeoM5CUHO
jwIDAQAB
-----END PUBLIC KEY-----`

	valid, err := verifySign(params, AlipayPublicKey)
	if err != nil || !valid {
		ctx.String(http.StatusBadRequest, "fail")
		return
	}

	// 处理业务逻辑，例如更新订单状态
	ctx.String(http.StatusOK, "success")
}

// 验证签名
func verifySign(params map[string]string, publicKeyStr string) (bool, error) {
	sign, ok := params["sign"]
	if !ok {
		return false, errors.New("sign not found")
	}
	delete(params, "sign")
	delete(params, "sign_type")

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr string
	for _, k := range keys {
		signStr += fmt.Sprintf("%s=%s&", k, params[k])
	}
	signStr = strings.TrimSuffix(signStr, "&")

	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return false, errors.New("failed to decode public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	hashed := crypto.SHA256.New()
	hashed.Write([]byte(signStr))
	signature, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA256, hashed.Sum(nil), signature)
	return err == nil, nil
}

/*

2025-02-20T07:07:52.041904112Z 2025-02-20 15:07:52.041 [INFO]  api/handler/self_dev.go:328      | ===== H5支付 支付宝回调 Return
2025-02-20T07:07:52.041931212Z 2025-02-20 15:07:52.041 [INFO]  api/handler/self_dev.go:329      | &{POST /selfDev/alipay/wap/return HTTP/1.0 1 0 map[Connection:[close] Content-Length:[1197] Content-Type:[application/x-www-form-urlencoded; charset=utf-8] User-Agent:[Mozilla/4.0] X-Forwarded-For:[203.209.246.132] X-Real-Ip:[203.209.246.132]] 0xc00023e580 <nil> 1197 [] true 127.0.0.1:8223 map[] map[] <nil> map[] 172.18.0.1:34052 /selfDev/alipay/wap/return <nil> <nil> <nil> 0xc000244a00 <nil> [] map[]}
2025-02-20T07:07:52.041942691Z 2025-02-20 15:07:52.041 [INFO]  api/handler/self_dev.go:334      | gmt_create=2025-02-20+15%3A07%3A50&charset=utf-8&seller_email=xuyongzhong%40nestclouds.cn&subject=%E6%B5%8B%E8%AF%95%E6%94%AF%E4%BB%98%E5%AE%9DH5%E6%94%AF%E4%BB%98&sign=Kfy47MLtVP1GjXry2u1rhEX0aLckdbAfF4F7mK%2F9TbIPmdbZ73P85er9LwcURmPdd5J2ZmOXYGMpBUKnp1jRZjaxjrE9SrrqEfO%2Bt6Fq5QHzwJ9Z7kFH9WN%2B8nWbL8WZMv2YOYTfDZf5H3rXnIFdVBTNJ4lYpaLtPCtd7YGcIQxljIOss0BxSsunzb7V4s%2BWnWvqupdI%2BgrdfbHWDtbVsyVD7G0M3HuwKvCp%2BdfkB%2BbfAMeBRcBQvFUzu35SQZCTDsN2ZTeL2VMA2Ic560fGoCjN2Zx16LywrF15hUOE8gYaqBPD9iIPN7wnNGmLjYnvIlM8nsUconHbhi6rjB%2FXAQ%3D%3D&buyer_open_id=027h3g6_t9q1OARScw_P2COuwlG8Pf5xMkNBLHqCTP_dos2&invoice_amount=0.01&notify_id=2025022001222150750055271499167962&fund_bill_list=%5B%7B%22amount%22%3A%220.01%22%2C%22fundChannel%22%3A%22ALIPAYACCOUNT%22%7D%5D&notify_type=trade_status_sync&trade_status=TRADE_SUCCESS&receipt_amount=0.01&buyer_pay_amount=0.01&app_id=2021004159649204&sign_type=RSA2&seller_id=2088350716490605&gmt_payment=2025-02-20+15%3A07%3A50&notify_time=2025-02-20+15%3A07%3A50&merchant_app_id=2021004159649204&version=1.0&out_trade_no=trade_1740035258&total_amount=0.01&trade_no=2025022022001455271439940795&auth_app_id=2021004159649204&buyer_logon_id=184****3083&point_amount=0.00
2025-02-20T07:07:52.042034650Z 2025-02-20 15:07:52.041 [INFO]  api/router/base.go:156   |  302 |     207.616µs | 203.209.246.132 | POST | /selfDev/alipay/wap/return
2025-02-20T07:07:59.919019874Z 2025-02-20 15:07:59.918 [INFO]  api/handler/self_dev.go:341      | =====H5支付 支付宝通知 Notify
2025-02-20T07:07:59.919616202Z 2025-02-20 15:07:59.919 [INFO]  api/handler/self_dev.go:342      | &{GET /selfDev/alipay/wap/notify?charset=utf-8&out_trade_no=trade_1740035258&method=alipay.trade.wap.pay.return&total_amount=0.01&sign=EquvuAvL4y5nhnt3d2bioSc%2B7io35PKxeuuUTODuC%2Fc1jkX1Syn6OqN9ocSEgOYDNa%2FQaidEngLgJ0%2F6peiOlNAUvwFEShOUW6zEireLkg9NgVz7Q4HoZ3v%2BRtauJBEnimspWGgeOOVVcV1K7HtIwEJIVJKCmw2RxAxF1fP87L6fUHF6gNTNaSvdKa38BFsvGSRGnm7%2F2y7PYJOkvJSiB%2Fw%2FTqDxkdS6t%2FJqv6TecoYyLOE%2FT4POY%2FDzc%2FS0fvmDG%2BPIIq7w7mkGXS%2F6UmdoaeF0ZcEGNxqJTRTXfWAh6w7xV%2BKxmGBDW%2FEjr14oDZRmjjukPBYdDNjGS8gFGHtWCQ%3D%3D&trade_no=2025022022001455271439940795&auth_app_id=2021004159649204&version=1.0&app_id=2021004159649204&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+15%3A07%3A50 HTTP/1.0 1 0 map[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,*\/*;q=0.8] Accept-Encoding:[gzip, deflate, br] Accept-Language:[zh-CN,zh-Hans;q=0.9] Connection:[close] Priority:[u=0, i] Sec-Fetch-Dest:[document] Sec-Fetch-Mode:[navigate] Sec-Fetch-Site:[none] User-Agent:[Mozilla/5.0 (iPhone; CPU iPhone OS 18_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Mobile/15E148 Safari/604.1] X-Forwarded-For:[171.213.205.80] X-Real-Ip:[171.213.205.80]] {} <nil> 0 [] true 127.0.0.1:8223 map[] map[] <nil> map[] 172.18.0.1:34066 /selfDev/alipay/wap/notify?charset=utf-8&out_trade_no=trade_1740035258&method=alipay.trade.wap.pay.return&total_amount=0.01&sign=EquvuAvL4y5nhnt3d2bioSc%2B7io35PKxeuuUTODuC%2Fc1jkX1Syn6OqN9ocSEgOYDNa%2FQaidEngLgJ0%2F6peiOlNAUvwFEShOUW6zEireLkg9NgVz7Q4HoZ3v%2BRtauJBEnimspWGgeOOVVcV1K7HtIwEJIVJKCmw2RxAxF1fP87L6fUHF6gNTNaSvdKa38BFsvGSRGnm7%2F2y7PYJOkvJSiB%2Fw%2FTqDxkdS6t%2FJqv6TecoYyLOE%2FT4POY%2FDzc%2FS0fvmDG%2BPIIq7w7mkGXS%2F6UmdoaeF0ZcEGNxqJTRTXfWAh6w7xV%2BKxmGBDW%2FEjr14oDZRmjjukPBYdDNjGS8gFGHtWCQ%3D%3D&trade_no=2025022022001455271439940795&auth_app_id=2021004159649204&version=1.0&app_id=2021004159649204&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+15%3A07%3A50 <nil> <nil> <nil> 0xc000095e00 <nil> [] map[]}
2025-02-20T07:07:59.919658590Z 2025-02-20 15:07:59.919 [INFO]  api/handler/self_dev.go:347      |
2025-02-20T07:07:59.919662131Z 2025-02-20 15:07:59.919 [INFO]  api/router/base.go:156   |  400 |     710.286µs |  171.213.205.80 | GET | /selfDev/alipay/wap/notify?charset=utf-8&out_trade_no=trade_1740035258&method=alipay.trade.wap.pay.return&total_amount=0.01&sign=EquvuAvL4y5nhnt3d2bioSc%2B7io35PKxeuuUTODuC%2Fc1jkX1Syn6OqN9ocSEgOYDNa%2FQaidEngLgJ0%2F6peiOlNAUvwFEShOUW6zEireLkg9NgVz7Q4HoZ3v%2BRtauJBEnimspWGgeOOVVcV1K7HtIwEJIVJKCmw2RxAxF1fP87L6fUHF6gNTNaSvdKa38BFsvGSRGnm7%2F2y7PYJOkvJSiB%2Fw%2FTqDxkdS6t%2FJqv6TecoYyLOE%2FT4POY%2FDzc%2FS0fvmDG%2BPIIq7w7mkGXS%2F6UmdoaeF0ZcEGNxqJTRTXfWAh6w7xV%2BKxmGBDW%2FEjr14oDZRmjjukPBYdDNjGS8gFGHtWCQ%3D%3D&trade_no=2025022022001455271439940795&auth_app_id=2021004159649204&version=1.0&app_id=2021004159649204&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+15%3A07%3A50
2025-02-20T07:08:00.663862273Z 2025-02-20 15:08:00.663 [INFO]  api/handler/self_dev.go:341      | =====H5支付 支付宝通知 Notify
2025-02-20T07:08:00.663904272Z 2025-02-20 15:08:00.663 [INFO]  api/handler/self_dev.go:342      | &{GET /selfDev/alipay/wap/notify?charset=utf-8&out_trade_no=trade_1740035258&method=alipay.trade.wap.pay.return&total_amount=0.01&sign=SFSpruOlHUxaEx4rsdyBWDYztPKVhiyJYJTW6ZwHmmp%2FWdj%2FprJTxbeOatJmP3ssjqu6MHzwD9M9WLGeousVCz2IZYSiWLuhg9MwMGfBYXCq19s8iVMMOC2G47vHcFpDKoQFShZLgm9mntf8jdvxUI43244VbETR72XsHOHuayETVbNlkDmJcye%2BfjKmgziZdH%2BQQ7xQGcgxFOdPIRprbSaWTru3ONCwNmhFm66wCnernt0yYhgPXEyCGhE12gOvRBAFcN9S%2FenkO5hjTsvu0QDQyWgHcYgwb2QGiyQl9HQkV5sslKkIVQLKehtW7aN9uoMjY3tHwIESVW%2FAX5KAlw%3D%3D&trade_no=2025022022001455271439940795&auth_app_id=2021004159649204&version=1.0&app_id=2021004159649204&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+15%3A08%3A00 HTTP/1.0 1 0 map[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,*\/*;q=0.8] Accept-Encoding:[gzip, deflate, br] Accept-Language:[zh-CN,zh-Hans;q=0.9] Connection:[close] Priority:[u=0, i] Referer:[https://mclient.alipay.com/] Sec-Fetch-Dest:[document] Sec-Fetch-Mode:[navigate] Sec-Fetch-Site:[cross-site] User-Agent:[Mozilla/5.0 (iPhone; CPU iPhone OS 18_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Mobile/15E148 Safari/604.1] X-Forwarded-For:[171.213.205.80] X-Real-Ip:[171.213.205.80]] {} <nil> 0 [] true 127.0.0.1:8223 map[] map[] <nil> map[] 172.18.0.1:34080 /selfDev/alipay/wap/notify?charset=utf-8&out_trade_no=trade_1740035258&method=alipay.trade.wap.pay.return&total_amount=0.01&sign=SFSpruOlHUxaEx4rsdyBWDYztPKVhiyJYJTW6ZwHmmp%2FWdj%2FprJTxbeOatJmP3ssjqu6MHzwD9M9WLGeousVCz2IZYSiWLuhg9MwMGfBYXCq19s8iVMMOC2G47vHcFpDKoQFShZLgm9mntf8jdvxUI43244VbETR72XsHOHuayETVbNlkDmJcye%2BfjKmgziZdH%2BQQ7xQGcgxFOdPIRprbSaWTru3ONCwNmhFm66wCnernt0yYhgPXEyCGhE12gOvRBAFcN9S%2FenkO5hjTsvu0QDQyWgHcYgwb2QGiyQl9HQkV5sslKkIVQLKehtW7aN9uoMjY3tHwIESVW%2FAX5KAlw%3D%3D&trade_no=2025022022001455271439940795&auth_app_id=2021004159649204&version=1.0&app_id=2021004159649204&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+15%3A08%3A00 <nil> <nil> <nil> 0xc0004a20f0 <nil> [] map[]}
2025-02-20T07:08:00.663922932Z 2025-02-20 15:08:00.663 [INFO]  api/handler/self_dev.go:347      |
2025-02-20T07:08:00.663926072Z 2025-02-20 15:08:00.663 [INFO]  api/router/base.go:156   |  400 |     196.626µs |  171.213.205.80 | GET | /selfDev/alipay/wap/notify?charset=utf-8&out_trade_no=trade_1740035258&method=alipay.trade.wap.pay.return&total_amount=0.01&sign=SFSpruOlHUxaEx4rsdyBWDYztPKVhiyJYJTW6ZwHmmp%2FWdj%2FprJTxbeOatJmP3ssjqu6MHzwD9M9WLGeousVCz2IZYSiWLuhg9MwMGfBYXCq19s8iVMMOC2G47vHcFpDKoQFShZLgm9mntf8jdvxUI43244VbETR72XsHOHuayETVbNlkDmJcye%2BfjKmgziZdH%2BQQ7xQGcgxFOdPIRprbSaWTru3ONCwNmhFm66wCnernt0yYhgPXEyCGhE12gOvRBAFcN9S%2FenkO5hjTsvu0QDQyWgHcYgwb2QGiyQl9HQkV5sslKkIVQLKehtW7aN9uoMjY3tHwIESVW%2FAX5KAlw%3D%3D&trade_no=2025022022001455271439940795&auth_app_id=2021004159649204&version=1.0&app_id=2021004159649204&sign_type=RSA2&seller_id=2088350716490605&timestamp=2025-02-20+15%3A08%3A00


*/
