package snow

import (
	"bytes"
	"errors"
	"github.com/spf13/cast"
	"math/rand"
	"strings"
	"time"
)

const (
	SeedChar = "a,b,c,d,e,f,g,h,j,k,m,n,p,q,r,s,t,u,v,w,x,y,z,A,B,C,D,E,F,G,H,J,K,M,N,P,Q,R,S,T,U,V,W,X,Y,Z,2,3,4,5,6,7,8,9"
)

/*
策略:
头N位固定乱序生成
1位时间位
2位顺序位
当随机种子设置为秒时,该算法理论上一秒内可生成不重复ID个数位 len(SendChar) * len(SendChar)
*/

type MySnow struct {
	LastTime int64
	NowTime  int64
	Len      int
	Mod      int32 //模
	Count    int32 //余
	Number   int32
}

func InitMySnow(len int) *MySnow {
	return &MySnow{Len: len, LastTime: time.Now().Unix(), NowTime: time.Now().Unix()}
}

func (p *MySnow) RandString(seed int64) (string, error) {
	if p.Len <= 5 {
		return "", errors.New("String len to short. ")
	}
	chars := strings.Split(SeedChar, ",")
	SCount := len(chars)
	rand.Seed(seed)
	buffer := &bytes.Buffer{}
	var err error
	for i := 0; i < p.Len-3; i++ {
		c := chars[rand.Intn(SCount)]
		_, err = buffer.WriteString(c)
		if err != nil {
			return "", err
		}
	}
	p.NowTime = time.Now().Unix()
	//时间位
	_, err = buffer.WriteString(chars[p.NowTime%cast.ToInt64(SCount)])
	//取模
	p.Mod = p.Number / cast.ToInt32(SCount)
	if p.Mod >= cast.ToInt32(SCount) {
		p.Mod = int32(SCount) - 1
	}
	//取余
	p.Count = p.Number % cast.ToInt32(SCount)

	//取模位
	_, err = buffer.WriteString(chars[p.Mod])
	//余位
	_, err = buffer.WriteString(chars[p.Count])

	if p.LastTime == p.NowTime {
		p.Number++
	} else {
		p.LastTime = p.NowTime
		p.Number = 0
	}
	if err != nil {
		return "", err
	}
	p.LastTime = time.Now().Unix()
	return buffer.String(), nil
}
