package email

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pinpoint"
	"github.com/aws/aws-sdk-go/service/pinpointemail"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/tool/config"
	"regexp"
)

type AwsPinPoint struct {
	region  string `json:"regin"`
	appId   string
	charset *string

	sess *session.Session

	pinPointClient *pinpointemail.PinpointEmail

	pinClient *pinpoint.Pinpoint
}

func StrSliceToAwsStrSlice(s []string) []*string {
	var re = []*string{}
	for _, i := range s {
		re = append(re, aws.String(i))
	}
	return re
}

var htmlTagPattern = regexp.MustCompile(`<[\w\d\-]+?>`)

func isHtml(s string) bool {
	return htmlTagPattern.MatchString(s)
}

// InitAwsPinPoint
//
//	@Description: 初始化 pin point 邮件服务
//	@param conf  body any true "-"
//	@return app
//	@return err
func InitAwsPinPoint(conf *config.EmailConfig) (app *AwsPinPoint, err error) {
	app = &AwsPinPoint{region: conf.Region, appId: conf.AppId, charset: aws.String("UTF-8")}
	app.sess, err = session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: credentials.NewStaticCredentials(conf.AccessID, conf.AccessKey, ""),
	})
	if err != nil {
		return app, err
	}

	app.pinPointClient = pinpointemail.New(app.sess)
	app.pinClient = pinpoint.New(app.sess)
	return app, err
}

// Send
//
//	@Description: 通过 pin point 发送邮件
//	@receiver a
//	@param ctx  body any true "-"
//	@param from  body any true "-"
//	@param to  body any true "-"
//	@param title  body any true "-"
//	@param content  body any true "-"
//	@return mailId
//	@return err
func (a *AwsPinPoint) Send(ctx context.Context, from string, to []string, title, content string) (mailId string, err error) {

	var msgInput = pinpointemail.SendEmailInput{
		Destination: &pinpointemail.Destination{
			ToAddresses: StrSliceToAwsStrSlice(to),
		},
		FromEmailAddress: aws.String(from),
		Content: &pinpointemail.EmailContent{
			Simple: &pinpointemail.Message{
				Body: &pinpointemail.Body{
					Html: nil,
					Text: nil,
				},
				Subject: &pinpointemail.Content{
					Charset: a.charset,
					Data:    aws.String(title),
				},
			},
		},
	}
	var body = aws.String(content)
	if isHtml(content) {
		msgInput.Content.Simple.Body.Html = &pinpointemail.Content{
			Charset: a.charset,
			Data:    body,
		}
	} else {
		msgInput.Content.Simple.Body.Text = &pinpointemail.Content{
			Charset: a.charset,
			Data:    body,
		}
	}
	if req, err := a.pinPointClient.SendEmailWithContext(ctx, &msgInput); err != nil {
		return "", err
	} else {
		return *req.MessageId, err
	}
}

// SendByTemplate
//
//	@Description: 通过 pin point 的模版和参数发送邮件
//	@receiver a
//	@param ctx  body any true "-"
//	@param from  body any true "-"
//	@param to  body any true "-"
//	@param templateArn  body any true "-"
//	@param params  body any true "-"
//	@return string
//	@return error
func (a *AwsPinPoint) SendByTemplate(ctx context.Context, from string, to []string,
	templateArn string, params map[string]string) (string, error) {

	kvString, err := jsoniter.MarshalToString(&params) // 序列化
	if err != nil {
		return "", err
	}
	var msgInput = &pinpointemail.SendEmailInput{
		Content: &pinpointemail.EmailContent{
			Raw:    nil,
			Simple: nil,
			Template: &pinpointemail.Template{
				TemplateArn:  aws.String(templateArn),
				TemplateData: aws.String(kvString),
			},
		},
		Destination: &pinpointemail.Destination{
			ToAddresses: StrSliceToAwsStrSlice(to),
		},
		FromEmailAddress: aws.String(from),
	}
	if out, err := a.pinPointClient.SendEmailWithContext(ctx, msgInput); err != nil {
		return "", err
	} else {
		return *out.MessageId, err
	}
}

// ListTemplates
//
//	@Description: 获取所有模版
//	@receiver a
//	@param ctx  body any true "-"
//	@return []*pinpoint.TemplateResponse
//	@return error
func (a *AwsPinPoint) ListTemplates(ctx context.Context) ([]*pinpoint.TemplateResponse, error) {
	var (
		reList    = []*pinpoint.TemplateResponse{}
		nextToken = ""
		input     = &pinpoint.ListTemplatesInput{
			PageSize:     aws.String("20"),
			TemplateType: aws.String(pinpoint.TemplateTypeEmail),
		}
	)

	for {
		if list, err := a.pinClient.ListTemplatesWithContext(ctx, input); err != nil {
			return nil, err
		} else if list.TemplatesResponse != nil {
			if list.TemplatesResponse.NextToken != nil && 0 < len(*list.TemplatesResponse.NextToken) {
				nextToken = *list.TemplatesResponse.NextToken
				input.NextToken = aws.String(nextToken)
				//fmt.Printf("next token is %s", nextToken)
			} else {
				// 如果没有下一页了，就退出返回
				return reList, err
			}
			reList = append(reList, list.TemplatesResponse.Item...)
		}
	}
	return reList, nil
}

func (a *AwsPinPoint) ListTemplates2(ctx context.Context) ([]*pinpoint.TemplateResponse, error) {
	// 调用ListTemplates API获取邮箱模板列表
	var input = &pinpoint.ListTemplatesInput{
		TemplateType: aws.String("EMAIL"), //指定模板类型为email
	}
	resp, err := a.pinClient.ListTemplatesWithContext(ctx, input)
	if err != nil {
		fmt.Println("Failed to list email templates", err)
		return nil, err
	} else {
		return resp.TemplatesResponse.Item, err
	}
	return nil, err
}

// CreateTemplate 创建邮件模版
//
//	@Description:
//	@receiver a
//	@param ctx  body any true "-"
//	@param name  body any true "-"
//	@param title  body any true "-"
//	@param html  body any true "-"
//	@param text  body any true "-"
//	@return error
func (a *AwsPinPoint) CreateTemplate(ctx context.Context, name string, title string, html, text string) error {
	input := &pinpoint.CreateEmailTemplateInput{
		TemplateName: aws.String(name),
		EmailTemplateRequest: &pinpoint.EmailTemplateRequest{
			HtmlPart: aws.String(html),
			Subject:  aws.String(title),
			TextPart: aws.String(text),
		}}
	if _, err := a.pinClient.CreateEmailTemplateWithContext(ctx, input); err != nil {
		return err
	} else {
		return nil
	}
}

// DeleteTemplate
//
//	@Description: 删除邮件模版
//	@receiver a
//	@param ctx  body any true "-"
//	@param name  body any true "-"
//	@return error
func (a *AwsPinPoint) DeleteTemplate(ctx context.Context, name string) error {
	input := &pinpoint.DeleteEmailTemplateInput{TemplateName: aws.String(name)}
	_, err := a.pinClient.DeleteEmailTemplateWithContext(ctx, input)
	return err
}

func (a *AwsPinPoint) SendByTemplate2(ctx context.Context, from string, to []string,
	templateName string, params map[string]string) (string, error) {

	var addrs = map[string]*pinpoint.AddressConfiguration{}
	for _, t := range to {
		addrs[t] = &pinpoint.AddressConfiguration{ChannelType: aws.String(pinpoint.ChannelTypeEmail)}
	}
	var msgInput = &pinpoint.SendMessagesInput{
		ApplicationId: aws.String(a.appId),
		MessageRequest: &pinpoint.MessageRequest{
			Addresses: addrs, // 收件人
			MessageConfiguration: &pinpoint.DirectMessageConfiguration{
				EmailMessage: &pinpoint.EmailMessage{
					FromAddress: aws.String(from), // 发件人
					Substitutions: map[string][]*string{
						"title": []*string{aws.String("cccc")},
						"h2":    []*string{aws.String("well")},
					},
				},
				DefaultMessage: &pinpoint.DefaultMessage{
					Substitutions: map[string][]*string{
						"title": []*string{aws.String("cccc")},
						"h2":    []*string{aws.String("well")},
					},
				},
			},
			TemplateConfiguration: &pinpoint.TemplateConfiguration{
				EmailTemplate: &pinpoint.Template{
					Name:    aws.String(templateName),
					Version: aws.String("2"),
				},
			},
		},
	}

	if re, err := a.pinClient.SendMessagesWithContext(ctx, msgInput); err != nil {
		return "", err
	} else if re.MessageResponse != nil {
		return *re.MessageResponse.RequestId, err
	} else {
		return "", err
	}
}
