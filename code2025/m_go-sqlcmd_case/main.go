package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"math/big"
	mathRand "math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	//fmt.Println(UrlExists("https://www.baidu.com"))
	////fmt.Println(UrlExists("http://128.3.112.12:80"))
	//
	//fmt.Println(IsLocalPortAvailable(80))

	case3()

}

// case1 通过请求http Head 请求验证url是否有效
func UrlExists(url string) (exists bool) {
	fmt.Printf("http.Head to %q", url)
	resp, err := http.Head(url)
	if err != nil {
		fmt.Printf("http.Head to %q failed with %v", url, err)
		return false
	}
	if resp.StatusCode != 200 {
		fmt.Printf("http.Head to %q returned status code %d", url, resp.StatusCode)
		return false
	}

	fmt.Printf("http.Head to %q succeeded", url)

	return true
}

// case2 检查本机端口是否被占用
func IsLocalPortAvailable(port int) (portAvailable bool) {
	timeout := time.Second

	hostPort := net.JoinHostPort("localhost", strconv.Itoa(port))
	fmt.Printf(
		"Checking if local port %#v is available using DialTimeout(tcp, %v, timeout: %d)",
		port,
		hostPort,
		timeout,
	)
	conn, err := net.DialTimeout(
		"tcp",
		hostPort,
		timeout,
	)
	if err != nil {
		fmt.Printf(
			"Expected connecting error '%v' on local port %#v, therefore port is available)",
			err,
			port,
		)
		portAvailable = true
	}
	if conn != nil {
		err := conn.Close()
		if err != nil {
			fmt.Printf(
				"Expected closing connection error '%v' on local port %#v, therefore port is available)",
				err,
				port,
			)
			portAvailable = true
		}
		fmt.Printf("Local port '%#v' is not available", port)
	} else {
		fmt.Printf("Local port '%#v' is available", port)
	}

	return
}

// case3 生成指定长度的随机密码。密码
//将至少包含指定数量的特殊字符，
//例如数字和大写字母。中的其余字符
//密码将从小写字母组合中选择，特殊
//例如字母、字符和数字。特殊字符选自
//提供的特殊字符集

func case3() {
	for i := 0; i < 100; i++ {
		fmt.Println(Generate(10, 4, 4, 2, "!@#$%^&*()"))
	}
}

const (
	lowerCharSet = "abcdedfghijklmnopqrstuvwxyz"
	upperCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberSet    = "0123456789"
)

func Generate(passwordLength, minSpecialChar, minNum, minUpperCase int, specialCharSet string) string {
	var password strings.Builder
	allCharSet := lowerCharSet + upperCharSet + specialCharSet + numberSet

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(specialCharSet))))
		checkErr(err)
		_, err = password.WriteString(string(specialCharSet[idx.Int64()]))
		checkErr(err)
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(numberSet))))
		checkErr(err)
		_, err = password.WriteString(string(numberSet[idx.Int64()]))
		checkErr(err)
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(upperCharSet))))
		checkErr(err)
		_, err = password.WriteString(string(upperCharSet[idx.Int64()]))
		checkErr(err)
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(allCharSet))))
		checkErr(err)
		_, err = password.WriteString(string(allCharSet[idx.Int64()]))
		checkErr(err)
	}

	inRune := []rune(password.String())
	mathRand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// case4 使用 simplifiedchinese 中文转码解码
func GbkToUtf8(s []byte) ([]byte, error) {
	//第二个参数为“transform.Transformer”接口，simplifiedchinese.GBK.NewDecoder()包含了该接口
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// case5  编码转换

type Charset string

// 中文
const (
	GBK     Charset = "GBK"
	GB18030         = "GB18030"
	GB2312          = "GB2312"
	Big5            = "Big5"
)

// 日文
const (
	EUCJP     Charset = "EUCJP"
	ISO2022JP         = "ISO2022JP"
	ShiftJIS          = "ShiftJIS"
)

// 韩文
const (
	EUCKR Charset = "EUCKR"
)

// Unicode
const (
	UTF_8    Charset = "UTF-8"
	UTF_16           = "UTF-16"
	UTF_16BE         = "UTF-16BE"
	UTF_16LE         = "UTF-16LE"
)

// 其他编码
const (
	Macintosh Charset = "macintosh"
	IBM               = "IBM*"
	Windows           = "Windows*"
	ISO               = "ISO-*"
)

var charsetAlias = map[string]string{
	"HZGB2312": "HZ-GB-2312",
	"hzgb2312": "HZ-GB-2312",
	"GB2312":   "HZ-GB-2312",
	"gb2312":   "HZ-GB-2312",
}

func Convert(dstCharset Charset, srcCharset Charset, src string) (dst string, err error) {
	if dstCharset == srcCharset {
		return src, nil
	}
	dst = src
	// Converting <src> to UTF-8.
	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewDecoder()),
			)
			if err != nil {
				return "", fmt.Errorf("%s to utf8 failed. %v", srcCharset, err)
			}
			src = string(tmp)
		} else {
			return dst, errors.New(fmt.Sprintf("unsupport srcCharset: %s", srcCharset))
		}
	}
	// Do the converting from UTF-8 to <dstCharset>.
	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := ioutil.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewEncoder()),
			)
			if err != nil {
				return "", fmt.Errorf("utf to %s failed. %v", dstCharset, err)
			}
			dst = string(tmp)
		} else {
			return dst, errors.New(fmt.Sprintf("unsupport dstCharset: %s", dstCharset))
		}
	} else {
		dst = src
	}
	return dst, nil
}

func ToUTF8(srcCharset Charset, src string) (dst string, err error) {
	return Convert("UTF-8", srcCharset, src)
}

func UTF8To(dstCharset Charset, src string) (dst string, err error) {
	return Convert(dstCharset, "UTF-8", src)
}

func getEncoding(charset Charset) encoding.Encoding {
	if c, ok := charsetAlias[string(charset)]; ok {
		charset = Charset(c)
	}
	if e, err := ianaindex.MIB.Encoding(string(charset)); err == nil && e != nil {
		return e
	}
	return nil
}
