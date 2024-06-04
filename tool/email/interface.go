package email

import "context"

type Email interface {
	Send(ctx context.Context, from string, to []string, title, content string) (string, error)
	SendByTemplate(ctx context.Context, from string, to []string, templateName string, params map[string]string) (string, error)
}
