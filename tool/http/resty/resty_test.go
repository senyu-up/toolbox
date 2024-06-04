package resty

import "testing"

func TestHttpGetResJson(t *testing.T) {
	userInfo := &struct {
		FirstName  string `json:"first_name,omitempty"`
		LastName   string `json:"last_name,omitempty"`
		ProfilePic string `json:"profile_pic,omitempty"`
		Id         string `json:"id,omitempty"`
	}{}
	url := "https://graph.facebook.com/5616624041714007"
	params := map[string]string{
		"access_token": "EAALKuhAVkqwBAKxn1mczb5dADz4dd4ZBgMAmb0dFA0jY0cEhau2VV6XPNhqR0eCLZCxcaWwonj1nkwbVdsQiS30s3joJ8w779mWYqWtzOFJl3WRociium7N6bsBkkqwszV5o70iKx5cTBbtimTW4BzWIR8VOYRy0PY2y69K4cmc9RJYiZAH",
	}
	res, err := HttpGetResJson(url, params, userInfo)
	if err != nil {
		t.Error(err)
	}
	t.Log(res.StatusCode())
	t.Log(userInfo.LastName + userInfo.FirstName)
}

func TestHttpSendJsonResJson(t *testing.T) {
	type SendMsgRes struct {
		RecipientId string `json:"recipient_id,omitempty"`
		MessageId   string `json:"message_id,omitempty"`
	}
	var respone SendMsgRes
	url := "https://graph.facebook.com/v15.0/me/messages?access_token=EAALKuhAVkqwBAPrzLJZBhXdGIcbT50Mbg9KdvBuLZCtN2oOFzQyqOAffin7nX2hZAaODPHGCZCMeZANYTcyD1EdVMvrklLKGPV9bBnixVkbZANzsJGD0tTlZC7epchi4vlCBrdldF3lb84BbHEZAOIDZAdCOJKiI35tmfxiwLhPHZBZBXWu2sqIxvugVJPSLOmLTE0ZD"
	paramsJson := `{"recipient":{"id":"5481046535350151"},"message":{"attachment":{"type":"image","payload":{"url":"https://platform-im-test.allxxx.com/2022/10/24/6df088ace3ba1a3c826c3e175e4da3fd.webp?X-Amz-Algorithm=AWS4-HMAC-SHA256\u0026X-Amz-Credential=AKIA2LZCRPKRKWIHSGFM%2F20221024%2Fus-west-2%2Fs3%2Faws4_request\u0026X-Amz-Date=20221024T103133Z\u0026X-Amz-Expires=604800\u0026X-Amz-SignedHeaders=host\u0026X-Amz-Signature=013c02752568bef6ff7a94d38c2bd9872c28b1af8d30ca246230a1ec00b75eed"}}}}`
	res, err := HttpSendJsonResJson(url, "POST", paramsJson, &respone)
	if err != nil {
		t.Error(err)
	}
	t.Log(res.StatusCode())
	t.Log(respone)
}

func TestHttpSendPost(t *testing.T) {
	paramsJson := `{"recipient": {
        "id": "5481046535350151"
    },
    "message": {
        "attachment": {
            "type": "template",
            "payload": {
                "template_type": "customer_feedback",
                "title": "服务评价",
                "subtitle": "你的评价对我们很重要",
                "button_title": "提交评价",
                "feedback_screens": [
                    {
                        "questions": [
                            {
                                "id": "hauydmns8",
                                "type": "csat",
                                "title": "请提交你对我们的评价",
                                "score_label": "neg_pos",                      "score_option": "five_stars", 
                                "follow_up": {
                                    "type": "free_form",
                                    "placeholder": "Give additional feedback" 
                                }
                            }
                        ]
                    }
                ],
                "business_privacy": {
                    "url": "https://www.example.com"
                },
                "expires_in_days": 3
            }
        }
    }
}`
	url := "https://graph.facebook.com/v15.0/me/messages?access_token=EAALKuhAVkqwBAPrzLJZBhXdGIcbT50Mbg9KdvBuLZCtN2oOFzQyqOAffin7nX2hZAaODPHGCZCMeZANYTcyD1EdVMvrklLKGPV9bBnixVkbZANzsJGD0tTlZC7epchi4vlCBrdldF3lb84BbHEZAOIDZAdCOJKiI35tmfxiwLhPHZBZBXWu2sqIxvugVJPSLOmLTE0ZD"
	var respone interface{}
	res, err := HttpSendJsonResJson(url, "POST", paramsJson, &respone)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
	t.Log(respone)
}
