package wework

import (
	"net/http"
)

type WxOption func(*WechatClient)

func OptWithHttpClient(hc *http.Client) WxOption {
	return func(option *WechatClient) {
		if nil != hc {
			option.hclient = hc
		}
	}
}

func OptWithAccess(ac *access) WxOption {
	return func(option *WechatClient) {
		if nil == ac {
			option.access = ac
		}
	}
}

// ac token storage driver
func OptWithAccessTokenStoreDriver(ac AccessTokenStorageDriver) WxOption {
	return func(option *WechatClient) {
		if nil != ac {
			option.acStoreDriver = ac
		}
	}
}
