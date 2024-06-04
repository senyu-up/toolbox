package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	"github.com/senyu-up/toolbox/tool/config"
)

type ConnLogger struct {
	sync.Mutex
	innerWriter    io.WriteCloser
	ReconnectOnMsg bool   `json:"reconnectOnMsg"`
	Reconnect      bool   `json:"reconnect"`
	Net            string `json:"net"`
	Addr           string `json:"addr"`
	Level          string `json:"level"`
	LogLevel       LogLevel
	illNetFlag     bool //网络异常标记
}

func (c *ConnLogger) Init(jsonConfig string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	fmt.Printf("consoleWriter Init:%s\n", jsonConfig)
	err := json.Unmarshal([]byte(jsonConfig), c)
	if err != nil {
		return err
	}
	if l, ok := LevelMap[c.Level]; ok {
		c.LogLevel = l
	}
	if c.innerWriter != nil {
		c.innerWriter.Close()
		c.innerWriter = nil
	}
	return nil
}

func (c *ConnLogger) InitByConf(conf config.ConnConfig) (err error) {
	fmt.Printf("consoleWriter Init:%+v\n", conf)
	copier.CopyWithOption(c, conf, copier.Option{IgnoreEmpty: true})
	if err != nil {
		return err
	}
	if l, ok := LevelMap[c.Level]; ok {
		c.LogLevel = l
	} else {
		return ErrInvalidLogLevel
	}
	if c.innerWriter != nil {
		c.innerWriter.Close()
		c.innerWriter = nil
	}
	return nil
}

func (c *ConnLogger) LogWrite(when time.Time, msg string, level LogLevel, extras []Field) (err error) {
	if level > c.LogLevel {
		return nil
	}
	msgObj := fieldsToLogInfo(extras)

	if c.needToConnectOnMsg() {
		err = c.connect()
		if err != nil {
			return
		}
		//重连成功
		c.illNetFlag = false
	}

	//每条消息都重连一次日志中心，适用于写日志频率极低的情况下的服务调用,避免长时间连接，占用资源
	if c.ReconnectOnMsg { // 频繁日志发送切勿开启
		defer c.innerWriter.Close()
	}

	//网络异常时，消息发出
	if !c.illNetFlag {
		err = c.println(when, msgObj) // TODO
		//网络异常，通知处理网络的go程自动重连
		if err != nil {
			c.illNetFlag = true
		}
	}

	return
}

func (c *ConnLogger) Destroy() {
	if c.innerWriter != nil {
		c.innerWriter.Close()
	}
}

func (c *ConnLogger) Name() string {
	return AdapterConn
}

func (c *ConnLogger) connect() error {
	if c.innerWriter != nil {
		c.innerWriter.Close()
		c.innerWriter = nil
	}
	addrs := strings.Split(c.Addr, ";")
	for _, addr := range addrs {
		conn, err := net.Dial(c.Net, addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "net.Dial error:%v\n", err)
			continue
			//return err
		}

		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
		}
		c.innerWriter = conn
		return nil
	}
	return fmt.Errorf("hava no valid logs service addr:%v", c.Addr)
}

func (c *ConnLogger) needToConnectOnMsg() bool {
	if c.Reconnect {
		c.Reconnect = false
		return true
	}

	if c.innerWriter == nil {
		return true
	}

	if c.illNetFlag {
		return true
	}
	return c.ReconnectOnMsg
}

func (c *ConnLogger) CurrentLevel() LogLevel {
	return c.LogLevel
}

func (c *ConnLogger) println(when time.Time, msg *loginfo) error {
	c.Lock()
	defer c.Unlock()
	ss, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = c.innerWriter.Write(append(ss, '\n'))

	//返回err，解决日志系统网络异常后的自动重连
	return err
}

func fieldsToLogInfo(fields []Field) *loginfo {
	msg := &loginfo{}
	for _, field := range fields {
		switch field.Key {
		case "name":
			msg.Name = field.String
		case "path":
			msg.Path = field.String
		case "level":
			msg.Level = levelPrefix[field.Integer]
		case "time":
			msg.Time = time.Unix(field.Integer, 0).Format("2006-01-02 15:04:05")
		case "content":
			msg.Content = field.String
		}
	}
	return msg
}
