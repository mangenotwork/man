package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// SSH配置结构体
type SSHConfig struct {
	Addr      string // SSH服务器地址
	User      string // SSH用户名
	Password  string // SSH密码或私钥路径（取决于认证方式）
	Port      int    // SSH端口，默认22
	LocalPort int    // 本地转发端口
}

// 数据库配置结构体
type DBConfig struct {
	User     string // 数据库用户名
	Password string // 数据库密码
	Name     string // 数据库名称
}

func createSSHClient(sshConf SSHConfig) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: sshConf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshConf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境中应避免使用InsecureIgnoreHostKey()
		Timeout:         60 * time.Second,
	}

	return ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshConf.Addr, sshConf.Port), sshConfig)
}

// 建立ssh 隧道
//
//	go SSHSD()
func SSHSD() {
	// 配置参数 - 修改这些值以匹配你的环境
	sshConfig := struct {
		Username   string
		Password   string // 密码认证
		PrivateKey string // 密钥文件路径
		RemoteHost string // SSH服务器地址
		RemotePort int64  // SSH服务器端口
		LocalPort  int64  // 本地监听端口
		TargetHost string // 目标服务器地址
		TargetPort int64  // 目标服务端口
	}{}

	// 配置SSH客户端
	config := &ssh.ClientConfig{
		User: sshConfig.Username,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// 注意：生产环境中应该验证主机密钥
			// 这里为了简化示例，接受所有密钥
			return nil
		},
		Timeout: 10 * time.Second,
	}

	// 设置认证方法
	if sshConfig.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(sshConfig.Password))
	}

	if sshConfig.PrivateKey != "" {
		key, err := os.ReadFile(sshConfig.PrivateKey)
		if err != nil {
			log.Fatalf("读取私钥文件失败: %v", err)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			log.Fatalf("解析私钥失败: %v", err)
		}

		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	if len(config.Auth) == 0 {
		log.Fatal("必须指定密码或私钥文件进行认证")
	}

	// 连接到SSH服务器
	sshAddr := fmt.Sprintf("%s:%d", sshConfig.RemoteHost, sshConfig.RemotePort)
	client, err := ssh.Dial("tcp", sshAddr, config)
	if err != nil {
		log.Fatalf("连接到SSH服务器失败: %v", err)
	}
	defer client.Close()

	log.Printf("已连接到SSH服务器 %s", sshAddr)

	// 本地监听
	localAddr := fmt.Sprintf("127.0.0.1:%d", sshConfig.LocalPort)
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("本地监听失败: %v", err)
	}
	defer listener.Close()

	log.Printf("本地监听已启动: %s", localAddr)
	log.Printf("转发规则: 本地 %s -> 远程 %s:%d",
		localAddr, sshConfig.TargetHost, sshConfig.TargetPort)

	// 设置信号处理，优雅关闭
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
		<-sigchan
		log.Println("接收到关闭信号，正在关闭 ssh 隧道...")
		listener.Close()
		client.Close()
		os.Exit(0)
	}()

	// 接受本地连接并转发
	for {
		localConn, err := listener.Accept()
		if err != nil {
			// 检查是否是因为监听关闭导致的错误
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				log.Printf("临时错误: %v", err)
				continue
			}
			log.Printf("接受连接失败: %v", err)
			break
		}
		defer localConn.Close()
		// 处理每个连接
		go handleConnection(client, localConn, sshConfig.TargetHost, sshConfig.TargetPort)
	}
}

// 处理每个连接的转发
func handleConnection(client *ssh.Client, localConn net.Conn, targetHost string, targetPort int64) {
	// 连接到远程目标
	targetAddr := fmt.Sprintf("%s:%d", targetHost, targetPort)
	remoteConn, err := client.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("连接到远程目标失败: %v", err)
		return
	}
	log.Printf("新连接已建立: %s <-> %s", localConn.RemoteAddr(), targetAddr)

	// 双向数据转发
	go forward(localConn, remoteConn)
	go forward(remoteConn, localConn)
}

// 数据转发函数
func forward(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()

	buf := make([]byte, 32*1024)
	_, err := io.CopyBuffer(dst, src, buf)
	if err != nil {
		log.Printf("转发数据失败: %v", err)
	}
}

//// NewORM  连接 orm
//func NewORM(database, user, password, host, port string, disablePrepared bool) (*gorm.DB, error) {
//
//	var (
//		orm *gorm.DB
//		err error
//	)
//
//	if isSSH, ok := conf.GetString("IsSSH"); ok && isSSH == "y" {
//		go func() {
//			SSHSD()
//		}()
//
//		time.Sleep(3 * time.Second)
//
//		host = "127.0.0.1"
//		localPort, _ := conf.GetInt64("LocalPort")
//		port = utils.AnyToString(localPort)
//	}
//
//	if database == "" || user == "" || password == "" || host == "" {
//		panic("数据库配置信息获取失败")
//	}
//	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database) + "?charset=utf8mb4&parseTime=true&loc=Local"
//	if disablePrepared {
//		str = str + "&interpolateParams=true"
//	}
//
//	orm, err = gorm.Open(mysql.Open(str), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: logger.NewGormLogger(),
//		//Logger: newLogger,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	db, err := orm.DB()
//	if err != nil {
//		return nil, err
//	}
//
//	db.SetMaxIdleConns(10) //空闲连接数
//	db.SetMaxOpenConns(20) //最大连接数
//	db.SetConnMaxLifetime(10 * time.Second)
//	db.SetConnMaxIdleTime(60 * time.Second) // 连接最大空闲时间
//
//	return orm, err
//}
