package wework

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// GetUserInfoRequest
//
//	@Description: 读取成员相关信息, api doc: https://developer.work.weixin.qq.com/document/path/90196
//	@receiver this
//	@param uid  body any true "-"
//	@return info
//	@return err
func (this *WechatClient) GetUserInfoRequest(uid string) (info *UserInfo, err error) {
	genUrl := func() string {
		var ac, _ = this.acStoreDriver.GetAccessToken()
		return fmt.Sprintf(appGetUserInfo, ac, uid, this.debug)
	}
	var resp respUserInfo
	err = this.autoRetryRequest(genUrl, http.MethodGet, nil, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.UserInfo, err
}

// GetUserIdByCode
//
//	@Description:  通过登陆code获取访客身份，api doc：https://developer.work.weixin.qq.com/document/path/91707 这个 /user/getuserinfo 是老接口，推荐用auth/getuserinfo，
//	@receiver this
//	@param code  body any true "-"
//	@return userId
//	@return err
func (this *WechatClient) GetUserIdByCode(code string) (userId string, err error) {
	genUrl := func() string {
		var ac, _ = this.acStoreDriver.GetAccessToken()
		return fmt.Sprintf(appchatGetUserId, ac, code, this.debug)
	}
	var resp RespUserId
	err = this.autoRetryRequest(genUrl, http.MethodGet, nil, &resp)
	if err != nil {
		return userId, err
	}
	if resp.ErrCode != 0 {
		return userId, errors.New(resp.ErrMsg)
	}
	userId = resp.UserId
	return userId, err
}

// 获取用户授权登录链接地址 _type 1:网页授权 2:扫码授权 redirectUri 授权成功后重定向地址
func (this *WechatClient) GetAuthUrl(_type int, redirectUri string) (authUrl string) {
	if _type == WebLogin {
		return fmt.Sprintf(appchatWebAuthorize, this.corpID, url.QueryEscape(redirectUri))
	}
	if _type == QrLogin {
		return fmt.Sprintf(appchatQrAuthorize, this.corpID, this.agentID, url.QueryEscape(redirectUri))
	}
	return authUrl
}
