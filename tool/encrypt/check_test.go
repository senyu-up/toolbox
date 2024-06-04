package encrypt

import (
	"testing"
)

const publicKey = `-----BEGIN RSA PUBLIC KEY-----
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBANnNLY2jLg3nCjznsk88p93eoV5bKznE
y2Yhx+Z0ymSMtEH5ywhbeDrUdnatzbr9PlKLN7D9f7bFIO+e5aZeqzkCAwEAAQ==
-----END RSA PUBLIC KEY-----
`

const priKey = `-----BEGIN RSA PRIVATE KEY-----
MIIBOQIBAAJBANnNLY2jLg3nCjznsk88p93eoV5bKznEy2Yhx+Z0ymSMtEH5ywhb
eDrUdnatzbr9PlKLN7D9f7bFIO+e5aZeqzkCAwEAAQJAYJUx7Cs+kv+vdI6ybZzS
O/srx/HZ20Y/hEzanoGP9sH0uGr4nsmn9e5bfjm0yu/soknXHXGYFrLdusteXjWS
uQIhAN+HxHxVmQCMAQ40OgttKzomTOEpU6e2e6Ye5tbr7ujnAiEA+XBgP+/azu2u
qYb9yH/rEnnzPWVkIqlalSOL/GRkpt8CIA7+YGOmqjirK3b0ceBKVlf0Mbv4ta/O
QcUG1Z0c/k2JAiAqTffVAC4BCGimEeH63k8VDB/H2ulXw5c8UhIM1U4IywIgVTnN
LmtehykxYmoTCxgzLgsHwWurROX4mN5RK2RLtjk=
-----END RSA PRIVATE KEY-----
`

func Test_Check(t *testing.T) {
	check, _ := NewCheck(publicKey)
	tag := &Tag{
		Ts:        1649970843,
		Seed:      5,
		Signature: "T6kTgPI91F6B0hohZ9uDQSS/Rp3sskqWklplr9+h1EBxAG07sDdvL+nRiJz15cBdekfRfQS4eD5wKWsNsgoC+w==",
		Data: map[string]interface{}{
			"xxxr": map[string]interface{}{
				"app_key":             "ap2p",
				"cp_account_id":       "3",
				"cp_active_player_id": "2431",
				"sdk_account_id":      "1",
			},
			"head": "head",
			"name": "测试",
			"type": "xxxr",
			"seed": 5,
			"ts":   1649970843,
		},
	}
	if ok, _ := check.Check(tag); !ok {
		t.Error("签名验证算法错误")
	}
}

func Test_ERRCheck(t *testing.T) {
	check, _ := NewCheck(publicKey)
	tag := &Tag{
		Ts:        1649970843,
		Seed:      5,
		Signature: "T6kTgPI91F6B0hohZ9uDQSS/Rp3sskqWklplr9+h1EBxAG07sDdvL+nRiJz15cBdekfRfQS4eD5wKWsNsgoC+w==",
		Data: map[string]interface{}{
			"xxxr": map[string]interface{}{
				"app_key":             "ap2p",
				"cp_account_id":       "3",
				"cp_active_player_id": "2431",
				"sdk_account_id":      "1",
			},
			"head": "head",
			"name": "测试",
			"type": "gam123er",
			"seed": 5,
			"ts":   1649970843,
		},
	}
	if ok, _ := check.Check(tag); ok {
		t.Error("签名验证算法错误")
	}
}

func TestNewSig(t *testing.T) {
	signature := "T6kTgPI91F6B0hohZ9uDQSS/Rp3sskqWklplr9+h1EBxAG07sDdvL+nRiJz15cBdekfRfQS4eD5wKWsNsgoC+w=="
	sig, _ := NewSig(priKey)
	signature2, _ := sig.Signature(&Tag{
		Ts:        1649970843,
		Seed:      5,
		Signature: signature,
		Data: map[string]interface{}{
			"xxxr": map[string]interface{}{
				"app_key":             "ap2p",
				"cp_account_id":       "3",
				"cp_active_player_id": "2431",
				"sdk_account_id":      "1",
			},
			"head": "head",
			"name": "测试",
			"type": "xxxr",
			"seed": 5,
			"ts":   1649970843,
		},
	})
	if signature != signature2 {
		t.Error("签名算法错误")
	}
}
