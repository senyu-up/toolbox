package google

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/encrypt"
	"google.golang.org/api/option"
)

func NewFirebase(ctx context.Context, cnf config.Firebase) (app *firebase.App, err error) {
	var opt option.ClientOption
	if cnf.CredentialsFile != "" {
		opt = option.WithCredentialsFile(cnf.CredentialsFile)
	} else {
		deData, err := encrypt.Base64Decode([]byte(cnf.CredentialsJson))
		if err != nil {
			return nil, err
		}
		opt = option.WithCredentialsJSON(deData)
	}

	c := &firebase.Config{
		AuthOverride:     nil,
		DatabaseURL:      cnf.DatabaseURL,
		ProjectID:        cnf.ProjectId,
		ServiceAccountID: cnf.AccountId,
		StorageBucket:    cnf.Bucket,
	}

	return firebase.NewApp(ctx, c, opt)
}
