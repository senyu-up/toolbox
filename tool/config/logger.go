package config

import (
	"time"
)

// Log 配置
type LogConfig struct {
	AppName    string `json:"-"     yaml:"-"`
	TimeFormat string `json:"TimeFormat"  yaml:"timeformat"`
	CallDepth  int    `json:"CallDepth"   yaml:"calldepth"` // 打印日志时，调用栈深度 跳过多少级
	UsePath    string `json:"UsePath"     yaml:"usepath"`   // 打印日志路径时，去掉的前缀

	// 默认配置项目[console,file,conn,zap]
	// 如果不填，下面三种配置中哪个有值就会用哪个, 如果多种配置都有效，则随机！
	DefaultLog string         `json:"DefaultLog,omitempty" yaml:"defaultlog,omitempty"`
	Console    *ConsoleConfig `json:"Console,omitempty"    yaml:"console,omitempty"`
	File       *FileConfig    `json:"File,omitempty"       yaml:"file,omitempty"`
	Conn       *ConnConfig    `json:"Conn,omitempty"       yaml:"conn,omitempty"`
	Zap        *ZapConfig     `json:"zap,omitempty"        yaml:"zap,omitempty"`
}

type FileConfig struct {
	Filename   string `json:"filename"`
	Append     bool   `json:"append"`
	MaxLines   int    `json:"maxlines"`
	MaxSize    int    `json:"maxsize"`
	Daily      bool   `json:"daily"`
	MaxDays    int64  `json:"maxdays"`
	Level      string `json:"level"` // 日志打印登记，选项：EMER，ALRT，CRIT，EROR，WARN，INFO，DEBG，TRAC
	PermitMask string `json:"permit"`

	MaxSizeCurSize   int
	MaxLinesCurLines int
	DailyOpenDate    int
	DailyOpenTime    time.Time
	FileNameOnly     string
	Suffix           string
}

type ConsoleConfig struct {
	Level    string `json:"level"` // 日志打印登记，选项：EMER，ALRT，CRIT，EROR，WARN，INFO，DEBG，TRAC
	Colorful bool   `json:"color"`
}

type ConnConfig struct {
	ReconnectOnMsg bool   `json:"reconnectOnMsg"`
	Reconnect      bool   `json:"reconnect"`
	Net            string `json:"net"`
	Addr           string `json:"addr"`
	Level          string `json:"level"` // 日志打印登记，选项：EMER，ALRT，CRIT，EROR，WARN，INFO，DEBG，TRAC
	illNetFlag     bool   //网络异常标记
}

type ZapConfig struct {
	Level    string `json:"level"`
	Colorful bool   `json:"color"`
	LogLevel int    `json:"log_level"`
	//Writer io.Writer
	// std(默认) or 具体的文件路径
	Output string `json:"output"`
	// 调用栈往上走的层数
	CallerSkip int `json:"-"` // zap 相较于其他日志库，要多2层，注意
}
