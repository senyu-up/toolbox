package email

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	jsoniter "github.com/json-iterator/go"
	"github.com/senyu-up/toolbox/tool/config"
)

type AwsSes struct {
	region string `json:"regin"`

	sess *session.Session

	ses *ses.SES
}

func InitAwsSes(conf *config.EmailConfig) (*AwsSes, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: credentials.NewStaticCredentials(conf.AccessID, conf.AccessKey, ""),
	})

	return &AwsSes{region: conf.Region, sess: sess, ses: ses.New(sess)}, err
}

// SESSend 发送邮件
func (a *AwsSes) Send(ctx context.Context, form string, target []string, title, content string) (mailId string, err error) {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			//抄送
			CcAddresses: []*string{},
			ToAddresses: StrSliceToAwsStrSlice(target),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data:    aws.String(content),
					Charset: aws.String("UTF-8"),
				},
				Text: nil,
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(title),
			},
		},
		Source: aws.String(form),
	}
	result, err := a.ses.SendEmailWithContext(ctx, input)
	if err != nil {
		return "", err
	}
	return *result.MessageId, nil
}

func (a *AwsSes) CreateTemplate(name string, title string, html, text string) error {
	input := &ses.CreateTemplateInput{Template: &ses.Template{
		HtmlPart:     aws.String(html),
		SubjectPart:  aws.String(title),
		TemplateName: aws.String(name),
		TextPart:     aws.String(text),
	}}
	_, err := a.ses.CreateTemplate(input)
	return err
}

func (a *AwsSes) DeleteTemplate(name string) error {
	input := &ses.DeleteTemplateInput{TemplateName: aws.String(name)}
	_, err := a.ses.DeleteTemplate(input)
	return err
}

func (a *AwsSes) ListTemplates() ([]*ses.Template, error) {
	input := &ses.ListTemplatesInput{
		MaxItems:  nil,
		NextToken: nil,
	}
	var names []*string
LOOP:
	outPut, err := a.ses.ListTemplates(input)
	if err != nil {
		return nil, err
	}
	for _, meta := range outPut.TemplatesMetadata {
		names = append(names, meta.Name)
	}
	if outPut.NextToken != nil {
		goto LOOP
	}
	var templates []*ses.Template
	var out *ses.GetTemplateOutput
	for _, name := range names {
		out, err = a.ses.GetTemplate(&ses.GetTemplateInput{TemplateName: name})
		if err != nil {
			return nil, err
		}
		templates = append(templates, out.Template)
	}
	return templates, nil
}

// SESSendByTemplate
// name 模板名称
// 模板标题参数
// 内容参数
func (a *AwsSes) SendByTemplate(ctx context.Context, from string, to []string, templateName string, params map[string]string) (string, error) {
	str, _ := jsoniter.MarshalToString(params)
	input := &ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			ToAddresses: StrSliceToAwsStrSlice(to),
		},
		Template:     aws.String(templateName),
		Source:       &from,
		TemplateData: &str,
	}
	out, err := a.ses.SendTemplatedEmailWithContext(ctx, input)
	if err != nil {
		return "", err
	}
	return *out.MessageId, nil
}
