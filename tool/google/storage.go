package google

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/senyu-up/toolbox/tool/config"
	"github.com/senyu-up/toolbox/tool/encrypt"
	"google.golang.org/api/option"
)

type Storage struct {
	*storage.BucketHandle
}

func NewStorageClient(ctx context.Context, cnf config.Gcs, opts ...option.ClientOption) (*Storage, error) {
	byteData, err := encrypt.Base64Decode([]byte(cnf.CredentialsJson))
	if err != nil {
		return nil, err
	}
	opts = append(opts, option.WithCredentialsJSON(byteData))
	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	// 使用 GCS 客户端执行操作
	cli := &Storage{}

	cli.BucketHandle = client.Bucket(cnf.Bucket)

	return cli, nil
}
