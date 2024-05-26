package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/cockroachdb/errors"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
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

func LogPath(id int32, suffix string) string {
	return EnvPath(fmt.Sprintf("logs/%d/%s", id, suffix))
}

func EnvPath(suffix string) string {
	return viper.GetString("ENVIRONMENT") + "/" + suffix
}

const (
	pathRaw    = "raw"
	pathParsed = "parsed"
)

func (c *Client) SaveEncounter(ctx context.Context, id int32, raw, parsed []byte) error {
	files := map[string][]byte{
		pathRaw:    raw,
		pathParsed: parsed,
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errs := make(chan error)
	for file, data := range files {
		go func(file string, data []byte) {
			reader := bytes.NewReader(data)
			_, err := c.client.PutObject(
				ctx,
				c.bucket,
				LogPath(id, file),
				reader,
				int64(reader.Len()),
				minio.PutObjectOptions{},
			)
			errs <- errors.Wrapf(err,
				"saving encounter %d (%s) to s3", id, file,
			)
		}(file, data)
	}

	for i := 0; i < 2; i++ {
		if err := <-errs; err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) FetchEncounter(ctx context.Context, id int32) ([]byte, error) {
	reader, err := c.client.GetObject(ctx, c.bucket, LogPath(id, pathParsed), minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "fetching encounter from s3")
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "reading object")
	}
	return raw, nil
}

func (c *Client) FetchEncounterOldRaw(ctx context.Context, id int32) ([]byte, error) {
	reader, err := c.client.GetObject(
		ctx, c.bucket,
		fmt.Sprintf("logs/%s/%d", viper.Get("ENVIRONMENT"), id),
		minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "fetching encounter from s3")
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "reading object")
	}
	return raw, nil
}

func (c *Client) DeleteEncounter(ctx context.Context, id int32) error {
	return c.client.RemoveObject(ctx, c.bucket, LogPath(id, ""), minio.RemoveObjectOptions{})
}

func (c *Client) FetchImage(ctx context.Context, path string) ([]byte, error) {
	reader, err := c.client.GetObject(ctx, c.bucket, "images/"+viper.GetString("ENVIRONMENT")+"/"+path, minio.GetObjectOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "fetching image from s3")
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrap(err, "reading object")
	}
	return raw, nil
}

func (c *Client) FetchAvatar(ctx context.Context, uuid string) ([]byte, error) {
	return c.FetchImage(ctx, "avatar/"+uuid)
}

func (c *Client) SaveImage(ctx context.Context, path string, raw []byte) error {
	reader := bytes.NewReader(raw)
	_, err := c.client.PutObject(ctx, c.bucket, "images/"+viper.GetString("ENVIRONMENT")+"/"+path, reader, int64(reader.Len()), minio.PutObjectOptions{})
	return err
}

func (c *Client) RemoveImage(ctx context.Context, path string) error {
	return c.client.RemoveObject(ctx, c.bucket, "images/"+viper.GetString("ENVIRONMENT")+path, minio.RemoveObjectOptions{})
}

//func (c *Client) PurgeLogs(ctx context.Context, ids []int32) error {
//	objectChan := make(chan minio.ObjectInfo)
//	removeErrs := c.client.RemoveObjects(ctx, c.bucket, objectChan, minio.RemoveObjectsOptions{})
//	for _, id := range ids {
//		objectChan <- minio.ObjectInfo{Key: logPath(id)}
//		objectChan <- minio.ObjectInfo{Key: "images/" + viper.GetString("ENVIRONMENT") + "/thumbnail/" + strconv.Itoa(int(id))}
//	}
//	close(objectChan)
//
//	errs := &multierror.Error{}
//	for err := range removeErrs {
//		errs = multierror.Append(errs, err.Err)
//	}
//	return errs
//}

func (c *Client) SaveAvatar(ctx context.Context, uuid string, raw []byte) error {
	return c.SaveImage(ctx, "avatar/"+uuid, raw)
}

func (c *Client) RemoveAvatar(ctx context.Context, uuid string) error {
	return c.RemoveImage(ctx, "avatar/"+uuid)
}
