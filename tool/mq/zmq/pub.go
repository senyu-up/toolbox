package zmq

//import (
//	"context"
//	"fmt"
//	"github.com/pebbe/zmq4"
//)
//
////发布订阅模式 ZMQ_PUB
////ZMQ_PUB套接字由于已达到订户的高水位标记而进入静音状态时，将发送给有问题的订户的任何消息都将被丢弃，直到静音状态结束为止
////通过应答机制，告知客户端服务端是否处于开启状态
//
//func PublishMessage(addr, listener string, topic string, sender <-chan string, ctx context.Context) error {
//	zmqctx, err := zmq4.NewContext()
//	if err != nil {
//		return err
//	}
//	defer func() {
//		_ = zmqctx.Term()
//	}()
//	go func() {
//		//应答控制
//		sync, err := zmqctx.NewSocket(zmq4.REP)
//		if err != nil {
//			fmt.Println("Publish listener err:", err)
//			return
//		}
//		err = sync.Bind(listener)
//		if err != nil {
//			fmt.Println("Publish listener err:", err)
//			return
//		}
//		defer func() {
//			_ = sync.Close()
//		}()
//		for {
//			subscriber, err := sync.Recv(0)
//			if err != nil {
//				continue
//			}
//			_, err = sync.Send(fmt.Sprintf("Successfully Subscribed: %s", subscriber), 0)
//		}
//	}()
//	pub, err := zmqctx.NewSocket(zmq4.PUB)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		//设置socket 延时关闭时间
//		_ = pub.SetLinger(-1)
//		_ = pub.Close()
//	}()
//	err = pub.Bind(addr)
//	if err != nil {
//		return err
//	}
//	for {
//		select {
//		case msg := <-sender:
//			_, err = pub.Send(fmt.Sprintf("%s:%s", topic, msg), 0)
//			if err != nil {
//				return err
//			}
//			fmt.Printf("topic: %s;send message success: %s\n", topic, msg)
//		case <-ctx.Done():
//			return nil
//		}
//	}
//}
