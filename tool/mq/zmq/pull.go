package zmq

//func ReceiveMessage(addr string, handler func(msg []byte), ctx context.Context) error {
//	receiver, err := zmq4.NewSocket(zmq4.PULL)
//	if err != nil {
//		return err
//	}
//	err = receiver.Connect(addr)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		_ = receiver.SetLinger(-1)
//		_ = receiver.Close()
//	}()
//	for {
//		select {
//		case <-ctx.Done():
//		default:
//			msg, err := receiver.RecvBytes(0)
//			if err != nil {
//				return err
//			}
//			handler(msg)
//		}
//	}
//}
