package s3

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"io"
	"strconv"
	"strings"
)

type Config struct {
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	SecretAccessKey string
}

type Client struct {
	client *minio.Client

	bucket string
}

func NewClient(c *Config) (*Client, error) {
	client, err := minio.New(c.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(c.AccessKeyID, c.SecretAccessKey, ""),
	})
	if err != nil {
		return nil, errors.Wrap(err, "creating minio s3 client")
	}

	return &Client{client: client, bucket: c.Bucket}, nil
}

func logPath(id int32) string {
	return "logs/" + viper.GetString("ENVIRONMENT") + "/" + strconv.Itoa(int(id))
}

func (c *Client) SaveEncounter(ctx context.Context, id int32, raw string) error {
	reader := strings.NewReader(raw)
	_, err := c.client.PutObject(ctx, c.bucket, logPath(id), reader, int64(reader.Len()), minio.PutObjectOptions{})
	return err
}

func (c *Client) FetchEncounter(ctx context.Context, id int32) ([]byte, error) {
	reader, err := c.client.GetObject(ctx, c.bucket, logPath(id), minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "fetching encounter from s3")
	}

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "reading object")
	}
	return raw, nil
}

func (c *Client) DeleteEncounter(ctx context.Context, id int32) error {
	return c.client.RemoveObject(ctx, c.bucket, logPath(id), minio.RemoveObjectOptions{})
}
