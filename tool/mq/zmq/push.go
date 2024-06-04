package zmq

//func SendMessage(addr string, sender <-chan []byte, ctx context.Context) error {
//	socket, err := zmq4.NewSocket(zmq4.PUSH)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		_ = socket.SetLinger(-1)
//		_ = socket.Close()
//	}()
//	err = socket.Connect(addr)
//	if err != nil {
//		return err
//	}
//	_ = socket.SetSndhwm(1)
//	for {
//		select {
//		case <-ctx.Done():
//			return nil
//		case msg, ok := <-sender:
//			if !ok {
//				logger.ERR("message sender channel closed!")
//				return nil
//			}
//			_, err = socket.SendBytes(msg, 0)
//			if err != nil {
//				logger.ERR("Send message err: ", err, "\n data: ", string(msg))
//				return err
//			}
//		}
//	}
//}
