package main

import (
	"bytes"
	"fmt"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"log"
	"net"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
)

func main() {
	//case1()

	//case2()

	//case4()

	//case5()

	//case6()

	//case7()

	case9()

	case11()

	case12()
}

// case1  golang.org/x/sys/windows 是 Go 语言中用于与 Windows 操作系统进行交互的标准库扩展包，它提供了许多与 Windows
// 系统底层功能相关的函数和类型，允许 Go 程序调用 Windows 特定的 API，从而实现与 Windows 操作系统的深度集成。
// 使用场景
// 系统级编程：当需要编写与 Windows 操作系统底层功能紧密相关的程序时，如文件系统操作、进程管理、服务控制等，golang.org/x/sys/windows 库可以提供必要的接口。
// 图形界面开发：在开发 Windows 图形界面应用程序时，可以使用该库调用 Windows 的图形 API，如 GDI（Graphics Device Interface），实现图形绘制、窗口管理等功能。
// COM 组件集成：如果需要在 Go 程序中集成 COM 组件，例如使用 Office 自动化功能，golang.org/x/sys/windows 库提供了与 COM 接口交互的能力。
// 网络编程：在进行 Windows 平台上的网络编程时，该库可以用于调用 Windows 的网络 API，实现网络连接、数据传输等功能。
func case1() {
	var size uint32 = 256
	buf := make([]uint16, size)
	err := windows.GetComputerName(&buf[0], &size)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Computer Name:", syscall.UTF16ToString(buf))
	}
}

// case2 创建windows守护进程
// 定义服务结构体
type myService struct{}

// 实现服务执行逻辑
func (m *myService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	// 定义服务可接受的控制命令
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	// 设置定时器，每隔30秒触发一次
	tick := time.Tick(30 * time.Second)

	// 发送服务启动信号
	status <- svc.Status{State: svc.StartPending}
	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	// 主循环，处理定时器事件和控制命令
loop:
	for {
		select {
		case <-tick:
			// 处理定时器事件，记录日志
			log.Println("Tick Handled...!")
		case c := <-r:
			// 处理控制命令
			switch c.Cmd {
			case svc.Interrogate:
				// 发送服务当前状态
				status <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// 停止服务
				log.Println("Shutting service...!")
				break loop
			case svc.Pause:
				// 暂停服务
				status <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				// 恢复服务
				status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				// 处理未知控制命令
				log.Printf("Unexpected service control request #%d", c)
			}
		}
	}

	// 发送服务停止信号
	status <- svc.Status{State: svc.StopPending}
	return false, 1
}

// 运行服务
func runService(name string, isDebug bool) {
	if isDebug {
		// 调试模式
		err := debug.Run(name, &myService{})
		if err != nil {
			log.Fatalln("Error running service in debug mode:", err)
		}
	} else {
		// 服务控制模式
		err := svc.Run(name, &myService{})
		if err != nil {
			log.Fatalln("Error running service in SC mode:", err)
		}
	}
}

func case2() {
	// 运行服务
	runService("./main.exe", true) // 将第二个参数改为true可进入调试模式
}

// case3 在 Win 环境下设置进程优先级的函数是 BOOL SetPriorityClass(HANDLE hProcess, DWORD dwPriorityClass)。
func case3() {
	cur_task_handle := windows.CurrentProcess()
	err := windows.SetPriorityClass(cur_task_handle, windows.HIGH_PRIORITY_CLASS)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		time.Sleep(10 * time.Second)
	}
}

// case4 操作Windows注册表
func case4() {

	// 打开根键，这里以 HKEY_CURRENT_USER 为例，你也可以选择其他根键如 HKEY_LOCAL_MACHINE 等
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\MyApp`, registry.ALL_ACCESS)
	if err != nil {
		if err == registry.ErrNotExist {
			// 如果键不存在，则尝试创建
			newKey, _, createErr := registry.CreateKey(registry.CURRENT_USER, `SOFTWARE\MyApp`, registry.ALL_ACCESS)
			if createErr != nil {
				fmt.Printf("创建注册表键失败: %v\n", createErr)
				return
			}
			defer newKey.Close()
			fmt.Println("注册表键创建成功")
		} else {
			fmt.Printf("打开注册表键失败: %v\n", err)
			return
		}
	} else {
		defer key.Close()
		fmt.Println("注册表键已存在")
	}

	//key, exists, err := registry.CreateKey(registry.CLASSES_ROOT, "\\HKEY_LOCAL_MACHINE\\SOFTWARE\\Hello Go\\", registry.ALL_ACCESS)
	//defer key.Close()
	//if err != nil {
	//	log.Println("err", err.Error())
	//}
	//// 判断是否已经存在了
	//if exists {
	//	println(`键已存在`)
	//} else {
	//	println(`新建注册表键`)
	//}

	// 写入：32位整形值
	err = key.SetDWordValue(`32位整形值`, uint32(123456))
	if err != nil {
		log.Println("写入：32位整形值 err", err.Error())
	}
	// 写入：64位整形值
	err = key.SetQWordValue(`64位整形值`, uint64(123456))
	if err != nil {
		log.Println("写入：64位整形值 err", err.Error())
	}
	// 写入：字符串
	err = key.SetStringValue(`字符串`, `hello`)
	if err != nil {
		log.Println("写入：字符串 err", err.Error())
	}
	// 写入：字符串数组
	err = key.SetStringsValue(`字符串数组`, []string{`hello`, `world`})
	if err != nil {
		log.Println("写入：字符串数组 err", err.Error())
	}
	// 写入：二进制
	err = key.SetBinaryValue(`二进制`, []byte{0x11, 0x22})
	if err != nil {
		log.Println("写入：二进制 err", err.Error())
	}

	// 读取：字符串
	s, _, err := key.GetStringValue(`字符串`)
	if err != nil {
		log.Println("读取：字符串 err", err.Error())
	}
	println(s)

	// 读取：一个项下的所有子项
	keys, err := key.ReadSubKeyNames(0)
	if err != nil {
		log.Println("读取：一个项下的所有子项 err", err.Error())
	}
	for _, key_subkey := range keys {
		// 输出所有子项的名字
		println(key_subkey)
	}

	// 创建：子项
	subkey, _, err := registry.CreateKey(key, `子项`, registry.ALL_ACCESS)
	if err != nil {
		log.Println("创建：子项 err", err.Error())
	}
	defer subkey.Close()

	// 删除：子项
	// 该键有子项，所以会删除失败
	// 没有子项，删除成功
	//registry.DeleteKey(key, `子项`)
}

// case5 windows单实例互斥锁
func case5() {
	name := "a"
	ptr, _ := syscall.UTF16PtrFromString(name)
	_, err := windows.CreateMutex(nil, true, ptr)
	if err != nil {
		log.Println("pr1 ", err.Error())
	}

	ptr2, _ := syscall.UTF16PtrFromString(name)
	_, err = windows.CreateMutex(nil, true, ptr2)
	if err != nil {
		log.Println("pr2 ", err.Error())
	}
}

// 指定字符分割 常用于邮箱分割  指定字符第一个分割
func case6() {
	splitTenantAndClientID := func(user string) (string, string) {
		at := strings.IndexRune(user, '@')
		if at < 1 || at >= (len(user)-1) {
			return user, ""
		}

		return user[0:at], user[at+1:]
	}

	log.Println(splitTenantAndClientID("123@qq.com"))
}

// case7 分割路由最后一个/
func case7() {
	splitAuthorityAndTenant := func(authorityURL string) (string, string) {
		separatorIndex := strings.LastIndex(authorityURL, "/")
		tenant := authorityURL[separatorIndex+1:]
		authority := authorityURL[:separatorIndex]
		return authority, tenant
	}

	log.Println(splitAuthorityAndTenant("/a/b/c/d"))
}

// case8 比较两个字符是否相等
func case8() {
	a := []byte{1, 2, 3, 4}
	b := []byte{1, 2, 3, 4}
	c := []byte{1, 2, 5, 4}

	fmt.Println(bytes.Equal(a, b)) // 输出: true
	fmt.Println(bytes.Equal(a, c)) // 输出: false
}

// case9 转字符串
func case9() {
	log.Println(asString(1230))
	log.Println(asString(false))
	log.Println(asString([]byte("aaa")))
}

func asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	return fmt.Sprintf("%v", src)
}

// case10 封装一个简单的超时连接
type timeoutConn struct {
	c       net.Conn
	timeout time.Duration
}

func newTimeoutConn(conn net.Conn, timeout time.Duration) *timeoutConn {
	return &timeoutConn{
		c:       conn,
		timeout: timeout,
	}
}

func (c *timeoutConn) Read(b []byte) (n int, err error) {
	if c.timeout > 0 {
		err = c.c.SetDeadline(time.Now().Add(c.timeout))
		if err != nil {
			return
		}
	}
	return c.c.Read(b)
}

func (c *timeoutConn) Write(b []byte) (n int, err error) {
	if c.timeout > 0 {
		err = c.c.SetDeadline(time.Now().Add(c.timeout))
		if err != nil {
			return
		}
	}
	return c.c.Write(b)
}

func (c timeoutConn) Close() error {
	return c.c.Close()
}

func (c timeoutConn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

func (c timeoutConn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

func (c timeoutConn) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

func (c timeoutConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

func (c timeoutConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}

// case11  字符串转换为UTF-16编码的[]字节 小端 littleEndian  自行实现
func str2ucs2(s string) []byte {
	res := utf16.Encode([]rune(s))
	ucs2 := make([]byte, 2*len(res))
	for i := 0; i < len(res); i++ {
		ucs2[2*i] = byte(res[i])
		ucs2[2*i+1] = byte(res[i] >> 8)
	}
	return ucs2
}

func case11() {
	log.Println(str2ucs2("你好"))
	s := str2ucs2("你好")
	log.Println(string(s))
	log.Println(string(str2ucs2("你好12312>?<?M}")))
}

// case12  Sscanf 解析指定值
func case12() {
	var driverVersion = "v1.8.0"
	var majorVersion uint32
	var minorVersion uint32
	var rev uint32
	_, _ = fmt.Sscanf(driverVersion, "v%d.%d.%d", &majorVersion, &minorVersion, &rev)
	log.Println(majorVersion)
	log.Println(minorVersion)
	log.Println(rev)
}
