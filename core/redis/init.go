//
//	初始化redis连接
//

package manredis

import (
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// 	定义错误 ErrRedisIPIsEmpty 连接传入的RedisIP为空　则使用这个错误
//	使用该错误信息的地方需要panic()
var ErrRedisIPIsEmpty = errors.New("input redis ip is empty.")

//	定义错误 ErrSSHInfoIsEmpty 使用了ssh通道连接redis，没有传入连接ssh的主机，密码　则使用这个错误
//	使用该错误信息的地方需要panic()
var ErrSSHInfoIsEmpty = errors.New("input ssh info is empty.")

//	连接失败定义的错误
var ErrConnFailed = errors.New("connect redis failed.")

//	连接redis的结构体定义
type RedisGo struct {

	//	redis连接地址
	RedisIP string

	//	redis连接端口,默认 6379
	RedisPort int

	//	redisn连接密码
	RedisPassword string

	//	使用ssh通道连接redis的 host地址
	SSHAddr string

	//	使用ssh通道连接redis的 host 账号
	SSHUser string

	//	使用ssh通道连接redis的 host 密码
	SSHPass string

	//	全局设置redis　db, 默认连接为0, 非单独使用所强制指定可以不用在这里设置
	GlobalDBNumber int

	//	最大闲置数，用于redis连接池
	RedisMaxIdle int

	//	最大连接数，
	RedisMaxActive int

	//	单条连接Timeout
	RedisIdleTimeoutSec int
}

func (r *RedisGo) init() {
	if r.RedisIP == "" {
		panic(ErrRedisIPIsEmpty)
	}

	if r.RedisPort == 0 {
		r.RedisPort = 6379
	}

	if r.RedisMaxActive == 0 {
		r.RedisMaxActive = 20
	}

	if r.RedisMaxActive == 0 {
		r.RedisMaxActive = 20
	}

	if r.RedisIdleTimeoutSec == 0 {
		r.RedisIdleTimeoutSec = 20
	}

}

//	connSSH 连接ssh
func (r *RedisGo) connSSH() (*ssh.Client, error) {
	if r.SSHAddr == "" || r.SSHPass == "" {
		panic(ErrSSHInfoIsEmpty)
	}

	if r.SSHUser == "" {
		r.SSHUser = "root"
	}

	config := &ssh.ClientConfig{
		User: r.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.RedisPassword),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	sshConn, err := net.Dial("tcp", r.SSHAddr)
	if nil != err {
		fmt.Println("net dial err: ", err)
		return nil, err
	}

	clientConn, chans, reqs, err := ssh.NewClientConn(sshConn, r.SSHAddr, config)
	if nil != err {
		sshConn.Close()
		fmt.Println("ssh client conn err: ", err)
		return nil, err
	}

	client := ssh.NewClient(clientConn, chans, reqs)
	return client, nil
}

// 	RConn 普通连接
// 	返回redis连接
func (r *RedisGo) RConn() (redis.Conn, error) {
	host := fmt.Sprintf("%s:%d", r.RedisIP, r.RedisPort)
	conn, err := redis.Dial("tcp", host)
	if nil != err {
		fmt.Println("dial to redis addr err: ", err)
		return nil, err
	}
	if r.RedisPassword != "" {
		if _, authErr := conn.Do("AUTH", r.RedisPassword); authErr != nil {
			fmt.Println("redis auth password error: ", authErr)
			return nil, fmt.Errorf("redis auth password error: %s", authErr)
		}
	}
	return conn, nil
}

// 	RSSHConn SSH普通连接
// 	返回redis连接
func (r *RedisGo) RSSHConn() (redis.Conn, error) {

	sshClient, err := r.connSSH()
	if nil != err {
		fmt.Println(err)
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", r.RedisIP, r.RedisPort)
	conn, err := sshClient.Dial("tcp", host)
	if nil != err {
		fmt.Println("dial to redis addr err: ", err)
		return nil, err
	}

	redisConn := redis.NewConn(conn, -1, -1)

	fmt.Println("RSSHConn SSH普通连接　　＝　", r.RedisPassword)

	if r.RedisPassword != "" {
		if _, authErr := redisConn.Do("AUTH", r.RedisPassword); authErr != nil {
			fmt.Println("redis auth password error: ", authErr)
			return nil, fmt.Errorf("redis auth password error: %s", authErr)
		}
	}

	return redisConn, nil
}

// RPool 连接池连接
// 返回redis连接池  *redis.Pool.Get() 获取redis连接
func (r *RedisGo) RPool() *redis.Pool {
	redisURL := fmt.Sprintf("redis://%s:%d", r.RedisIP, r.RedisPort)
	fmt.Println("redisURL -> ", redisURL)
	return &redis.Pool{
		MaxIdle:     RedisMaxIdle,
		MaxActive:   RedisMaxActive,
		IdleTimeout: time.Duration(RedisIdleTimeoutSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}

			fmt.Println("RPool 连接池连接　　＝　", r.RedisPassword)

			//验证redis密码
			if r.RedisPassword != "" {
				if _, authErr := c.Do("AUTH", r.RedisPassword); authErr != nil {
					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}

			return nil
		},
	}
}

// RSSHPool SSH连接池连接
// addr : SSH主机地址, 如: 127.0.0.1:22
// user : SSH用户
// pass : SSH密码
// ip : redis服务地址
// port:  Redis 服务端口
// password  Redis 服务密码
// 配置参数  RedisMaxIdle 最大连接
// 配置参数  RedisMaxActive 最大连接数
// 配置参数  RedisIdleTimeoutSec 设置超时
// 返回redis连接池  调用: c := RSSHPool().Get() 返回redis连接
func　(r *RedisGo) RSSHPool() *redis.Pool {
	sshClient, err := r.connSSH()
	if nil != err {
		fmt.Println(err)
		return nil
	}

	redisURL := fmt.Sprintf("%s:%d", r.RedisIP, r.RedisPort)
	return &redis.Pool{
		MaxIdle:     RedisMaxIdle,
		MaxActive:   RedisMaxActive,
		IdleTimeout: time.Duration(RedisIdleTimeoutSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := sshClient.Dial("tcp", redisURL)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}

			fmt.Println(c)
			redisc := redis.NewConn(c, -1, -1)

			fmt.Println("RSSHPool SSH连接池连接　", r.RedisPassword)

			if r.RedisPassword != "" {
				if _, authErr := redisc.Do("AUTH", r.RedisPassword); authErr != nil {
					return nil, fmt.Errorf("redis auth password error: %s", authErr)
				}
			}

			return redisc, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}

			return nil
		},
	}
}
