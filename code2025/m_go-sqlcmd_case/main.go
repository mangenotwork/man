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
	"math"
	"math/big"
	mathRand "math/rand"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {

	//fmt.Println(UrlExists("https://www.baidu.com"))
	////fmt.Println(UrlExists("http://128.3.112.12:80"))
	//
	//fmt.Println(IsLocalPortAvailable(80))

	//case3()

	//case4()

	//case8()

	//case9()

	case10()

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

// 分割字符指定前缀后缀
func extractSection(s string, start, end rune) (prefix, body, suffix string, found bool) {
	if strings.HasPrefix(s, string(start)) {
		// no prefix
		body = s[1:]
	} else {
		a := strings.SplitN(s, string(start), 2)
		if len(a) != 2 {
			return "", "", s, false
		}
		prefix = a[0]
		body = a[1]
	}
	a := strings.SplitN(body, string(end), 2)
	if len(a) != 2 {
		return "", "", "", false
	}
	return prefix, a[0], a[1], true
}

func case4() {
	fmt.Println(extractSection("aa(123)cc", '(', ')'))
}

// CopyMap returns a copy of map[string]string
func CopyMap[T comparable, R interface{}](items map[T]R) map[T]R {
	result := make(map[T]R)
	for idx, item := range items {
		result[idx] = item
	}
	return result
}

// CopyStringMap returns a copy of map[string]string
func CopyStringMap(items map[string]string) map[string]string {
	result := make(map[string]string)
	for idx, item := range items {
		result[idx] = item
	}
	return result
}

func case5() {
	source := map[string]int{"foo": 1, "bar": 2}
	duplicate := CopyMap(source)
	fmt.Println(duplicate)

	source2 := map[string]string{"foo": "1", "bar": "2"}
	duplicate2 := CopyStringMap(source2)
	fmt.Println(duplicate2)
}

// 给切片的每个元素应用一个函数
func CollectionApply[T any, R interface{}](collection []T, mutator func(t T) R) []R {
	cast := make([]R, len(collection))
	for i, v := range collection {
		cast[i] = mutator(v)
	}
	return cast
}

func case6() {
	fmt.Println(CollectionApply([]int{1, 2, 3}, func(i int) int {
		return i * 2
	}))
}

var LocalhostStrings = [4]string{"localhost", "[::1]", "::1", "127.0.0.1"}

func isLocalhost(host string) bool {
	normalizedHost := strings.ToLower(host)
	for _, localhostString := range LocalhostStrings {
		if strings.HasPrefix(normalizedHost, localhostString) {
			return isValidRemainder(strings.TrimPrefix(normalizedHost, localhostString))
		}
	}

	return false
}

func isValidRemainder(remainder string) bool {
	return remainder == "" || strings.HasPrefix(remainder, ":")
}

// 判断是否是本地地址
func case7() {
	fmt.Println(isLocalhost("localhost"))
	fmt.Println(isLocalhost("127.0.0.1"))
	fmt.Println(isLocalhost("127.0.0.1:8080"))
	fmt.Println(isLocalhost("127.0.0.1:8080/"))
}

func loadFields(value string) (map[string]string, error) {
	result := make(map[string]string)
	if len(value) == 0 {
		return result, nil
	}
	urlParse := strings.Split(value, "?")
	if len(urlParse) < 2 {
		return result, nil
	}
	parts := strings.Split(urlParse[1], "&")
	for _, part := range parts {
		keyValue := strings.Split(part, "=")
		if len(keyValue) == 2 {
			key, err := sanitizeKey(keyValue[0])
			if err != nil {
				return nil, err
			}
			if result[key] == "" {
				result[key] = keyValue[1]
			} else {
				result[key] += "," + keyValue[1]
			}
		}
	}
	return result, nil
}

func sanitizeKey(key string) (string, error) {
	if key == "" {
		return "", nil
	}
	res, err := url.QueryUnescape(key)
	if err != nil {
		return "", err
	}
	return strings.Trim(res, " "), nil
}

// 提取url中的参数
func case8() {
	rawURL := "https://www.baidu.com/s?wd=aaa&ie=utf-8,1&tn=15007414_9_dg"
	values, err := loadFields(rawURL)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range values {
		fmt.Println(k, v)
	}

	// 标准库实现
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println(err)
	}
	for k, v := range u.Query() {
		fmt.Println(k, v)
	}

}

func isNil(a interface{}) bool {
	if a == nil {
		return true
	}
	val := reflect.ValueOf(a)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			return true
		}
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

// 空判断
func case9() {
	fmt.Println(isNil(nil))
	fmt.Println(isNil([]int{}))
	fmt.Println(isNil(map[string]int{}))
	fmt.Println(isNil(1))
	fmt.Println(isNil("1"))
	fmt.Println(isNil([]int{1}))
	fmt.Println(isNil(map[string]int{"1": 1}))
	fmt.Println(isNil(1))
	fmt.Println(isNil("1"))
	fmt.Println(isNil([]int{1}))
}

func as[T any](in interface{}, out T) error {
	// No point in trying anything if already nil
	if isNil(in) {
		return nil
	}

	// Make sure nothing is a pointer
	valValue := reflect.ValueOf(in)
	for valValue.Kind() == reflect.Ptr {
		valValue = valValue.Elem()
		in = valValue.Interface()
	}

	outVal := reflect.ValueOf(out)
	if outVal.Kind() != reflect.Pointer || isNil(out) {
		return fmt.Errorf("out is not pointer or is nil")
	}

	nestedOutVal := outVal.Elem()
	// Handle the case where out is a pointer to an interface
	if nestedOutVal.Kind() == reflect.Interface && !nestedOutVal.IsNil() {
		nestedOutVal = nestedOutVal.Elem()
	}

	outType := nestedOutVal.Type()

	if !isCompatible(in, outType) {
		return fmt.Errorf("value '%v' is not compatible with type %T", in, nestedOutVal.Interface())
	}

	outVal.Elem().Set(valValue.Convert(outType))
	return nil
}

func isNumericType(in interface{}) bool {

	if in == nil {
		return false
	}

	tp, ok := in.(reflect.Type)
	if !ok {
		tp = reflect.TypeOf(in)
	}

	switch tp.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isCompatible(value interface{}, tp reflect.Type) bool {
	// Can't join with lower, number types are always "convertible" just not losslessly.
	if isNumericType(value) && isNumericType(tp) {
		//NOTE: no need to check if number is compatible with another, always yes, just overflows
		//Check if number value is TRULY compatible
		return isCompatibleInt(value, tp)
	}

	return reflect.TypeOf(value).ConvertibleTo(tp)
}

func isCompatibleInt(in interface{}, tp reflect.Type) bool {
	if !isNumericType(in) || !isNumericType(tp) {
		return false
	}

	inFloat := reflect.ValueOf(in).Convert(reflect.TypeOf(float64(0))).Float()
	hasDecimal := hasDecimalPlace(inFloat)

	if rangeInfo, ok := numericTypeRanges[tp.Kind()]; ok {
		if inFloat >= rangeInfo.min && inFloat <= rangeInfo.max {
			return rangeInfo.allowDecimal || !hasDecimal
		}
	}
	return false
}

type numericRange struct {
	min          float64
	max          float64
	allowDecimal bool
}

var (
	numericTypeRanges = map[reflect.Kind]numericRange{
		reflect.Int8:    {math.MinInt8, math.MaxInt8, false},
		reflect.Uint8:   {0, math.MaxUint8, false},
		reflect.Int16:   {math.MinInt16, math.MaxInt16, false},
		reflect.Uint16:  {0, math.MaxUint16, false},
		reflect.Int32:   {math.MinInt32, math.MaxInt32, false},
		reflect.Uint32:  {0, math.MaxUint32, false},
		reflect.Int64:   {math.MinInt64, math.MaxInt64, false},
		reflect.Uint64:  {0, math.MaxUint64, false},
		reflect.Float32: {-math.MaxFloat32, math.MaxFloat32, true},
		reflect.Float64: {-math.MaxFloat64, math.MaxFloat64, true},
	}
)

func hasDecimalPlace(value float64) bool {
	return value != float64(int64(value))
}

// 类型转换
func case10() {
	var a int = 10
	var b float64
	var c string
	err := as(&a, &b)
	if err != nil {
		fmt.Println(err)
	}
	err = as(&a, &c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%f", b)
	fmt.Printf("%T", b)
	fmt.Printf("%s", c)
	fmt.Printf("%T", c)
}
