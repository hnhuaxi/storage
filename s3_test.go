package storage

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

func TestS3ObjectStorage_List(t *testing.T) {
	var (
		key    = os.Getenv("AWS_ACCESS_KEY_ID")
		secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
		cred   = credentials.NewStaticCredentials(key, secret, "")
		cfg    = &aws.Config{Credentials: cred,
			Region:   aws.String("cn-hangzhou"),
			Endpoint: aws.String("http://oss-cn-hangzhou.aliyuncs.com"),
		}
		err error
	)

	store, err := NewS3(key,
		secret,
		"pdls2",
		session.New(cfg),
		S3Endpoint("http://oss-cn-hangzhou.aliyuncs.com"),
	)
	assert.NoError(t, err)

	gotObjects, err := store.List("clouds/pushnode/_autoupdate")
	assert.GreaterOrEqual(t, len(gotObjects), 1)

}

func TestS3ObjectStorage_Get(t *testing.T) {
	var (
		key    = os.Getenv("AWS_ACCESS_KEY_ID")
		secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
		cred   = credentials.NewStaticCredentials(key, secret, "")
		cfg    = &aws.Config{Credentials: cred,
			Region:   aws.String("cn-hangzhou"),
			Endpoint: aws.String("http://oss-cn-hangzhou.aliyuncs.com"),
		}
		err error
	)

	store, err := NewS3(key,
		secret,
		"pdls2",
		session.New(cfg),
		S3Endpoint("http://oss-cn-hangzhou.aliyuncs.com"),
	)
	assert.NoError(t, err)

	content, err := store.Get("clouds/pushnode/_autoupdate/LATEST")
	t.Logf("content %s", string(content))
	assert.NoError(t, err)
	assert.Greater(t, len(content), 1)
	assert.Equal(t, uint8(content[0]), uint8('v'))
}

func TestS3ObjectStorage_Put(t *testing.T) {
	var (
		key    = os.Getenv("AWS_ACCESS_KEY_ID")
		secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
		cred   = credentials.NewStaticCredentials(key, secret, "")
		cfg    = &aws.Config{Credentials: cred,
			Region:   aws.String("cn-hangzhou"),
			Endpoint: aws.String("http://oss-cn-hangzhou.aliyuncs.com"),
		}
		err error
	)

	store, err := NewS3(key,
		secret,
		"pdls2",
		session.New(cfg),
		S3Endpoint("http://oss-cn-hangzhou.aliyuncs.com"),
	)
	assert.NoError(t, err)

	err = store.Put("clouds/test/_hello", []byte("hello world"))
	assert.NoError(t, err)
}

func TestS3ObjectStorage_Move(t *testing.T) {
	var (
		key    = os.Getenv("AWS_ACCESS_KEY_ID")
		secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
		cred   = credentials.NewStaticCredentials(key, secret, "")
		cfg    = &aws.Config{Credentials: cred,
			Region:   aws.String("cn-hangzhou"),
			Endpoint: aws.String("http://oss-cn-hangzhou.aliyuncs.com"),
		}
		err error
	)

	store, err := NewS3(key,
		secret,
		"pdls2",
		session.New(cfg),
		S3Endpoint("http://oss-cn-hangzhou.aliyuncs.com"),
	)
	assert.NoError(t, err)

	err = store.Put("hello3.txt", []byte("hello world"))
	assert.NoError(t, err)

	err = store.Move("/hello4.txt", "hello3.txt")
	assert.NoError(t, err)

	assert.True(t, store.Exist("hello4.txt"))
	assert.False(t, store.Exist("hello3.txt"))

	err = store.Put("test/hello3.txt", []byte("hello world"))
	assert.NoError(t, err)

	err = store.Move("hello4.txt", "test/hello3.txt")
	assert.NoError(t, err)

	assert.True(t, store.Exist("hello4.txt"))
	assert.False(t, store.Exist("test/hello3.txt"))
}

func TestS3ObjectStorage_Remove(t *testing.T) {
	var (
		key    = os.Getenv("AWS_ACCESS_KEY_ID")
		secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
		cred   = credentials.NewStaticCredentials(key, secret, "")
		cfg    = &aws.Config{Credentials: cred,
			Region:   aws.String("cn-hangzhou"),
			Endpoint: aws.String("http://oss-cn-hangzhou.aliyuncs.com"),
		}
		err error
	)

	store, err := NewS3("LTAI4FhtiW3CF5BuDdNVu6jG",
		"qfCRObTXX8X1ovWgWwfkR2UWixY5et",
		"pdls2",
		session.New(cfg),
		S3Endpoint("http://oss-cn-hangzhou.aliyuncs.com"),
	)
	assert.NoError(t, err)
	const filename = "clouds/test/_hello"
	err = store.Put(filename, []byte("hello world"))
	assert.NoError(t, err)
	err = store.Remove(filename)
	assert.NoError(t, err)
	if ok := store.Exist(filename); ok {
		t.Fatalf("must deleted %s file", filename)
	}
}

func TestS3ObjectStorage_Exist(t *testing.T) {
	var (
		key    = os.Getenv("AWS_ACCESS_KEY_ID")
		secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
		cred   = credentials.NewStaticCredentials(key, secret, "")
		cfg    = &aws.Config{Credentials: cred,
			Region:   aws.String("cn-hangzhou"),
			Endpoint: aws.String("http://oss-cn-hangzhou.aliyuncs.com"),
		}
		err error
	)

	store, err := NewS3(key,
		secret,
		"pdls2",
		session.New(cfg),
		S3Endpoint("http://oss-cn-hangzhou.aliyuncs.com"),
	)
	assert.NoError(t, err)

	if ok := store.Exist("clouds/test/is_not_exist.file"); ok {
		t.Fatalf("can't found %s file", "clouds/test/is_not_exist.file")
	}

}
