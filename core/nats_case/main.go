package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

func main() {
	go func() {
		NatEncodedPublish()
	}()

	go func() {
		time.Sleep(5 * time.Second)
		NatEncodedSubscribe()
	}()
	select {}
}

// NatEncodedSubscribe 订阅
func NatEncodedSubscribe() {
	natsConn, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		fmt.Errorf("nats连接错误%#v", err)
		return
	}
	c, _ := nats.NewEncodedConn(natsConn, "json")
	_, err = c.Subscribe("rd", func(m *nats.Msg) {
		fmt.Printf("接收到nats消息: %+v\n", string(m.Data))
	})
	if err != nil {
		log.Println(err)
	}
}

// NatEncodedPublish 发布
func NatEncodedPublish() {
	natsConn, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		fmt.Errorf("nats连接错误%#v", err)
		return
	}
	c, _ := nats.NewEncodedConn(natsConn, "json")
	for {
		time.Sleep(1 * time.Second)
		err := c.Publish("rd", "data")
		if err != nil {
			fmt.Println("消息发布错误：", err.Error())
		}
		fmt.Println("消息发送完毕")
	}
}
