package qwrobot

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/panjf2000/ants/v2"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/http/req"
)

const DefaultFreqLimit = 100

const (
	UnitHour   = "H"
	UnitMinute = "M"
	UnitSecond = "S"
)

const DefaultFreqUnit = UnitHour

const (
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

var luaScript = redis.NewScript(`
local i = redis.call("INCR", KEYS[1])
if i <= 1 then
	redis.call("EXPIRE", KEYS[1], ARGV[1])
end
return i 
`)

type Message struct {
	// 标题
	Title string
	// 内容
	Content string
	// 通知的成员列表, 可不填。 填了谁，发出的消息就会@谁， 传 []string{"@all"} 表示@所有人
	UserList []string
	level    string
}

// New
//
//	@Description: 初始化企业微信机器人， api doc：https://developer.work.weixin.qq.com/tutorial/detail/54
//	@param conf  body any true "-"
//	@return *QWRobot
func New(conf *config.QwRobotConfig) *QWRobot {
	robot := &QWRobot{
		config:   *conf,
		redisCli: conf.RedisCli,
	}
	robot.infoFreqLimit, robot.infoFreqTTL = robot.parseFreqLimit(conf.InfoFreqLimit)
	robot.warnFreqLimit, robot.warnFreqTTL = robot.parseFreqLimit(conf.WarnFreqLimit)
	robot.errFreqLimit, robot.errFreqTTL = robot.parseFreqLimit(conf.ErrorFreqLimit)
	// 清理可能已有的脏数据
	if robot.redisCli != nil {
		robot.redisCli.Del(robot.getCacheKey(LevelInfo))
		robot.redisCli.Del(robot.getCacheKey(LevelWarn))
		robot.redisCli.Del(robot.getCacheKey(LevelError))
	}

	if conf.RedisCli != nil && conf.Webhook != "" {
		robot.ready = true
	}

	opt := ants.Options{
		ExpiryDuration:   time.Hour,
		PreAlloc:         false,
		MaxBlockingTasks: 10000,
		Nonblocking:      false,
		PanicHandler:     nil,
		Logger:           nil,
	}
	options := ants.WithOptions(opt)
	pool, _ := ants.NewPoolWithFunc(100, robot.send, options)
	robot.pool = pool

	return robot
}

type QWRobot struct {
	config        config.QwRobotConfig
	redisCli      redis.UniversalClient
	infoFreqLimit int
	infoFreqTTL   int

	warnFreqLimit int
	warnFreqTTL   int

	errFreqLimit int
	errFreqTTL   int

	hostName string
	stage    string
	ip       string

	msgCh chan Message
	ready bool
	pool  *ants.PoolWithFunc
}

func (q *QWRobot) getCacheKey(level string) string {
	return q.config.Prefix + "_" + q.config.MessageType + "_" + level
}

func (q *QWRobot) Incr(key string, ttl int) int {
	cnt, err := luaScript.Run(q.redisCli, []string{key}, ttl).Int()
	if err != nil {
		fmt.Printf("redis lua error key:%s ttl:%d err:%s\n", key, ttl, err.Error())
		return 0
	}

	return cnt
}

type qwMarkdownMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Content *string `json:"content"`
	} `json:"markdown"`
}

type qwTextMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content       *string  `json:"content"`
		MentionedList []string `json:"mentioned_list"`
	} `json:"text"`
}

func (q *QWRobot) getMessage(msg Message, level string) []byte {
	// 主动拼接环境信息
	hostName := q.hostName
	if hostName == "" {
		hostName = "-"
	}
	ip := q.ip
	if ip == "" {
		ip = "-"
	}

	msg.Content += fmt.Sprintf("\n[host: %s] [ip: %s] [stage: %s]", hostName, ip, q.stage)
	var qwMsg interface{}
	var content = new(string)
	if msg.UserList != nil && len(msg.UserList) > 0 {
		*content = "[" + q.config.MessageType + "] " + msg.Title + ", " + msg.Content
		qwMsg = qwTextMessage{
			MsgType: "text",
			Text: struct {
				Content       *string  `json:"content"`
				MentionedList []string `json:"mentioned_list"`
			}{Content: content, MentionedList: msg.UserList},
		}
	} else {
		var color string
		if level == LevelInfo {
			color = "blue"
		} else if level == LevelWarn {
			color = "yellow"
		} else {
			color = "red"
		}

		*content = `<font color="` + color + `">[` + q.config.MessageType + "] " + msg.Title + "</font>\n>" + msg.Content

		qwMsg = qwMarkdownMessage{
			MsgType: "markdown",
			Markdown: struct {
				Content *string `json:"content"`
			}{
				Content: content,
			},
		}
	}
	if len(*content) > 2048 {
		*content = (*content)[0:2048]
	}

	data, _ := jsoniter.Marshal(qwMsg)

	return data
}

func (q *QWRobot) parseFreqLimit(freq string) (times int, ttl int) {
	items := strings.SplitN(freq, "/", 2)
	var limit = DefaultFreqLimit
	var unit = DefaultFreqUnit

	if len(items) == 2 {
		limit, _ = strconv.Atoi(items[0])
		if limit < 1 {
			limit = DefaultFreqLimit
		}
		unit = strings.ToUpper(items[1])
		if unit != UnitMinute && unit != UnitSecond && unit != UnitHour {
			unit = DefaultFreqUnit
		}
	}

	if unit == UnitMinute {
		ttl = 60
	} else if unit == UnitSecond {
		ttl = 1
	} else {
		ttl = 3600
	}

	return limit, ttl
}

func (q *QWRobot) checkLimit(key string, ttl, limit int) bool {
	curN := q.Incr(key, ttl)
	if curN <= limit {
		return true
	}

	return false
}

func (q *QWRobot) getLimitAndTTL(level string) (int, int) {
	if level == LevelInfo {
		return q.infoFreqLimit, q.infoFreqTTL
	} else if level == LevelWarn {
		return q.warnFreqLimit, q.warnFreqTTL
	} else {
		return q.errFreqLimit, q.errFreqTTL
	}
}

func (q *QWRobot) send(in interface{}) {
	msg := in.(Message)
	if q.config.Webhook == "" {
		fmt.Println("qw robot webhook undefined ignored!!!")
		return
	}

	if q.config.RedisCli == nil {
		fmt.Println("qw robot redis cli undefined ignored!!!")
		return
	}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("qw robot send error", err)
		}
	}()
	// 1. 检查频率
	key := q.getCacheKey(msg.level)
	limit, ttl := q.getLimitAndTTL(msg.level)
	content := q.getMessage(msg, msg.level)

	if !q.checkLimit(key, ttl, limit) {
		fmt.Println("frequency limit exceeded level:"+msg.level+" type:"+q.config.MessageType, "content:"+string(content))
		return
	}

	err := q.sendToQW(q.config.Webhook, content)
	if err != nil {
		fmt.Println("webhook=" + q.config.Webhook + " content:" + string(content) + " err:" + err.Error())
	}
}

func (q *QWRobot) sendToQW(webhook string, data []byte) error {
	response, err := req.New(context.Background()).Header(map[string]string{"Content-Type": "application/json"}).Body(data).Post(webhook)
	if err != nil {
		return err
	}
	if response != nil && response.IsSuccess() {
		return nil
	} else {
		if response != nil {
			return fmt.Errorf("code:%d", response.StatusCode)
		}
		return errors.New("request error")
	}
}

func (q *QWRobot) Info(msg Message) {
	if q.ready {
		msg.level = LevelInfo
		_ = q.pool.Invoke(msg)
	}
}

func (q *QWRobot) Warn(msg Message) {
	if q.ready {
		msg.level = LevelWarn
		_ = q.pool.Invoke(msg)
	}
}

func (q *QWRobot) Error(msg Message) {
	if q.ready {
		msg.level = LevelError
		_ = q.pool.Invoke(msg)
	}
}

var qwRobotInst *QWRobot
var qwRobotOnce sync.Once

func Get() *QWRobot {
	return qwRobotInst
}

// Init
//
//	@Description: 初始化 全局 qw 机器人
//	@param conf  body any true "-"
//	@param opts  body any true "-"
//	@return qwr
//	@return err
func Init(conf *config.QwRobotConfig, opts ...QwrOption) (qwr *QWRobot, err error) {
	if qwRobotInst == nil {
		qwRobotOnce.Do(func() {
			msgType := conf.MessageType
			if msgType == "" {
				err = errors.New("QWRobot config.MessageType not set ")
			}
			qwRobotInst = New(&config.QwRobotConfig{
				Webhook:        conf.Webhook,
				InfoFreqLimit:  conf.InfoFreqLimit,
				WarnFreqLimit:  conf.WarnFreqLimit,
				ErrorFreqLimit: conf.ErrorFreqLimit,
				MessageType:    msgType,
				Prefix:         conf.Prefix,
				RedisCli:       conf.RedisCli,
			})

			// 应用 options
			for _, opt := range opts {
				opt(qwRobotInst)
			}
		})
	}

	return qwRobotInst, err
}
