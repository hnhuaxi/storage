package storage

import (
	"os"
)

// Storage 简易对象存储接口
type Storage interface {
	List(prefix string) ([]os.FileInfo, error)
	Get(key string) ([]byte, error)
	PutFile(key string, file string) error
	Put(key string, val []byte) error
	Move(dest string, from string) error
	Remove(key string) error
	Exist(key string) bool
	BucketName() string
	WebURL(key string) (string, error)
	BucketURI(key string) BucketURI
}

// FastdfsStorage fastdfs 对象存储
type FastdfsStorage struct {
	Endpoint string
	Token    string
}
