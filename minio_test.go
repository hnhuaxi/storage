package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinioStorage_List(t *testing.T) {
	var (
		err    error
		key    = os.Getenv("MINIO_ACCESS_KEY")
		secret = os.Getenv("MINIO_SECRET_KEY")
		bucket = os.Getenv("MINIO_BUCKET")
	)

	store, err := NewMinio(key, secret, bucket, MinioEndpoint("localhost:9000"))
	assert.NoError(t, err)

	gotObjects, err := store.List("/")
	assert.GreaterOrEqual(t, len(gotObjects), 1)
}

func TestMinioStorage_Get(t *testing.T) {
	var (
		key    = os.Getenv("MINIO_ACCESS_KEY")
		secret = os.Getenv("MINIO_SECRET_KEY")
		bucket = os.Getenv("MINIO_BUCKET")
	)
	store, err := NewMinio("minioadmin", "minioadmin", "forensics", MinioEndpoint("localhost:9000"))
	assert.NoError(t, err)

	content, err := store.Get("/3.jpg")
	assert.NoError(t, err)
	assert.Greater(t, len(content), 1)
	// assert.Equal(t, uint8(content[0]), uint8('v'))
}

func TestMinioStorage_Put(t *testing.T) {
	var (
		key    = os.Getenv("MINIO_ACCESS_KEY")
		secret = os.Getenv("MINIO_SECRET_KEY")
		bucket = os.Getenv("MINIO_BUCKET")
	)
	store, err := NewMinio(key, secret, bucket, MinioEndpoint("localhost:9000"))
	assert.NoError(t, err)

	err = store.Put("/hello.txt", []byte("hello world"))
	assert.NoError(t, err)

	content, err := store.Get("/hello.txt")
	assert.NoError(t, err)
	assert.Greater(t, len(content), 1)

	content, err = store.Get("hello.txt")
	assert.NoError(t, err)
	assert.Greater(t, len(content), 1)
}

func TestMinioStorage_Move(t *testing.T) {
	var (
		key    = os.Getenv("MINIO_ACCESS_KEY")
		secret = os.Getenv("MINIO_SECRET_KEY")
		bucket = os.Getenv("MINIO_BUCKET")
	)
	store, err := NewMinio(key, secret, bucket, MinioEndpoint("localhost:9000"))
	assert.NoError(t, err)

	err = store.Put("hello3.txt", []byte("hello world"))
	assert.NoError(t, err)

	err = store.Move("hello4.txt", "hello3.txt")
	assert.NoError(t, err)

	assert.True(t, store.Exist("hello4.txt"))
	assert.False(t, store.Exist("hello3.txt"))

	err = store.Put("/test/hello3.txt", []byte("hello world"))
	assert.NoError(t, err)

	err = store.Move("hello4.txt", "test/hello3.txt")
	assert.NoError(t, err)

	assert.True(t, store.Exist("hello4.txt"))
	assert.False(t, store.Exist("test/hello3.txt"))
}

func TestMinioStorage_Remove(t *testing.T) {
	var (
		key    = os.Getenv("MINIO_ACCESS_KEY")
		secret = os.Getenv("MINIO_SECRET_KEY")
		bucket = os.Getenv("MINIO_BUCKET")
	)
	store, err := NewMinio(key, secret, bucket, MinioEndpoint("localhost:9000"))
	assert.NoError(t, err)

	err = store.Put("hello2.txt", []byte("hello world"))
	assert.NoError(t, err)

	err = store.Remove("hello2.txt")
	assert.NoError(t, err)

}

func TestMinioStorage_Exist(t *testing.T) {
	var (
		key    = os.Getenv("MINIO_ACCESS_KEY")
		secret = os.Getenv("MINIO_SECRET_KEY")
		bucket = os.Getenv("MINIO_BUCKET")
	)
	store, err := NewMinio(key, secret, bucket, MinioEndpoint("localhost:9000"))
	assert.NoError(t, err)

	exist := store.Exist("hello3.txt")
	assert.Equal(t, exist, false)

	err = store.Put("hello4.txt", []byte("hello world"))
	assert.NoError(t, err)

	exist = store.Exist("hello4.txt")
	assert.Equal(t, exist, true)
}

func TestMinioStorage_WebURL(t *testing.T) {
	var (
		key    = os.Getenv("MINIO_ACCESS_KEY")
		secret = os.Getenv("MINIO_SECRET_KEY")
		bucket = os.Getenv("MINIO_BUCKET")
	)
	store, err := NewMinio(key, secret, bucket,
		MinioEndpoint("localhost:9000"),
		MinioWebPrefix("http://localhost:9000/"),
	)
	assert.NoError(t, err)

	url, err := store.WebURL("test.jpg")
	assert.NoError(t, err)
	assert.Equal(t, url, "http://localhost:9000/test.jpg")
}
