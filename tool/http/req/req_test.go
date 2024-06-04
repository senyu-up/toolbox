package req

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

// TestGet
// @description GET 调用
func TestGet(t *testing.T) {
	resp, err := New(context.Background()).Get("http://www.baidu.com")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(resp.String())
}

func TestDebug(t *testing.T) {
	//  开启debug
	_, err := New(context.Background()).Debug().Get("http://www.baidu.com")
	if err != nil {
		t.Error(err)
		return
	}
}

type RespData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Items []interface{} `json:"items"`
	} `json:"data"`
	RequestId string `json:"request_id"`
}

// TestPost
// @description POST 调用 并 赋值
func TestPost(t *testing.T) {
	// body 参数
	data := map[string]interface{}{
		"app_key": []string{},
	}
	// cookie
	cookie := []*http.Cookie{
		{
			Name:  "grafana_session",
			Value: "d4a66588fb24b60dc71ecda91ed0838c",
		},
		{
			Name:  "vue_admin_template_token",
			Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIxNjU0NzYwODczIiwiaXRhIjoiMTY1NDE1NjA3MyIsInBhc3N3b3JkIjoiIiwicm9sZUlkIjoiMSIsInV1aWQiOiJjYWM2bWFjc2FkNWx2NnBxN3YxZyJ9.3aF97M6ROM8uc_LbVO1JiqyyhRtovB0T3U_v6KOA28w",
		},
	}
	// 链式调用发起一个请求
	resp, err := New(context.Background()).Debug().Header(map[string]string{
		"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjoie1wiaWRcIjoxLFwid3hfdXNlcl9pZFwiOlwiUWluZ0NoZW5nXCIsXCJuYW1lXCI6XCLpnZLln45cIixcInBvc2l0aW9uXCI6XCLlkI7nq6_nqIvluo9cIixcIm1vYmlsZVwiOlwiXCIsXCJhdmF0YXJcIjpcImh0dHA6Ly93ZXdvcmsucXBpYy5jbi9iaXptYWlsL2tVbGljT1R2dzNRdU5mQmFWQjZQWTBndUt1dEpJdjlkM1d1MzVVNFJoSnlJejNYdk1IaWFuTEZBLzBcIixcInRodW1iX2F2YXRhclwiOlwiaHR0cDovL3dld29yay5xcGljLmNuL2Jpem1haWwva1VsaWNPVHZ3M1F1TmZCYVZCNlBZMGd1S3V0Skl2OWQzV3UzNVU0UmhKeUl6M1h2TUhpYW5MRkEvMTAwXCIsXCJhbGlhc1wiOlwiXCIsXCJvcGVuX3VzZXJpZFwiOlwiXCIsXCJnZW5kZXJcIjoxLFwibWFpbl9kZXBhcnRtZW50XCI6NDM1LFwicHJvamVjdF90ZWFtXCI6XCJTMDHlhazlhbHmlK_mjIFcIixcImNyZWF0ZWRfYXRcIjoxNjQxNDU5MDYzLFwidXBkYXRlZF9hdFwiOjE2NDE0NTkwNjMsXCJzdXBlclwiOnRydWUsXCJzZXRfcGVybWlzc2lvblwiOjEsXCJzdGF0dXNcIjoxLFwiZGVwYXJ0bWVudF9uYW1lXCI6XCJcIixcImRlcGFydG1lbnRfdXNlcnNcIjpudWxsLFwiYWRtaW5fcm9sZV9hY3Rpb25fcGVybWlzc2lvbnNcIjpudWxsLFwiYWRtaW5fcm9sZV9zZWxlY3RfcGVybWlzc2lvbnNcIjpudWxsfSIsImV4cCI6MTY3MjcxNzI2NiwibmJmIjoxNjQ2Nzk3MjY2LCJpYXQiOjE2NDY3OTcyNjZ9.lwJi65EpmeDO_dmfZYiHsa7qyvWcGogapz2_8NjnIR8",
	}).Cookie(cookie).BodyJson(data).Post("http://platform-develop.outer.xxxxxx.com:32142/log_track/overview/number")
	if err != nil {
		t.Error(err)
		return
	}
	resp.String()
	//  判断请求是否成功
	if resp.IsSuccess() {
		// 解析响应数据
		trackData := &RespData{}
		err = resp.Unmarshal(trackData)
		if err != nil {
			return
		}
		t.Logf("%+v\n", trackData)
	} else {
		t.Error("请求失败")
	}
}
