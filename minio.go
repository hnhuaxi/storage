package storage

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Endpoint   string
	AccessKey  string
	AppSecret  string
	Region     string
	Bucket     string
	HttpPrefix string
	UseSSL     bool

	client *minio.Client
}

func MinioWebPrefix(url string) MinioOptionFunc {
	return func(minio *MinioStorage) error {
		minio.HttpPrefix = url
		return nil
	}
}

func MinioEndpoint(url string) MinioOptionFunc {
	return func(minio *MinioStorage) error {
		minio.Endpoint = url
		return nil
	}
}

func MinioUseSSL(useSSL bool) MinioOptionFunc {
	return func(minio *MinioStorage) error {
		minio.UseSSL = useSSL
		return nil
	}
}

type MinioOptionFunc func(*MinioStorage) error

func NewMinio(appkey, secret string, bucket string, opts ...MinioOptionFunc) (store *MinioStorage, err error) {
	store = &MinioStorage{
		Bucket:    bucket,
		AccessKey: appkey,
		AppSecret: secret,
		// client:    client,
	}
	for _, set := range opts {
		set(store)
	}

	register("minio", store, store.Hostname())

	client, err := minio.New(store.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4("", appkey, secret),
		Secure: store.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	store.client = client
	return store, nil
}

func (store *MinioStorage) Hostname() string {
	if Empty(store.HttpPrefix) {
		uri, err := url.Parse(store.Endpoint)
		if err != nil {
			return ""
		}
		return uri.String()
	} else {
		return hostname(store.HttpPrefix)
	}
}

func NewMinioV2(appkey, secret string, bucket string, opts ...MinioOptionFunc) (store *MinioStorage, err error) {
	store = &MinioStorage{
		Bucket:    bucket,
		AccessKey: appkey,
		AppSecret: secret,
		// client:    client,
	}
	for _, set := range opts {
		set(store)
	}

	cred := credentials.NewStaticV2(appkey, secret, "")
	client, err := minio.New(store.Endpoint, &minio.Options{
		Creds: cred,
	})
	// client, err := minio.New(store.Endpoint, appkey, secret, store.UseSSL)
	if err != nil {
		return nil, err
	}
	store.client = client
	return store, nil
}

func (store *MinioStorage) List(prefix string) ([]os.FileInfo, error) {
	var (
		ctx = context.Background()
	)

	// Indicate to our routine to exit cleanly upon return.

	var result = make([]os.FileInfo, 0)
	isRecursive := true

	objectCh := store.client.ListObjects(ctx, store.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: isRecursive,
	})
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return nil, object.Err
		}

		var obj = &ObjectInfo{key: object.Key, size: object.Size, time: object.LastModified}
		if obj.key[len(obj.key)-1] == '/' {
			obj.isDir = true
		}

		result = append(result, obj)
	}
	return result, nil
}

func (store *MinioStorage) Get(key string) ([]byte, error) {
	var (
		ctx = context.Background()
	)

	object, err := store.client.GetObject(ctx, store.Bucket, key, minio.GetObjectOptions{})
	key = strings.TrimPrefix(key, "/")
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(object)
}

func (store *MinioStorage) PutFile(key string, file string) error {
	var (
		ctx = context.Background()
	)
	// 使用FPutObject上传一个zip文件。
	key = strings.TrimPrefix(key, "/")
	_, err := store.client.FPutObject(ctx, store.Bucket, key, file, minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (store *MinioStorage) Put(key string, val []byte) error {
	var (
		ctx = context.Background()
		buf = bytes.NewBuffer(val)
	)
	key = strings.TrimPrefix(key, "/")
	_, err := store.client.PutObject(ctx, store.Bucket, key, buf, int64(len(val)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (store *MinioStorage) Move(dest string, from string) error {
	var (
		ctx = context.Background()
	)
	srcOpts := minio.CopySrcOptions{
		Bucket: store.Bucket,
		Object: from,
	}

	dest = strings.TrimPrefix(dest, "/")
	from = strings.TrimPrefix(from, "/")

	dstOpts := minio.CopyDestOptions{
		Bucket: store.Bucket,
		Object: dest,
	}

	// dstOpts, err := minio.NewDestinationInfo(store.Bucket, dest, nil, nil)
	// if err != nil {
	// 	return err
	// }

	_, err := store.client.CopyObject(ctx, dstOpts, srcOpts)
	if err != nil {
		return err
	}

	return store.Remove(from)
}

func (store *MinioStorage) Remove(key string) error {
	var (
		ctx = context.Background()
	)
	err := store.client.RemoveObject(ctx, store.Bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (store *MinioStorage) Exist(key string) bool {
	var (
		ctx = context.Background()
	)
	_, err := store.client.StatObject(ctx, store.Bucket, key, minio.StatObjectOptions{})
	if err != nil {
		return false
	}
	return true
}

func (store *MinioStorage) hasHttpPrefix(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

func (store *MinioStorage) WebURL(key string) (string, error) {
	var (
		u   *url.URL
		err error
	)

	if store.hasHttpPrefix(store.HttpPrefix) {
		u, err = url.Parse(store.HttpPrefix)
		if err != nil {
			return "", err
		}
	} else {
		u, err = url.Parse("http://" + store.HttpPrefix)
		if err != nil {
			return "", err
		}
	}

	u.Path = path.Join(u.Path, key)
	return u.String(), nil
}

func (store *MinioStorage) BucketName() string {
	return store.Bucket
}

func (store *MinioStorage) BucketURI(key string) BucketURI {
	return BucketURI(fmt.Sprintf("%s://%s/%s", "minio", store.Bucket, key))
}

var _ Storage = &MinioStorage{}
