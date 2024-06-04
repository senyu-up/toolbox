package zmq

//
//import (
//	"context"
//	"fmt"
//	"github.com/pebbe/zmq4"
//	"strings"
//)
//
//func SubscribeMessage(addr, listener, topic string, handler func(msg string), ctx context.Context) error {
//	//判断服务器是否运行中
//	req, err := zmq4.NewSocket(zmq4.REQ)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		_ = req.Close()
//	}()
//	err = req.Connect(listener)
//	if err != nil {
//		return err
//	}
//	_, err = req.Send(topic, 0)
//	if err != nil {
//		return err
//	}
//	rsp, err := req.Recv(0)
//	if err != nil {
//		return err
//	}
//	fmt.Println(rsp)
//	sub, err := zmq4.NewSocket(zmq4.SUB)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		//设置socket 延时关闭时间
//		_ = sub.SetLinger(-1)
//		_ = sub.Close()
//		_ = sub.SetUnsubscribe(topic)
//	}()
//	err = sub.Connect(addr)
//	if err != nil {
//		return err
//	}
//	err = sub.SetSubscribe(topic)
//	var msg string
//	for {
//		select {
//		case <-ctx.Done():
//		default:
//			msg, err = sub.Recv(0)
//			if err != nil {
//				return err
//			}
//			handler(strings.TrimPrefix(msg, topic+":"))
//		}
//	}
//}
