package storage

import (
	"os"
	"time"
)

// ObjectInfo 对象文件信息
type ObjectInfo struct {
	key   string
	size  int64
	time  time.Time
	isDir bool
}

func (obj *ObjectInfo) Name() string {
	return obj.key
}

func (obj *ObjectInfo) Size() int64 {
	return obj.size
}

func (obj *ObjectInfo) Mode() os.FileMode {
	return os.ModePerm
}

func (obj *ObjectInfo) ModTime() time.Time {
	return obj.time
}

func (obj *ObjectInfo) IsDir() bool {
	return obj.isDir
}

func (obj *ObjectInfo) Sys() interface{} {
	return nil
}
