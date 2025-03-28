package main

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	gt "github.com/mangenotwork/gathertool"
)

func main() {
	svc := &Intercept{
		Ip: "0.0.0.0:18111",
		HttpPackageFunc: func(pack *HttpPackage) {
			// 查看 ContentType
			//log.Println("ContentType = ", pack.ContentType)
			//// 获取数据包数据为 http,json等 TXT格式的数据
			//log.Println("Txt = ", pack.Html())
			//// 获取数据包为图片，将图片转为 base64
			//log.Println("img base64 = ", pack.Img2Base64())
			//// 获取数据包为图片，存储图片
			//log.Println(pack.SaveImage(""))
		},
	}
	svc.RunServer()

	//main3()
}

// Intercept http/s 代理与抓包
type Intercept struct {
	Ip              string
	HttpPackageFunc func(pack *HttpPackage)
}

// RunServer 启动 http/s 代理与抓包服务
func (ipt *Intercept) RunServer() {

	cert, err := tls.LoadX509KeyPair("ssl/server.pem", "ssl/privkey.pem")
	if err != nil {
		gt.Panic(err)
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	gt.Info("启动代理&抓包 <ManGe代理&抓包> ......... ")
	gt.Info(" - HTTPS代理 : 只支持代理转发  -> ", ipt.Ip)
	gt.Info(" - HTTP代理: 支持数据包处理与代理转发  -> ", ipt.Ip)

	//cert, err := genCertificate()
	//if err != nil {
	//	panic(err)
	//}

	server := &http.Server{
		Addr:      ipt.Ip,
		TLSConfig: cfg, //&tls.Config{Certificates: []tls.Certificate{cert}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			id := gt.IDMd5()

			var reqBodyBytes []byte
			reqBodyBytes, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
			gt.Info("请求 body:  ", string(reqBodyBytes))
			gt.Info("请求 Method :  ", r.Method)
			gt.Info("请求 : ", r)

			if r.Method == http.MethodConnect {
				gt.Info("HTTPS 请求 | ", r.URL.String())

				destConn, err := net.DialTimeout("tcp", r.Host, 60*time.Second)
				if err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}

				w.WriteHeader(http.StatusOK)
				hijacker, ok := w.(http.Hijacker)
				if !ok {
					http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
					return
				}

				clientConn, _, err := hijacker.Hijack()
				if err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
				}

				//				// 写入数据到客户端
				//				_, err = clientConn.Write([]byte(`HTTP/1.1 200 OK
				//Date: Mon, 23 May 2023 20:00:00 GMT
				//Content-Type: text/html
				//Content-Length: 152
				//Last-Modified: Thu, 12 Mar 2021 18:23:50 GMT
				//Server: Apache 2.4.1 (Unix)
				//
				//<!DOCTYPE html>
				//<html>
				//<head>
				//    <title>Example Domain</title>
				//</head>
				//<body>
				//<div>
				//    <h1>Example Domain</h1>
				//    <p>This domain is established to be used for illustrative examples in documents. You may use this
				//    domain in examples without prior coordination or asking for permission.</p>
				//</div>
				//</body>
				//</html>\n`))
				//				if err != nil {
				//					return
				//				}
				//
				//				reader := bufio.NewReader(clientConn)
				//
				//				go func() {
				//					for {
				//						line, err := reader.ReadString('\n')
				//						if err != nil {
				//							return
				//						}
				//						gt.Info("line --> ", string(line))
				//					}
				//				}()
				//
				//				time.Sleep(2 * time.Second)
				//				clientConn.Close()

				//var rwb []byte
				//_, _ = rw.Read(rwb)
				//gt.Info("rwb --> ", string(rwb))
				//
				//b, _ := io.ReadAll(clientConn)
				//gt.Info("b --> ", string(b))

				//go func() {
				//	// 读取客户端数据
				//	reader := bufio.NewReader(clientConn)
				//	for {
				//		line, _, err := reader.ReadLine()
				//		if err != nil {
				//			gt.Info("读clientConn 关闭 clientConn")
				//			destConn.Close()
				//			return
				//		}
				//
				//		println(string(line))
				//
				//		_, err = destConn.Write(line)
				//		if err != nil {
				//			gt.Info("读clientConn 关闭 destConn")
				//			clientConn.Close()
				//			return
				//		}
				//	}
				//}()

				//go func() {
				//	reader := bufio.NewReader(destConn)
				//	for {
				//		line, _, err := reader.ReadLine()
				//		if err != nil {
				//			gt.Info("读destConn 关闭 destConn")
				//			clientConn.Close()
				//			return
				//		}
				//
				//		println(string(line))
				//
				//		_, err = clientConn.Write(line)
				//		if err != nil {
				//			gt.Info("读destConn 关闭 clientConn")
				//			destConn.Close()
				//			return
				//		}
				//	}
				//}()

				go func() {

					//fmt.Println("1")
					//
					//b, _ := io.ReadAll(destConn)
					//
					//gt.Info(string(b))
					//
					//_, _ = destConn.Write(b)

					//var srcBuffer bytes.Buffer
					//_, err = io.Copy(&srcBuffer, destConn)
					//if err != nil {
					//	gt.Error(err)
					//}
					//
					//gt.Info(srcBuffer.String())
					//
					//var reader io.Reader
					//reader = &srcBuffer

					//var srcBuffer io.Writer = io.Discard

					//srcBuffer := io.MultiWriter(os.Stdout, &buf)

					psth := id + "-resp.txt"

					gt.Info(psth)

					file, err := os.Create(psth)
					if err != nil {
						gt.Error(err)
						return
					}
					defer file.Close()

					// 创建一个Base64编码的Writer
					encoder := base64.NewEncoder(base64.StdEncoding, file)

					readerTee := io.TeeReader(destConn, encoder)
					defer encoder.Close()

					//readerTee := io.TeeReader(destConn, os.Stdout)

					_, _ = io.Copy(clientConn, readerTee)
					//if n > 0 {
					//	b, err := io.ReadAll(clientConn)
					//	if err != nil {
					//		fmt.Println("Error reading data:", err)
					//		return
					//	}
					//	// 输出读取到的数据
					//	fmt.Println(string(b))
					//}
					//gt.Info("111")
					//b := make([]byte, 0)
					//n, _ := clientConn.Read(b)
					//
					//fmt.Println("读取到的数据：", string(b[:n]))

				}()

				go func() {
					//fmt.Println("2")
					//var srcBuffer io.Writer = io.Discard
					//readerTee := io.TeeReader(clientConn, srcBuffer)
					//gt.Info(srcBuffer)

					psth := id + "-req.txt"
					gt.Info(psth)

					file, err := os.Create(psth)
					if err != nil {
						gt.Error(err)
						return
					}
					defer file.Close()
					readerTee := io.TeeReader(clientConn, file)
					//readerTee := io.TeeReader(clientConn, os.Stdout)
					_, _ = io.Copy(destConn, readerTee)

				}()

				//go func() {
				//	file, err := os.Create("filename")
				//	if err != nil {
				//		return
				//	}
				//	defer file.Close()
				//
				//	_, err = io.Copy(file, destConn)
				//	if err != nil {
				//		return
				//	}
				//
				//	return
				//}()

			} else {
				gt.Info("HTTP 请求 | ", r.URL.String())
				res, err := http.DefaultTransport.RoundTrip(r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusServiceUnavailable)
					return
				}
				defer func() {
					_ = res.Body.Close()
				}()
				for k, vv := range res.Header {
					for _, v := range vv {
						w.Header().Add(k, v)
					}
				}

				var bodyBytes []byte
				bodyBytes, _ = io.ReadAll(res.Body)
				res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				w.WriteHeader(res.StatusCode)
				httpPackage := &HttpPackage{
					Url:         r.URL,
					ContentType: res.Header.Get("Content-Type"),
					Body:        bodyBytes,
					Header:      res.Header,
				}
				ipt.HttpPackageFunc(httpPackage)
				_, _ = io.Copy(w, res.Body)
				_ = res.Body.Close()
			}
		}),
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

var errInvalidWrite = errors.New("invalid write result")
var ErrShortWrite = errors.New("short write")
var EOF = errors.New("EOF")

func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {

	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		gt.Info("WriteTo")
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		gt.Info("ReadFrom")
		return rt.ReadFrom(src)
	}
	if buf == nil {
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != EOF {
				err = er
			}
			break
		}
	}
	return written, err
}

// RunHttpIntercept 启动 http/s 代理与抓包服务
func (ipt *Intercept) RunHttpIntercept() {
	gt.Info("启动抓包 <ManGe抓包> ......... ")
	gt.Info("目前只支持HTTP, HTTPS还在开发中 ......... ")
	gt.Info("请在系统设置代理 HTTP代理  ", ipt.Ip)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		gt.Info("\n\n___________________________________________________________________________")
		gt.Info("代理请求信息： ", r.RemoteAddr, r.Method, r.URL.String())
		transport := http.DefaultTransport
		outReq := new(http.Request)
		*outReq = *r
		if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
				clientIP = strings.Join(prior, ", ") + ", " + clientIP
			}
			outReq.Header.Set("X-Forwarded-For", clientIP)
		}
		res, err := transport.RoundTrip(outReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		for key, value := range res.Header {
			for _, v := range value {
				w.Header().Add(key, v)
			}
		}
		var bodyBytes []byte
		bodyBytes, _ = io.ReadAll(res.Body)
		res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		w.WriteHeader(res.StatusCode)
		httpPackage := &HttpPackage{
			Url:         r.URL,
			ContentType: res.Header.Get("Content-Type"),
			Body:        bodyBytes,
			Header:      res.Header,
		}
		ipt.HttpPackageFunc(httpPackage)
		_, _ = io.Copy(w, res.Body)
		_ = res.Body.Close()
	})
	err := http.ListenAndServe(ipt.Ip, nil)
	if err != nil {
		panic(err)
	}
}

func genCertificate() (cert tls.Certificate, err error) {
	rawCert, rawKey, err := generateKeyPair()
	if err != nil {
		return
	}
	return tls.X509KeyPair(rawCert, rawKey)

}

func generateKeyPair() (rawCert, rawKey []byte, err error) {
	// Create private key and self-signed certificate
	// Adapted from https://golang.org/src/crypto/tls/generate_cert.go

	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	validFor := time.Hour * 24 * 365 * 10 // ten years
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"ManGe-gatherTool"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &private.PublicKey, private)
	if err != nil {
		return
	}

	rawCert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	rawKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(private)})

	keyFileName := "ca.key"
	certFIleName := "ca.crt"
	kf, err := os.Create(keyFileName)
	if err != nil {
		return
	}
	defer kf.Close()

	_, _ = kf.Write(rawKey)

	cf, err := os.Create(certFIleName)
	if err != nil {
		return
	}
	defer cf.Close()
	_, _ = cf.Write(rawCert)

	return
}

// HttpPackage 代理服务抓取到的HTTP的包
type HttpPackage struct {
	Url         *url.URL
	Body        []byte
	ContentType string
	Header      map[string][]string
}

// Img2Base64 如果数据类型是image 就转换成base64的图片输出
func (pack *HttpPackage) Img2Base64() string {
	if strings.Index(pack.ContentType, "image") != -1 {
		return base64.StdEncoding.EncodeToString(pack.Body)
	}
	return ""
}

// Html 数据类型是html
func (pack *HttpPackage) Html() string {
	if strings.Index(pack.ContentType, "html") != -1 {
		rdata := strings.NewReader(string(pack.Body))
		r, err := gzip.NewReader(rdata)
		if err == nil {
			s, _ := io.ReadAll(r)
			return string(s)
		}
	}
	return ""
}

// SaveImage 如果数据类型是image 就保存图片
func (pack *HttpPackage) SaveImage(path string) error {
	if strings.Index(pack.ContentType, "image") != -1 {
		idx := strings.LastIndex(pack.Url.String(), "/")
		if idx < 0 {
			path += pack.Url.String()
		} else {
			path += pack.Url.String()[idx+1:]
		}
		return os.WriteFile(path, pack.Body, 0666)
	}
	return fmt.Errorf("ContentType not image")
}

// Json 数据类型是json
func (pack *HttpPackage) Json() string {
	if strings.Index(pack.ContentType, "json") != -1 {
		return string(pack.Body)
	}
	return ""
}

// Txt 数据类型是txt
func (pack *HttpPackage) Txt() string {
	if strings.Index(pack.ContentType, "txt") != -1 {
		rdata := strings.NewReader(string(pack.Body))
		r, err := gzip.NewReader(rdata)
		if err == nil {
			s, _ := io.ReadAll(r)
			return string(s)
		}
	}
	return ""
}

// ToFile 抓取到的数据类型保存到文件
func (pack *HttpPackage) ToFile(path string) error {
	ext := ContentType[pack.ContentType]
	path = path + gt.Any2String(time.Now().UnixNano()) + ext
	return os.WriteFile(path, pack.Body, 0666)
}

// ContentType 数据类型
var ContentType = map[string]string{
	"application/octet-stream":            ".*",
	"application/x-001":                   ".001",
	"application/x-301":                   ".301",
	"text/h323":                           ".323",
	"application/x-906":                   ".906",
	"drawing/907":                         ".907",
	"application/x-a11":                   ".a11",
	"audio/x-mei-aac":                     ".acp",
	"application/postscript":              ".ai",
	"audio/aiff":                          ".aif",
	"application/x-anv":                   ".anv",
	"text/asa":                            ".asa",
	"video/x-ms-asf":                      ".asf",
	"text/asp":                            ".asp",
	"audio/basic":                         ".au",
	"video/avi":                           ".avi",
	"application/vnd.adobe.workflow":      ".awf",
	"text/xml":                            ".biz",
	"application/x-bmp":                   ".bmp",
	"application/x-bot":                   ".bot",
	"application/x-c4t":                   ".c4t",
	"application/x-c90":                   ".c90",
	"application/x-cals":                  ".cal",
	"application/vnd.ms-pki.seccat":       ".cat",
	"application/x-netcdf":                ".cdf",
	"application/x-cdr":                   ".cdr",
	"application/x-cel":                   ".cel",
	"application/x-x509-ca-cert":          ".cer",
	"application/x-g4":                    ".cg4",
	"application/x-cgm":                   ".cgm",
	"application/x-cit":                   ".cit",
	"java/*":                              ".class",
	"application/x-cmp":                   ".cmp",
	"application/x-cmx":                   ".cmx",
	"application/x-cot":                   ".cot",
	"application/pkix-crl":                ".crl",
	"application/x-csi":                   ".csi",
	"text/css":                            ".css",
	"application/x-cut":                   ".cut",
	"application/x-dbf":                   ".dbf",
	"application/x-dbm":                   ".dbm",
	"application/x-dbx":                   ".dbx",
	"application/x-dcx":                   ".dcx",
	"application/x-dgn":                   ".dgn",
	"application/x-dib":                   ".dib",
	"application/x-msdownload":            ".exe",
	"application/msword":                  ".doc",
	"application/x-drw":                   ".drw",
	"Model/vnd.dwf":                       ".dwf",
	"application/x-dwf":                   ".dwf",
	"application/x-dwg":                   ".dwg",
	"application/x-dxb":                   ".dxb",
	"application/x-dxf":                   ".dxf",
	"application/vnd.adobe.edn":           ".edn",
	"application/x-emf":                   ".emf",
	"message/rfc822":                      ".eml",
	"application/x-epi":                   ".epi",
	"application/x-ps":                    ".eps",
	"application/x-ebx":                   ".etd",
	"image/fax":                           ".fax",
	"application/vnd.fdf":                 ".fdf",
	"application/fractals":                ".fif",
	"application/x-frm":                   ".frm",
	"application/x-gbr":                   ".gbr",
	"application/x-gcd":                   ".gcd",
	"image/gif":                           ".gif",
	"application/x-gl2":                   ".gl2",
	"application/x-gp4":                   ".gp4",
	"application/x-hgl":                   ".hgl",
	"application/x-hmr":                   ".hmr",
	"application/x-hpgl":                  ".hpg",
	"application/x-hpl":                   ".hpl",
	"application/mac-binhex40":            ".hqx",
	"application/x-hrf":                   ".hrf",
	"application/hta":                     ".hta",
	"text/x-component":                    ".htc",
	"text/html":                           ".html",
	"text/webviewhtml":                    ".htt",
	"application/x-icb":                   ".icb",
	"image/x-icon":                        ".ico",
	"application/x-ico":                   ".ico",
	"application/x-iff":                   ".iff",
	"application/x-igs":                   ".igs",
	"application/x-iphone":                ".iii",
	"application/x-img":                   ".img",
	"application/x-internet-signup":       ".ins",
	"video/x-ivf":                         ".IVF",
	"image/jpeg":                          ".jpg",
	"application/x-jpe":                   ".jpe",
	"application/x-jpg":                   ".jpg",
	"application/x-javascript":            ".js",
	"audio/x-liquid-file":                 ".la1",
	"application/x-laplayer-reg":          ".lar",
	"application/x-latex":                 ".latex",
	"audio/x-liquid-secure":               ".lavs",
	"application/x-lbm":                   ".lbm",
	"audio/x-la-lms":                      ".lmsff",
	"application/x-ltr":                   ".ltr",
	"video/x-mpeg":                        ".m1v",
	"audio/mpegurl":                       ".m3u",
	"video/mpeg4":                         ".m4e",
	"application/x-mac":                   ".mac",
	"application/x-troff-man":             ".man",
	"application/msaccess":                ".mdb",
	"application/x-mdb":                   ".mdb",
	"application/x-shockwave-flash":       ".mfp",
	"application/x-mi":                    ".mi",
	"audio/mid":                           ".mid",
	"application/x-mil":                   ".mil",
	"audio/x-musicnet-download":           ".mnd",
	"audio/x-musicnet-stream":             ".mns",
	"video/x-sgi-movie":                   ".movie",
	"audio/mp1":                           ".mp1",
	"audio/mp2":                           ".mp2",
	"video/mpeg":                          ".mp2v",
	"audio/mp3":                           ".mp3",
	"video/x-mpg":                         ".mpa",
	"application/vnd.ms-project":          ".mpd",
	"video/mpg":                           ".mpeg",
	"audio/rn-mpeg":                       ".mpga",
	"application/x-mmxp":                  ".mxp",
	"image/pnetvue":                       ".net",
	"application/x-nrf":                   ".nrf",
	"text/x-ms-odc":                       ".odc",
	"application/x-out":                   ".out",
	"application/pkcs10":                  ".p10",
	"application/x-pkcs12":                ".p12",
	"application/x-pkcs7-certificates":    ".p7b",
	"application/pkcs7-mime":              ".p7c",
	"application/x-pkcs7-certreqresp":     ".p7r",
	"application/pkcs7-signature":         ".p7s",
	"application/x-pc5":                   ".pc5",
	"application/x-pci":                   ".pci",
	"application/x-pcl":                   ".pcl",
	"application/x-pcx":                   ".pcx",
	"application/pdf":                     ".pdf",
	"application/vnd.adobe.pdx":           ".pdx",
	"application/x-pgl":                   ".pgl",
	"application/x-pic":                   ".pic",
	"application/vnd.ms-pki.pko":          ".pko",
	"application/x-perl":                  ".pl",
	"audio/scpls":                         ".pls",
	"application/x-plt":                   ".plt",
	"image/png":                           ".png",
	"application/x-png":                   ".png",
	"application/vnd.ms-powerpoint":       ".ppt",
	"application/x-ppm":                   ".ppm",
	"application/x-ppt":                   ".ppt",
	"application/x-pr":                    ".pr",
	"application/pics-rules":              ".prf",
	"application/x-prn":                   ".prn",
	"application/x-prt":                   ".prt",
	"application/x-ptn":                   ".ptn",
	"text/vnd.rn-realtext3d":              ".r3t",
	"audio/vnd.rn-realaudio":              ".ra",
	"audio/x-pn-realaudio":                ".ram",
	"application/x-ras":                   ".ras",
	"application/rat-file":                ".rat",
	"application/vnd.rn-recording":        ".rec",
	"application/x-red":                   ".red",
	"application/x-rgb":                   ".rgb",
	"application/vnd.rn-realsystem-rjs":   ".rjs",
	"application/vnd.rn-realsystem-rjt":   ".rjt",
	"application/x-rlc":                   ".rlc",
	"application/x-rle":                   ".rle",
	"application/vnd.rn-realmedia":        ".rm",
	"application/vnd.adobe.rmf":           ".rmf",
	"application/vnd.rn-realsystem-rmj":   ".rmj",
	"application/vnd.rn-rn_music_package": ".rmp",
	"application/vnd.rn-realmedia-secure": ".rms",
	"application/vnd.rn-realmedia-vbr":    ".rmvb",
	"application/vnd.rn-realsystem-rmx":   ".rmx",
	"application/vnd.rn-realplayer":       ".rnx",
	"image/vnd.rn-realpix":                ".rp",
	"audio/x-pn-realaudio-plugin":         ".rpm",
	"application/vnd.rn-rsml":             ".rsml",
	"text/vnd.rn-realtext":                ".rt",
	"application/x-rtf":                   ".rtf",
	"video/vnd.rn-realvideo":              ".rv",
	"application/x-sam":                   ".sam",
	"application/x-sat":                   ".sat",
	"application/sdp":                     ".sdp",
	"application/x-sdw":                   ".sdw",
	"application/x-stuffit":               ".sit",
	"application/x-slb":                   ".slb",
	"application/x-sld":                   ".sld",
	"drawing/x-slk":                       ".slk",
	"application/smil":                    ".smi",
	"application/x-smk":                   ".smk",
	"text/plain":                          ".sol",
	"application/futuresplash":            ".spl",
	"application/streamingmedia":          ".ssm",
	"application/vnd.ms-pki.certstore":    ".sst",
	"application/vnd.ms-pki.stl":          ".stl",
	"application/x-sty":                   ".sty",
	"application/x-tdf":                   ".tdf",
	"application/x-tg4":                   ".tg4",
	"application/x-tga":                   ".tga",
	"image/tiff":                          ".tif",
	"application/x-tif":                   ".tif",
	"drawing/x-top":                       ".top",
	"application/x-bittorrent":            ".torrent",
	"application/x-icq":                   ".uin",
	"text/iuls":                           ".uls",
	"text/x-vcard":                        ".vcf",
	"application/x-vda":                   ".vda",
	"application/vnd.visio":               ".vdx",
	"application/x-vpeg005":               ".vpg",
	"application/x-vsd":                   ".vsd",
	"application/x-vst":                   ".vst",
	"audio/wav":                           ".wav",
	"audio/x-ms-wax":                      ".wax",
	"application/x-wb1":                   ".wb1",
	"application/x-wb2":                   ".wb2",
	"application/x-wb3":                   ".wb3",
	"image/vnd.wap.wbmp":                  ".wbmp",
	"application/x-wk3":                   ".wk3",
	"application/x-wk4":                   ".wk4",
	"application/x-wkq":                   ".wkq",
	"application/x-wks":                   ".wks",
	"video/x-ms-wm":                       ".wm",
	"audio/x-ms-wma":                      ".wma",
	"application/x-ms-wmd":                ".wmd",
	"application/x-wmf":                   ".wmf",
	"text/vnd.wap.wml":                    ".wml",
	"video/x-ms-wmv":                      ".wmv",
	"video/x-ms-wmx":                      ".wmx",
	"application/x-ms-wmz":                ".wmz",
	"application/x-wp6":                   ".wp6",
	"application/x-wpd":                   ".wpd",
	"application/x-wpg":                   ".wpg",
	"application/vnd.ms-wpl":              ".wpl",
	"application/x-wq1":                   ".wq1",
	"application/x-wr1":                   ".wr1",
	"application/x-wri":                   ".wri",
	"application/x-wrk":                   ".wrk",
	"application/x-ws":                    ".ws",
	"text/scriptlet":                      ".wsc",
	"video/x-ms-wvx":                      ".wvx",
	"application/vnd.adobe.xdp":           ".xdp",
	"application/vnd.adobe.xfd":           ".xfd",
	"application/vnd.adobe.xfdf":          ".xfdf",
	"application/vnd.ms-excel":            ".xls",
	"application/x-xls":                   ".xls",
	"application/x-xlw":                   ".xlw",
	"application/x-xwd":                   ".xwd",
	"application/x-x_b":                   ".x_b",
	"application/x-x_t":                   ".x_t",
	"application/json":                    ".json",
	"text/x-json":                         ".json",
	"application/andrew-inset":            ".ez",
	"application/mac-compactpro":          ".cpt",
	"application/oda":                     ".oda",
	"application/vnd.mif":                 ".mif",
	"application/vnd.wap.wbxml":           ".wbxml",
	"application/vnd.wap.wmlc":            ".wmlc",
	"application/vnd.wap.wmlscriptc":      ".wmlsc",
	"application/x-bcpio":                 ".bcpio",
	"application/x-cdlink":                ".vcd",
	"application/x-chess-pgn":             ".pgn",
	"application/x-cpio":                  ".cpio",
	"application/x-csh":                   ".csh",
	"application/x-director":              ".dcr",
	"application/x-dvi":                   ".dvi",
	"application/x-futuresplash":          ".spl",
	"application/x-gtar":                  ".gtar",
	"application/x-hdf":                   ".hdf",
	"application/x-koan":                  ".skp",
	"application/x-sh":                    ".sh",
	"application/x-shar":                  ".shar",
	"application/x-sv4cpio":               ".sv4cpio",
	"application/x-sv4crc":                ".sv4crc",
	"application/x-tar":                   ".tar",
	"application/x-tcl":                   ".tcl",
	"application/x-tex":                   ".tex",
	"application/x-texinfo":               ".texinfo",
	"application/x-troff":                 ".t",
	"application/x-troff-me":              ".me",
	"application/x-troff-ms":              ".ms",
	"application/x-ustar":                 ".ustar",
	"application/x-wais-source":           ".src",
	"application/zip":                     ".zip",
	"audio/midi":                          ".mid",
	"audio/mpeg":                          ".mpga",
	"audio/x-aiff":                        ".aif",
	"audio/x-mpegurl":                     ".m3u",
	"audio/x-realaudio":                   ".ra",
	"audio/x-wav":                         ".wav",
	"chemical/x-pdb":                      ".pdb",
	"chemical/x-xyz":                      ".xyz",
	"image/bmp":                           ".bmp",
	"image/ief":                           ".ief",
	"image/vnd.djvu":                      ".djvu",
	"image/x-cmu-raster":                  ".ras",
	"image/x-portable-anymap":             ".pnm",
	"image/x-portable-bitmap":             ".pbm",
	"image/x-portable-graymap":            ".pgm",
	"image/x-portable-pixmap":             ".ppm",
	"image/x-rgb":                         ".rgb",
	"image/x-xbitmap":                     ".xbm",
	"image/x-xpixmap":                     ".xpm",
	"image/x-xwindowdump":                 ".xwd",
	"model/iges":                          ".igs",
	"model/mesh":                          ".msh",
	"model/vrml":                          ".wrl",
	"text/richtext":                       ".rtx",
	"text/rtf":                            ".rtf",
	"text/sgml":                           ".sgml",
	"text/tab-separated-values":           ".tsv",
	"text/vnd.wap.wmlscript":              ".wmls",
	"text/x-setext":                       ".etx",
	"video/quicktime":                     ".qt",
	"video/vnd.mpegurl":                   ".mxu",
	"video/x-msvideo":                     ".avi",
	"x-conference/x-cooltalk":             ".ice",
}

// SocketProxy 启动一个socket5代理
func SocketProxy(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		gt.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			gt.Panic(err)
		}
		go handleClientRequest2(client)
	}
}

func handleClientRequest2(client net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			gt.Error(err)
		}
	}()
	gt.Info("socket 请求 : ", client.RemoteAddr(), " --> ", client.LocalAddr())
	if client == nil {
		return

	}
	defer func() {
		_ = client.Close()
	}()
	var b [1024 * 100]byte
	n, err := client.Read(b[:])
	if err != nil {
		gt.Error(err)
		return
	}
	if b[0] == 0x05 { //只处理Socks5协议
		//客户端回应：Socks服务端不需要验证方式
		_, _ = client.Write([]byte{0x05, 0x00})
		n, err = client.Read(b[:])
		var host, port string

		switch b[3] {
		case 0x01: //IP V4
			host = net.IPv4(b[4], b[5], b[6], b[7]).String()

		case 0x03: //域名
			host = string(b[5 : n-2]) //b[4]表示域名的长度

		case 0x04: //IP V6
			host = net.IP{b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15], b[16], b[17], b[18], b[19]}.String()
		}

		port = strconv.Itoa(int(b[n-2])<<8 | int(b[n-1]))

		server, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			gt.Error(err)
			return

		}
		defer func() {
			_ = server.Close()
		}()
		_, _ = client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) //响应客户端连接成功
		//进行转发
		go func() {
			_, _ = io.Copy(server, client)
		}()
		_, _ = io.Copy(client, server)
	}
}

func main3() {
	// 原始数据
	originalData := "Hello, World!"
	// 创建一个字符串Reader
	reader := strings.NewReader(originalData)

	// 创建一个bytes.Buffer作为TeeReader的第二个Writer
	var buffer bytes.Buffer

	// 创建TeeReader
	readerWithBuffer := io.TeeReader(reader, &buffer)

	// 读取数据
	copiedData, err := io.ReadAll(readerWithBuffer)
	if err != nil {
		panic(err)
	}

	// 打印从originalData复制过来的数据
	println(string(copiedData))

	// 打印buffer内容，应该与originalData相同
	println(buffer.String())
}
