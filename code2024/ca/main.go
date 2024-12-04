package main

import (
	"bytes"
	cr "crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"
)

type myhandler struct {
}

func (h *myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("ServeHTTP...")
	fmt.Fprintf(w, "Hi, This is an example of http service in golang!\n")
}

func main() {
	//main1()

	//mains()

	//main2()

	mainc()
}

func mains() {

	pool := x509.NewCertPool()
	caCertPath := "conf/ca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	s := &http.Server{
		Addr:    "0.0.0.0:16061",
		Handler: &myhandler{},
		TLSConfig: &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	log.Println("启动服务")

	err = s.ListenAndServeTLS("conf/server.crt", "conf/server.key")
	if err != nil {
		fmt.Println("ListenAndServeTLS err:", err)
	}
}

func mainc() {
	pool := x509.NewCertPool()
	caCertPath := "conf/ca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}

	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair("conf/server.crt", "conf/server.key")
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	client := &http.Client{Transport: tr}
	//这里的ip地址需要在生成自签名证书的时候指定,否则ssl验证不通过。
	resp, err := client.Get("https://127.0.0.1:16061")
	if err != nil {
		log.Println("err = ", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func main2() {
	// 确保你的证书文件(cert.pem和key.pem)在可执行文件同一目录下
	cert, err := tls.LoadX509KeyPair("ssl/RootCA_Cert.pem", "ssl/RootCA_Cert_Key.pem")
	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	cfg.InsecureSkipVerify = true

	srv := &http.Server{
		Addr:      ":8181",
		Handler:   http.DefaultServeMux,
		TLSConfig: cfg,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var reqBodyBytes []byte
		reqBodyBytes, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		log.Println("请求 body:  ", string(reqBodyBytes))
		log.Println("请求 Method :  ", r.Method)
		log.Println("请求 : ", r)

		log.Println(r.Body)

		_, _ = w.Write([]byte("Hello, this is an example server over HTTPS!"))
	})

	log.Println("Serving at https://localhost/")
	err = srv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}

func main1() {
	subj := &pkix.Name{
		CommonName:    "127.0.0.1:16061",
		Organization:  []string{"aac Company, INC."},
		Country:       []string{"US"},
		Province:      []string{""},
		Locality:      []string{"San Francisco"},
		StreetAddress: []string{"Golden Gate Bridge"},
		PostalCode:    []string{"94016"},
	}
	ca, err := CreateCA(subj, 10)
	if err != nil {
		log.Panic(err)
	}

	Write(ca, "./ca")

	crt, err := Req(ca.CSR, subj, 1, []string{"127.0.0.1"}, []net.IP{
		net.IPv4(127, 0, 0, 1),
		net.IPv4(0, 0, 0, 0),
	})

	if err != nil {
		log.Panic(err)
	}

	Write(crt, "./tls")

	Write(crt, "./client")
}

type CERT struct {
	CERT       []byte
	CERTKEY    *rsa.PrivateKey
	CERTPEM    *bytes.Buffer
	CERTKEYPEM *bytes.Buffer
	CSR        *x509.Certificate
}

func CreateCA(sub *pkix.Name, expire int) (*CERT, error) {
	var (
		ca  = new(CERT)
		err error
	)

	if expire < 1 {
		expire = 1
	}
	// 为ca生成私钥
	ca.CERTKEY, err = rsa.GenerateKey(cr.Reader, 4096)
	if err != nil {
		return nil, err
	}

	// 对证书进行签名
	ca.CSR = &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63n(2000)),
		Subject:      *sub,
		NotBefore:    time.Now(),                       // 生效时间
		NotAfter:     time.Now().AddDate(expire, 0, 0), // 过期时间
		IsCA:         true,                             // 表示用于CA
		// openssl 中的 extendedKeyUsage = clientAuth, serverAuth 字段
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		// openssl 中的 keyUsage 字段
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	// 创建证书
	// caBytes 就是生成的证书
	ca.CERT, err = x509.CreateCertificate(cr.Reader, ca.CSR, ca.CSR, &ca.CERTKEY.PublicKey, ca.CERTKEY)
	if err != nil {
		return nil, err
	}
	ca.CERTPEM = new(bytes.Buffer)
	pem.Encode(ca.CERTPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca.CERT,
	})
	ca.CERTKEYPEM = new(bytes.Buffer)
	pem.Encode(ca.CERTKEYPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(ca.CERTKEY),
	})

	// 进行PEM编码，编码就是直接cat证书里面内容显示的东西
	return ca, nil
}

func Req(ca *x509.Certificate, sub *pkix.Name, expire int, dns []string, ip []net.IP) (*CERT, error) {
	var (
		cert = &CERT{}
		err  error
	)
	cert.CERTKEY, err = rsa.GenerateKey(cr.Reader, 4096)
	if err != nil {
		return nil, err
	}
	if expire < 1 {
		expire = 1
	}
	cert.CSR = &x509.Certificate{
		SerialNumber: big.NewInt(rand.Int63n(2000)),
		Subject:      *sub,
		IPAddresses:  ip,
		DNSNames:     dns,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(expire, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	cert.CERT, err = x509.CreateCertificate(cr.Reader, cert.CSR, ca, &cert.CERTKEY.PublicKey, cert.CERTKEY)
	if err != nil {
		return nil, err
	}

	cert.CERTPEM = new(bytes.Buffer)
	pem.Encode(cert.CERTPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.CERT,
	})
	cert.CERTKEYPEM = new(bytes.Buffer)
	pem.Encode(cert.CERTKEYPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(cert.CERTKEY),
	})
	return cert, nil
}

func Write(cert *CERT, file string) error {
	keyFileName := file + ".key"
	certFIleName := file + ".crt"
	kf, err := os.Create(keyFileName)
	if err != nil {
		return err
	}
	defer kf.Close()

	if _, err := kf.Write(cert.CERTKEYPEM.Bytes()); err != nil {
		return err
	}

	cf, err := os.Create(certFIleName)
	if err != nil {
		return err
	}
	if _, err := cf.Write(cert.CERTPEM.Bytes()); err != nil {
		return err
	}
	return nil
}
