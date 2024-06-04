package wework

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/logger"
)

type WechatClient struct {
	hclient    *http.Client
	corpID     string `json:"corpid"`
	corpSecret string `json:"corpsecret"`
	agentID    string `json:"agent_id"`  //应用id
	tryTimes   int    `json:"try_times"` // 请求失败后，重试次数
	//deprecated	建议使用 acStoreDriver
	access *access
	debug  int

	retryOrReturn    RetryOrReturnFunc        // 重试或者返回错误判断函数
	acStoreDriver    AccessTokenStorageDriver // access store driver
	refreshTokenLock sync.Mutex               // 更新token 锁, 一个去更新，其他的等着用就行了
	upAccessLock     sync.Mutex               // 更新access token 锁
	refreshInterval  time.Duration            // token 更新间隔
}

const (
	successCode     = 0
	defaultTryTimes = 3
	//获取accessToken
	appGetToken = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	//获取部门列表
	appGetDepartment = "https://qyapi.weixin.qq.com/cgi-bin/department/list?access_token=%s&id=%d&debug=%d"
	//获取部门成员
	appGetDepartmentUserSimpleList = "https://qyapi.weixin.qq.com/cgi-bin/user/simplelist?access_token=%s&department_id=%d&fetch_child=%d&debug=%d"
	//获取部门成员-详情
	appGetDepartmentUserList = "https://qyapi.weixin.qq.com/cgi-bin/user/list?access_token=%s&department_id=%d&fetch_child=%d&debug=%d"
	//获取用户信息
	appGetUserInfo = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s&debug=%d"
	//创建群聊会话
	appchatCreate = "https://qyapi.weixin.qq.com/cgi-bin/appchat/create?access_token=%s&debug=%d"
	//修改群聊会话
	appchatUpdate = "https://qyapi.weixin.qq.com/cgi-bin/appchat/update?access_token=%s&debug=%d"
	//发送会话信息
	appchatSendMsg = "https://qyapi.weixin.qq.com/cgi-bin/appchat/send?access_token=%s&debug=%d"
	//获取访问用户userId
	appchatGetUserId = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s&debug=%d"
	//网页链接授权登录
	appchatWebAuthorize = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_base&state=STATE#wechat_redirect"
	//二维码授权登录
	appchatQrAuthorize = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=%s&agentid=%s&redirect_uri=%s&state=STATE"
	//发送应用消息
	sendMsg = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s&debug=%d"
)

// deprecated
func NewWechatClient(client *http.Client, corpid string, secret, agentId string, debug bool) *WechatClient {
	wclient := &WechatClient{
		hclient:       client,
		corpID:        corpid,
		agentID:       agentId,
		corpSecret:    secret,
		tryTimes:      defaultTryTimes,
		access:        &access{},
		retryOrReturn: DefaultRetryOrReturnFunc,
		acStoreDriver: &DefaultAccessTokenStorageDriver{},
	}
	if debug {
		wclient.debug = 1
	}
	return wclient
}

// InitByConfig
//
//	@Description: 通过配置文件初始化 企微客户端
//	@param conf  body any true "-"
//	@param opts  body any true "-"
//	@return *WechatClient
//	@return error
func InitByConfig(conf *config.WeWorkConfig, opts ...WxOption) (*WechatClient, error) {
	wclient := &WechatClient{
		hclient:         &http.Client{},
		corpID:          conf.CorpId,
		agentID:         conf.AgentId,
		corpSecret:      conf.Secret,
		access:          &access{},
		tryTimes:        defaultTryTimes,
		refreshInterval: 7000 * time.Second,
		retryOrReturn:   DefaultRetryOrReturnFunc,
		acStoreDriver:   &DefaultAccessTokenStorageDriver{},
	}
	if conf.Debug {
		wclient.debug = 1
	}
	if conf.RefreshInterval > 0 {
		wclient.refreshInterval = time.Duration(conf.RefreshInterval) * time.Second
	}
	if conf.TryTimes > 0 {
		wclient.tryTimes = int(conf.TryTimes)
	}
	// 应用 option
	for _, opt := range opts {
		opt(wclient)
	}

	return wclient, nil
	// 下面的不用了, 调用接口时再获取 access token， 不用定时获取了
	err := wclient.getSetToken()
	if err != nil {
		return wclient, err
	}
	//创建一个时间周期去更新access_token
	go func() {
		tk := time.NewTicker(wclient.refreshInterval)
		for {
			select {
			case <-tk.C:
				if err := wclient.getSetToken(); err != nil {
					log.Println(err)
				}
			}
		}
	}()
	return wclient, err
}

// CreateGroupChatRequest 创建群聊会话
func (this *WechatClient) CreateGroupChatRequest(input *CreateGroupChatInput) (chatid string, err error) {
	genUrl := func() string {
		var ac, _ = this.acStoreDriver.GetAccessToken()
		return fmt.Sprintf(appchatCreate, ac, this.debug)
	}
	var resp respCreateChat
	err = this.autoRetryRequest(genUrl, http.MethodPost, input, &resp)
	if err != nil {
		return
	}
	if resp.ErrCode != successCode {
		err = errors.New(resp.ErrMsg)
	}
	chatid = resp.ChatID
	return
}

// RemoveGroupChatUserRequest 修改群成员信息
func (this *WechatClient) RemoveGroupChatUserRequest(input *RemoveGroupChatUserInput) (err error) {
	genUrl := func() string {
		var ac, _ = this.acStoreDriver.GetAccessToken()
		return fmt.Sprintf(appchatUpdate, ac, this.debug)
	}
	var resp respCommonParam
	err = this.autoRetryRequest(genUrl, http.MethodPost, input, &resp)
	if err != nil {
		return
	}
	if resp.ErrCode != successCode {
		err = errors.New(resp.ErrMsg)
	}
	return
}

// 消息发送
func (this *WechatClient) ChatSendMsg(msg *MsgBodyInput) (err error) {

	switch msg.MsgType {
	case MsgMd:
		if msg.MarkDown == nil {
			return errors.New("text msg is nil")
		}
	case MsgText:
		if msg.Text == nil {
			return errors.New("text msg is nil")
		}
	case MsgVideo:
		if msg.Video == nil {
			return errors.New("video msg is nil")
		}
	case MsgFile:
		if msg.File == nil {
			return errors.New("file msg is nil")
		}
	case MsgCard:
		if msg.TextCard == nil {
			return errors.New("card msg is nil")
		}
	case MsgVoice:
		if msg.Voice == nil {
			return errors.New("voice msg is nil")
		}
	case MsgNews:
		if msg.Voice == nil {
			return errors.New("news msg is nil")
		}
	default:
		return errors.New("not found " + msg.MsgType)
	}
	genUrl := func() string {
		var ac, _ = this.acStoreDriver.GetAccessToken()
		return fmt.Sprintf(appchatSendMsg, ac, this.debug)
	}
	var resp respCommonParam
	err = this.autoRetryRequest(genUrl, http.MethodPost, msg, &resp)
	if err != nil {
		return
	}
	if resp.ErrCode != successCode {
		err = errors.New(resp.ErrMsg)
	}
	return
}

// GetGroupChatRequest 获取 access token 强制，不判断过期与否
func (this *WechatClient) refreshSetTokenForce() error {
	logger.Debug("wework force get token")
	var resp respAccessToken
	if this.refreshTokenLock.TryLock() {
		defer this.refreshTokenLock.Unlock()
		logger.Debug("wework force get token, get lock")
		// 获取到了锁，则执行更新逻辑
		urlAddr := fmt.Sprintf(appGetToken, this.corpID, this.corpSecret)
		_, err := this.doRequest(urlAddr, http.MethodGet, nil, &resp)
		if err != nil {
			return err
		}
		if resp.ErrCode != 0 {
			return errors.New(resp.ErrMsg)
		}
		// ac store
		this.acStoreDriver.SetAccessToken(resp.AccessToken, time.Second*time.Duration(resp.ExpiresIn))
		this.access.AccessToken = resp.AccessToken
		this.access.Expiration = time.Now().Unix() + resp.ExpiresIn
	} else {
		// 没有获取到, 则说明其他服务在更新，那么就等待
		logger.Debug("wework wait for token refresh")
		// 但也不能直接跳出, 等待其他协程更新完毕
		this.refreshTokenLock.Lock()
		this.refreshTokenLock.Unlock()
	}
	return nil
}

// deprecated
func (this *WechatClient) getSetToken() error {
	now := time.Now().Unix()
	if this.access.Expiration > now {
		return nil
	}
	this.access.Expiration = now + 7200
	//this.access.Lock()
	//defer this.access.Unlock()
	logger.Debug("get tokend")
	urlAddr := fmt.Sprintf(appGetToken, this.corpID, this.corpSecret)
	var resp respAccessToken
	_, err := this.doRequest(urlAddr, http.MethodGet, nil, &resp)
	if err != nil {
		return err
	}
	if resp.ErrCode != 0 {
		return errors.New(resp.ErrMsg)
	}
	this.access.AccessToken = resp.AccessToken
	//this.access.Expiration = time.Now().Add(2 * time.Hour).Unix()
	this.access.Expiration = now + resp.ExpiresIn
	return nil
}

func (this *WechatClient) autoRetryRequest(genUri func() string, method string, body interface{}, reader interface{}) (err error) {
	var i = 0
	defer func() {
		if i >= this.tryTimes {
			logger.Notify().Error("WechatClient autoRetryRequest, retry times: %d", i)
		}
	}()
	for i = 0; i < this.tryTimes; i++ {
		comResp, err := this.doRequest(genUri(), method, body, reader)
		if this.retryOrReturn != nil {
			if retry, tryErr := this.retryOrReturn(comResp, err); err != nil {
				// 如果 retry 判断返回 err则说明不需要重试，直接返回 err
				return tryErr
			} else if retry {
				// err 为nil 且判断要重试，重试前更新 access token
				if err = this.refreshSetTokenForce(); err != nil {
					return err
				}
			} else {
				// 无报错，无重试，直接返回
				return nil
			}
		} else {
			return ErrRetryFuncNotSet
		}
	}
	return
}

func (this *WechatClient) doRequest(uri string, method string, body interface{}, reader interface{}) (comResp respCommonParam, err error) {
	logger.Debug("WechatClient doRequest:", uri, "-", method, "-", body, "-", reader)
	var request *http.Request
	if method == http.MethodGet {
		request, err = http.NewRequest(method, uri, nil)
		if err != nil {
			return comResp, err
		}
	} else {
		buf := new(bytes.Buffer)
		b, err := jsoniter.Marshal(body)
		if err != nil {
			return comResp, err
		}
		buf.Write(b)
		request, err = http.NewRequest(method, uri, buf)
		if err != nil {
			return comResp, err
		}
	}
	resp, err := this.hclient.Do(request)
	if err != nil {
		return comResp, err
	}
	if resp == nil {
		return comResp, errors.New("resp is nil")
	}
	defer resp.Body.Close()
	//respBody := make([]byte, resp.ContentLength)
	//_, err = resp.Body.Read(respBody)
	respBody, err := ioutil.ReadAll(resp.Body)
	logger.Debug("WechatClient doRequest rsp:", string(respBody))
	if err != nil {
		return comResp, err
	}

	err = jsoniter.Unmarshal(respBody, reader)
	if err != nil {
		return comResp, err
	}
	comResp = respCommonParam{}
	return comResp, jsoniter.Unmarshal(respBody, &comResp)
}

// SendMsg 发送应用消息 https://developer.work.weixin.qq.com/document/path/90236
func (this *WechatClient) SendMsg(msg *ToMsg, id *ToId) (err error) {

	switch msg.MsgType {
	case MsgMd:
		if msg.MarkDown == nil {
			return errors.New("text msg is nil")
		}
	case MsgText:
		if msg.Text == nil {
			return errors.New("text msg is nil")
		}
	case MsgVideo:
		if msg.Video == nil {
			return errors.New("video msg is nil")
		}
	case MsgFile:
		if msg.File == nil {
			return errors.New("file msg is nil")
		}
	case MsgCard:
		if msg.TextCard == nil {
			return errors.New("card msg is nil")
		}
	case MsgVoice:
		if msg.Voice == nil {
			return errors.New("voice msg is nil")
		}
	case MsgNews:
		if msg.Voice == nil {
			return errors.New("news msg is nil")
		}
	default:
		return errors.New("not found " + msg.MsgType)
	}
	toSend := new(ToMsgBodyInput)
	_ = copier.Copy(toSend, msg)
	switch id.tp {
	case ToIdTypeUser:
		toSend.ToUser = id.String()
	case ToIdTypeParty:
		toSend.ToParty = id.String()
	case ToIdTypeTag:
		toSend.ToTag = id.String()
	}
	toSend.AgentId = this.agentID
	genUrl := func() string {
		var ac, _ = this.acStoreDriver.GetAccessToken()
		return fmt.Sprintf(sendMsg, ac, this.debug)
	}
	var resp RespSendMsg
	err = this.autoRetryRequest(genUrl, http.MethodPost, toSend, &resp)
	if err != nil {
		return
	}
	if resp.ErrCode != successCode {
		err = errors.New(fmt.Sprintf("%+v", resp))
	}
	return
}
