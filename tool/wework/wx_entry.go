package wework

import (
	"fmt"
	"github.com/senyu-up/toolbox/tool/logger"
	"strings"
	"sync"
	"time"
)

const (
	MsgText  = "text"
	MsgVoice = "voice"
	MsgVideo = "video"
	MsgFile  = "file"
	MsgCard  = "textcard"
	MsgNews  = "news"
	MsgMd    = "markdown"
)
const (
	//企业微信登录类型 1:网页链接授权 2:二维码登录
	WebLogin = iota + 1 //授权链接登录
	QrLogin             //二维码登录
)
const (
	//用户信息激活状态: 1=已激活，2=已禁用，4=未激活，5=退出企业。
	UserStatusActivate   = iota + 1 //已激活
	UserStatusBind                  //已禁用
	UserStatusNoActivate            //未激活
	UserStatusOut                   //退出企业
)

var (
	ErrAccessTokenExpired = fmt.Errorf("ToolBox WeWork access token expired")
	ErrRetryFuncNotSet    = fmt.Errorf("ToolBox WeWork retry func not set")
)

type AccessTokenStorageDriver interface {
	// GetAccessToken 获取 access token
	GetAccessToken() (string, error)
	// SetAccessToken 设置 access token
	SetAccessToken(token string, expire time.Duration) error
}

type DefaultAccessTokenStorageDriver struct {
	accessToken string
	expire      time.Duration
	lastTime    time.Time
}

func (sd *DefaultAccessTokenStorageDriver) GetAccessToken() (string, error) {
	if sd.accessToken == "" || time.Now().Sub(sd.lastTime) > sd.expire {
		return "", ErrAccessTokenExpired
	}
	return sd.accessToken, nil
}

func (sd *DefaultAccessTokenStorageDriver) SetAccessToken(token string, expire time.Duration) error {
	sd.accessToken = token
	sd.expire = expire
	sd.lastTime = time.Now()
	return nil
}

// 判断重试，或者返回错误
type RetryOrReturnFunc func(respCommonParam, error) (bool, error)

// 公共响应参数
type respCommonParam struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// 获取access_token
type respAccessToken struct {
	respCommonParam
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int64  `json:"expires_in,omitempty"`
}

// 获取部门列表
type respDepartmentList struct {
	respCommonParam
	GetDepartmentListOutput []*Department `json:"department,omitempty"`
}

// 获取部门信息
type respDepartmentInfo struct {
	respCommonParam
	GetDepartmentSimpleOutPut []*DepartmentSimpleUserInfo `json:"userlist,omitempty"`
}

// 获取部门信息-详情
type respDepartmentUserInfo struct {
	respCommonParam
	GetDepartmentUserListOutPut []*DepartmentUserInfo `json:"userlist,omitempty"`
}

// 创建会话
type respCreateChat struct {
	respCommonParam
	ChatID string `json:"chatid"`
}

// 获取用户信息
type respUserInfo struct {
	respCommonParam
	UserInfo
}

// token信息存储
type access struct {
	sync.Mutex
	Expiration  int64
	AccessToken string
}

// 创建群聊请求
type CreateGroupChatInput struct {
	Name     string   `json:"name"`
	Owner    string   `json:"owner"`
	Userlist []string `json:"userlist"`
	ChatID   string   `json:"chatid,omitempty"`
}

// 修改群信息
type RemoveGroupChatUserInput struct {
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	ChatID      string   `json:"chatid,omitempty"`
	AddUserList []string `json:"add_user_list"`
	DelUserList []string `json:"del_user_list"`
}

// Department 获取部门列表
type Department struct {
	ID       int32  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name     string `json:"name"`
	NameEn   string `json:"name_en"`
	Parentid int    `json:"parentid"`
	Order    int    `json:"order"`
}
type DepartmentSimpleUserInfo struct {
	Userid     string `json:"userid"`
	Name       string `json:"name"`
	Department []int  `json:"department"`
	OpenUserid string `json:"open_userid"`
}

type DepartmentUserInfo struct {
	Userid         string   `json:"userid"`
	Name           string   `json:"name"`
	Alias          string   `json:"alias"`
	Department     []int    `json:"department"`
	Order          []int    `json:"order"`
	OpenUserid     string   `json:"open_userid"`
	Position       string   `json:"position"`
	Mobile         string   `json:"mobile"`
	Email          string   `json:"email"`
	IsLeaderInDept []int32  `json:"is_leader_in_dept"`
	DirectLeader   []string `json:"direct_leader"`
	Avatar         string   `json:"avatar"`
	ThumbAvatar    string   `json:"thumb_avatar"`
	QrCode         string   `json:"qr_code"`
	Status         int      `json:"status"`
	Gender         string   `json:"gender"`
	MainDepartment int      `json:"main_department"`
}

type MsgBodyInput struct {
	ChatID   string    `json:"chatid"`
	MsgType  string    `json:"msgtype"`
	Text     *Text     `json:"text,omitempty"`
	MarkDown *Text     `json:"markdown,omitempty"`
	Image    *Image    `json:"image,omitempty"`
	Voice    *Voice    `json:"voice,omitempty"`
	Video    *Video    `json:"video,omitempty"`
	File     *File     `json:"file,omitempty"`
	TextCard *TextCard `json:"textcard,omitempty"`
	News     *News     `json:"news,omitempty"`
	Safe     int       `json:"safe"`
}

type ToMsg struct {
	MsgType  string    `json:"msgtype"`
	Text     *Text     `json:"text,omitempty"`
	MarkDown *Text     `json:"markdown,omitempty"`
	Image    *Image    `json:"image,omitempty"`
	Voice    *Voice    `json:"voice,omitempty"`
	Video    *Video    `json:"video,omitempty"`
	File     *File     `json:"file,omitempty"`
	TextCard *TextCard `json:"textcard,omitempty"`
	News     *News     `json:"news,omitempty"`
	Safe     int       `json:"safe"`
}

type ToMsgBodyInput struct {
	//三个to必须要有一个有值
	ToUser  string `json:"touser,omitempty"`
	ToParty string `json:"toparty,omitempty"`
	ToTag   string `json:"totag,omitempty"`

	AgentId  string    `json:"agentid"`
	MsgType  string    `json:"msgtype"`
	Text     *Text     `json:"text,omitempty"`
	MarkDown *Text     `json:"markdown,omitempty"`
	Image    *Image    `json:"image,omitempty"`
	Voice    *Voice    `json:"voice,omitempty"`
	Video    *Video    `json:"video,omitempty"`
	File     *File     `json:"file,omitempty"`
	TextCard *TextCard `json:"textcard,omitempty"`
	News     *News     `json:"news,omitempty"`
	Safe     int       `json:"safe"`
}

type ToIdType string

const (
	ToIdTypeUser  ToIdType = "user"
	ToIdTypeParty ToIdType = "party"
	ToIdTypeTag   ToIdType = "tag"
)

type ToId struct {
	ids []string
	tp  ToIdType
}

func (t *ToId) SetType(in ToIdType) {
	t.tp = in
}

func (t *ToId) Add(in ...string) {
	t.ids = append(t.ids, in...)
}

func (t *ToId) String() string {
	return strings.Join(t.ids, "|")
}

type Text struct {
	Content string `json:"content"`
}

type Image struct {
	MediaID string `json:"media_id"`
}
type Voice struct {
	MediaID string `json:"media_id"`
}
type Video struct {
	MediaID     string `json:"media_id"`
	Description string `json:"description"`
	Title       string `json:"title"`
}
type File struct {
	MediaID string `json:"media_id"`
}
type TextCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}
type Articles struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PicUrl      string `json:"picurl"`
}
type News struct {
	Articles []Articles `json:"articles"`
}

type UserInfo struct {
	Errcode        int    `json:"errcode"`
	Errmsg         string `json:"errmsg"`
	Userid         string `json:"userid"`
	Name           string `json:"name"`
	Department     []int  `json:"department"`
	Order          []int  `json:"order"`
	Position       string `json:"position"`
	Mobile         string `json:"mobile"`
	Gender         string `json:"gender"`
	Email          string `json:"email"`
	IsLeaderInDept []int  `json:"is_leader_in_dept"`
	Avatar         string `json:"avatar"`
	ThumbAvatar    string `json:"thumb_avatar"`
	Telephone      string `json:"telephone"`
	Alias          string `json:"alias"`
	Address        string `json:"address"`
	OpenUserid     string `json:"open_userid"`
	MainDepartment int    `json:"main_department"`
	Extattr        struct {
		Attrs []struct {
			Type int    `json:"type"`
			Name string `json:"name"`
			Text struct {
				Value string `json:"value"`
			} `json:"text,omitempty"`
			Web struct {
				URL   string `json:"url"`
				Title string `json:"title"`
			} `json:"web,omitempty"`
		} `json:"attrs"`
	} `json:"extattr"`
	Status           int    `json:"status"`
	QrCode           string `json:"qr_code"`
	ExternalPosition string `json:"external_position"`
	ExternalProfile  struct {
		ExternalCorpName string `json:"external_corp_name"`
		WechatChannels   struct {
			Nickname string `json:"nickname"`
			Status   int    `json:"status"`
		} `json:"wechat_channels"`
		ExternalAttr []struct {
			Type int    `json:"type"`
			Name string `json:"name"`
			Text struct {
				Value string `json:"value"`
			} `json:"text,omitempty"`
			Web struct {
				URL   string `json:"url"`
				Title string `json:"title"`
			} `json:"web,omitempty"`
			Miniprogram struct {
				Appid    string `json:"appid"`
				Pagepath string `json:"pagepath"`
				Title    string `json:"title"`
			} `json:"miniprogram,omitempty"`
		} `json:"external_attr"`
	} `json:"external_profile"`
}

// 获取访问用户身份
type RespUserId struct {
	respCommonParam
	UserId   string `json:"UserId,omitempty"`
	DeviceId string `json:"DeviceId,omitempty"`
}

/*
invaliduser	不合法的userid，不区分大小写，统一转为小写
invalidparty	不合法的partyid
invalidtag	不合法的标签id
msgid	消息id，用于撤回应用消息
response_code	仅消息类型为“按钮交互型”，“投票选择型”和“多项选择型”的模板卡片消息返回，应用可使用response_code调用更新模版卡片消息接口，24小时内有效，且只能使用一次
*/
type RespSendMsg struct {
	respCommonParam
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
	MsgId        string `json:"msgid"`
	ResponseCode string `json:"response_code"`
}

// https://developer.work.weixin.qq.com/document/path/90313
func DefaultRetryOrReturnFunc(val respCommonParam, inErr error) (bool, error) {
	if val.ErrCode == 0 {
		return false, inErr
	}
	if val.ErrCode == 42001 || // 过期
		val.ErrCode == 40014 || // 不合法
		val.ErrCode == 41001 || // token missing
		val.ErrCode == 40082 || // 不合法
		val.ErrCode == 42009 { // 过期
		logger.Warn("token需要重新获取, param: %v", val)
		return true, nil
	}
	return false, inErr
}
