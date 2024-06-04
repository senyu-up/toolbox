package main

import (
	"context"
	"fmt"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/email"
	"github.com/senyu-up/toolbox/tool/storage"
	"github.com/spf13/cast"
	"time"
)

func SesDemo() {
	// aws ses 配置
	var conf = &config.EmailConfig{
		Region:    "region",
		AccessID:  "AccessID",
		AccessKey: "AccessKey",
	}
	var ctx = context.Background()

	// 初始化 aws ses 服务
	ses, err := email.InitAwsSes(conf)
	if err != nil {
		fmt.Printf("err:%s", err)
		return
	}

	// 一般的邮件发送
	id, err := ses.Send(ctx, "AgeOfColossus@allxxx.com", []string{"tianyu@xxxxxx.com"}, "Age of Colossus 邮件验证", "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n  <meta charset=\"UTF-8\">\n  <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n  <title>Document</title>\n  <style>\n    body {\n      padding: 20px 30px;\n    }\n  </style>\n</head>\n<body>\n  <div>\n    <h3>校验码：</h3>\n    <p style=\"color: #409eff;font-weight: 800;font-size: 22px;\">693841</p>\n    <div class=\"border-cls\"></div>\n    <p>验证码的有效期为一小时，为不影响您的正常操作，请您及时完成验证。</p>\n    <p>如非本人操作，请忽略此操作。</p>\n  </div>\n</body>\n</html>")
	if err != nil {
		fmt.Printf("err:%s", err)
		return
	} else {
		fmt.Printf("send email return id:%s", id)
	}

	// 使用模版，发送消息
	// 需要填充的内容
	var params = map[string]string{
		"title_user": "tianhai",
		"code":       "123456",
		"expire":     "1小时",
		"note":       "如非本人操作，请忽略此操作",
	}
	// 邮件模版名称
	var tmpName = "email_tmp_12"
	id, err = ses.SendByTemplate(ctx, "AgeOfColossus@allxxx.com", []string{"tianyu@xxxxxx.com"}, tmpName, params)
	if err != nil {
		fmt.Printf("err:%s", err)
		return
	} else {
		fmt.Printf("send email by template sucess return id:%s", id)
	}

}

func GetClient() (*email.AwsPinPoint, error) {
	var conf = &config.EmailConfig{
		Region: storage.AWSUwRegion,
		AppId:  "9643ae5ab829496db7fcec784bd1932d", // pinpoint 创建一个项目，项目的 ID
	}
	ppClient, err := email.InitAwsPinPoint(conf)
	if err != nil {
		fmt.Printf("err:%s", err)
	}
	return ppClient, err
}

func PinPointDemo(ppClient *email.AwsPinPoint) {
	var ctx = context.Background()

	msgId, err := ppClient.Send(ctx,
		"tianhai@xxxxxx.com",
		//"AgeOfColossus@allxxx.com",
		[]string{"773821422@qq.com"},
		"Age of Colossus 邮件验证",
		"<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n  <meta charset=\"UTF-8\">\n  <meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n  <title>Document</title>\n  <style>\n    body {\n      padding: 20px 30px;\n    }\n  </style>\n</head>\n<body>\n  <div>\n    <h3>校验码：</h3>\n    <p style=\"color: #409eff;font-weight: 800;font-size: 22px;\">693841</p>\n    <div class=\"border-cls\"></div>\n    <p>验证码的有效期为一小时，为不影响您的正常操作，请您及时完成验证。</p>\n    <p>如非本人操作，请忽略此操作。</p>\n  </div>\n</body>\n</html>")
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	} else {
		fmt.Printf("send success re msg %s \n", msgId)
	}
}

func ListTemp(ppClient *email.AwsPinPoint) {
	var ctx = context.Background()
	if temps, err := ppClient.ListTemplates(ctx); err != nil {
		fmt.Printf("err %v", err)
	} else {
		fmt.Printf("get temps %+v \n", temps)
	}
}

func ListTemp2(ppClient *email.AwsPinPoint) {
	var ctx = context.Background()
	if temps, err := ppClient.ListTemplates2(ctx); err != nil {
		fmt.Printf("err temp2  %v", err)
	} else {
		fmt.Printf("get temps 2 %+v \n", temps)
	}
}

func PinPointSendEmailByTemp(ppClient *email.AwsPinPoint) {
	var ctx = context.Background()

	msgId, err := ppClient.SendByTemplate(ctx,
		"tianhai@xxxxxx.com",
		//"AgeOfColossus@allxxx.com",
		[]string{"773821422@qq.com", "wojiaoju@outlook.com"},
		"arn:aws:mobiletargeting:us-west-2:533150811650:templates/my_msg_tmp/EMAIL",
		map[string]string{"title": "cccc_t", "h2": "well~"})
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	} else {
		fmt.Printf("send success re msg %s \n", msgId)
	}
}

func PinPointSendEmailByTemp2(ppClient *email.AwsPinPoint) {
	var ctx = context.Background()

	msgId, err := ppClient.SendByTemplate2(ctx,
		"tianhai@xxxxxx.com",
		//"AgeOfColossus@allxxx.com",
		[]string{"773821422@qq.com"},
		"my_msg_tmp",
		map[string]string{"title": "cccc", "h2": "well~"})
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	} else {
		fmt.Printf("send success re msg %s \n", msgId)
	}
}

func PinPointCreateTemplate(ppClient *email.AwsPinPoint) {
	var ctx = context.Background()

	err := ppClient.CreateTemplate(ctx,
		"create_by_program_"+cast.ToString(time.Now().Minute()),
		"Hi {{name}}",
		"<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n  <meta charset=\"UTF-8\">\n <title>{{title}}</title></head><body><h1>{{text}}</h1></body></html>",
		"Hi {{title}} , welcome to {{text}}")
	if err != nil {
		fmt.Printf("err:%v", err)
		return
	} else {
		fmt.Printf("send success re msg \n ")
	}
}

func PinPointDeleteTemp(ppClient *email.AwsPinPoint) {
	var ctx = context.Background()

	err := ppClient.DeleteTemplate(ctx, "create_by_program_34")
	if err != nil {
		fmt.Printf("delete err:%v", err)
		return
	} else {
		fmt.Printf("delete email sucess \n ")
	}
}

func main() {
	//SesDemo	()

	client, _ := GetClient()

	PinPointSendEmailByTemp2(client)
	fmt.Printf("send email by temp 2\n")
	time.Sleep(30 * time.Second)

	PinPointCreateTemplate(client)
	fmt.Printf("create email temp\n")

	PinPointDeleteTemp(client)
	fmt.Printf("delete email temp \n")

	PinPointSendEmailByTemp(client)
	fmt.Printf("send email by temp \n")
	time.Sleep(30 * time.Second)

	PinPointDemo(client)
	fmt.Printf("send msg end\n")
	time.Sleep(30 * time.Second)

	ListTemp2(client)
	fmt.Printf("list temp2 end\n")
	time.Sleep(2 * time.Second)

	ListTemp(client)
	fmt.Printf("list temp end\n")
	time.Sleep(2 * time.Second)
}
