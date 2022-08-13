package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3ObjectStorage S3 简单存储对象
type S3ObjectStorage struct {
	Endpoint   string
	AccessKey  string
	AppSecret  string
	Region     string
	Bucket     string
	HttpPrefix string

	svc *s3.S3
}

func S3WebPrefix(url string) S3OptionFunc {
	return func(s3 *S3ObjectStorage) error {
		s3.HttpPrefix = url
		return nil
	}
}

func NewS3(appkey, secret string, bucket string, sess *session.Session, opts ...S3OptionFunc) (store *S3ObjectStorage, err error) {
	store = &S3ObjectStorage{
		Bucket:    bucket,
		AccessKey: appkey,
		AppSecret: secret,
		svc:       s3.New(sess),
	}

	for _, opt := range opts {
		opt(store)
	}

	register("s3", store, store.Hostname())

	return store, nil
}

func hostname(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return s
	}
	return u.Host
}

func (store *S3ObjectStorage) Hostname() string {
	if Empty(store.HttpPrefix) {
		uri, err := store.BucketWebsite()
		if err != nil {
			return ""
		}
		return hostname(uri)
	} else {
		return hostname(store.HttpPrefix)
	}
}

// List 列出 S3 Object 清单
func (store *S3ObjectStorage) List(prefix string) (objects []os.FileInfo, err error) {
	input := &s3.ListObjectsInput{
		Bucket: &store.Bucket,
		// Delimiter: aws.String("/"),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(100),
	}

	result, err := store.svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return nil, err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, err
		}
	}

	objects = make([]os.FileInfo, 0, len(result.Contents))
	for _, cont := range result.Contents {

		var obj = &ObjectInfo{key: *cont.Key, size: *cont.Size, time: *cont.LastModified}
		if obj.key[len(obj.key)-1] == '/' {
			obj.isDir = true
		}

		objects = append(objects, obj)
	}
	return
}

// Get 获取 S3 Object 对象
func (store *S3ObjectStorage) Get(key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(store.Bucket),
		Key:    aws.String(key),
		// Range:  aws.String("bytes=0-9"),
	}

	result, err := store.svc.GetObject(input)
	if err != nil {

		return nil, err
	}

	return ioutil.ReadAll(result.Body)
}

func (store *S3ObjectStorage) Put(key string, val []byte) error {
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewBuffer(val)),
		Bucket: aws.String(store.Bucket),
		Key:    aws.String(key),
		// ServerSideEncryption: aws.String("AES256"),
		// StorageClass:         aws.String("STANDARD_IA"),
	}

	_, err := store.svc.PutObject(input)
	if err != nil {

		return err
	}

	return nil
}

func (store *S3ObjectStorage) PutFile(key string, file string) error {
	f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(f),
		Bucket: aws.String(store.Bucket),
		Key:    aws.String(key),
		// ServerSideEncryption: aws.String("AES256"),
		// StorageClass:         aws.String("STANDARD_IA"),
	}

	_, err = store.svc.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}

func (store *S3ObjectStorage) Move(dest string, from string) error {

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(store.Bucket),
		CopySource: aws.String(path.Join(store.Bucket, from)),
		Key:        aws.String(dest),
		// ServerSideEncryption: aws.String("AES256"),
		// StorageClass:         aws.String("STANDARD_IA"),
	}
	_, err := store.svc.CopyObject(input)
	if err != nil {

		return err
	}

	return store.Remove(from)
}

func (store *S3ObjectStorage) Remove(key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(store.Bucket),
		Key:    aws.String(key),
	}

	_, err := store.svc.DeleteObject(input)
	if err != nil {

		return err
	}
	return nil
}

func (store *S3ObjectStorage) Exist(key string) bool {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(store.Bucket),
		Key:    aws.String(key),
	}

	_, err := store.svc.HeadObject(input)
	if err != nil {
		return false
	}

	return true
}

func (store *S3ObjectStorage) WebURL(key string) (string, error) {
	if Empty(store.HttpPrefix) {
		return store.BucketWebsite()
	}

	u, err := url.Parse(store.HttpPrefix)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, key)
	return u.String(), nil
}

func (store *S3ObjectStorage) BucketWebsite() (string, error) {
	result, err := store.svc.GetBucketWebsite(&s3.GetBucketWebsiteInput{
		Bucket: aws.String(store.Bucket),
	})
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func (store *S3ObjectStorage) BucketName() string {
	return store.Bucket
}

func (store *S3ObjectStorage) BucketURI(key string) BucketURI {
	return BucketURI(fmt.Sprintf("%s://%s/%s", "qiniu", store.Bucket, key))

}

var _ Storage = &S3ObjectStorage{}
