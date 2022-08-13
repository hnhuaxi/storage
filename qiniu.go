package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/hysios/log"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/sms/bytes"
	"github.com/qiniu/go-sdk/v7/storage"
)

// QiniuStorage qiniu 云对象存储
type QiniuStorage struct {
	Config QiniuConfig
	mac    *auth.Credentials
}

type QiniuConfig struct {
	AppKey     string
	Secret     string
	Bucket     string
	ParentDir  string
	Region     string
	HttpPrefix string
}

var qiniuRegionMap = map[string]storage.Region{
	"huadong":  storage.ZoneHuadong,
	"huabei":   storage.ZoneHuabei,
	"huanan":   storage.ZoneHuanan,
	"beimei":   storage.ZoneBeimei,
	"xinjiapo": storage.ZoneXinjiapo,
}

func NewQiniuStorage(cfg *QiniuConfig) *QiniuStorage {

	store := &QiniuStorage{
		Config: *cfg,
		mac:    qbox.NewMac(cfg.AppKey, cfg.Secret),
	}

	register("qiniu", store, cfg.Bucket)
	return store
}

func (qiniu *QiniuStorage) getKey(key string) string {
	return path.Join(qiniu.Config.ParentDir, key)
}

func (qiniu *QiniuStorage) List(prefix string) ([]os.FileInfo, error) {
	return nil, errors.New("not implemented")
}

func (qiniu *QiniuStorage) Get(key string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

// PutFile 上传一个文件
func (qiniu *QiniuStorage) PutFile(key string, localfile string) error {
	bucket := qiniu.Config.Bucket

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	upToken := putPolicy.UploadToken(qiniu.mac)
	cfg := storage.Config{}
	// 空间对应的机房

	if region, ok := qiniuRegionMap[qiniu.Config.Region]; ok {
		cfg.Zone = &region
	} else {
		cfg.Zone = &storage.ZoneHuadong
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: nil,
	}
	err := formUploader.PutFile(context.Background(), &ret, upToken, qiniu.getKey(key), localfile, &putExtra)
	if err != nil {
		return err
	}
	log.Debugf("upload to bucket %s -> %s", bucket, ret.Key)
	return nil
}

// PutFile 上传一段 Bytes 数据流
func (qiniu *QiniuStorage) Put(key string, b []byte) error {
	bucket := qiniu.Config.Bucket

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}

	upToken := putPolicy.UploadToken(qiniu.mac)
	cfg := storage.Config{}
	// 空间对应的机房

	if region, ok := qiniuRegionMap[qiniu.Config.Region]; ok {
		cfg.Zone = &region
	} else {
		cfg.Zone = &storage.ZoneHuadong
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{
		Params: nil,
	}

	var rd = bytes.NewReader(b)
	err := formUploader.Put(context.Background(), &ret, upToken, qiniu.getKey(key), rd, int64(len(b)), &putExtra)
	if err != nil {
		return err
	}
	log.Infof("upload to bucket %s -> %s", bucket, qiniu.getKey(key))
	return nil
}

// Move 移动目标到指定位置
func (qiniu *QiniuStorage) Move(dest string, from string) error {
	bucket := qiniu.Config.Bucket
	cfg := storage.Config{}
	// 空间对应的机房

	if region, ok := qiniuRegionMap[qiniu.Config.Region]; ok {
		cfg.Zone = &region
	} else {
		cfg.Zone = &storage.ZoneHuadong
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	bucketManager := storage.NewBucketManager(qiniu.mac, &cfg)

	return bucketManager.Move(bucket, from, bucket, dest, true)
}

// Exist 存储空间存在一个文件
func (qiniu *QiniuStorage) Exist(key string) bool {
	bucket := qiniu.Config.Bucket
	cfg := storage.Config{}
	// 空间对应的机房

	if region, ok := qiniuRegionMap[qiniu.Config.Region]; ok {
		cfg.Zone = &region
	} else {
		cfg.Zone = &storage.ZoneHuadong
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	bucketManager := storage.NewBucketManager(qiniu.mac, &cfg)
	fileInfo, err := bucketManager.Stat(bucket, qiniu.getKey(key))
	if err != nil {
		return false
	}
	return fileInfo.Fsize > 0
}

func (qiniu *QiniuStorage) Remove(key string) error {
	bucket := qiniu.Config.Bucket
	cfg := storage.Config{}
	// 空间对应的机房

	if region, ok := qiniuRegionMap[qiniu.Config.Region]; ok {
		cfg.Zone = &region
	} else {
		cfg.Zone = &storage.ZoneHuadong
	}
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	bucketManager := storage.NewBucketManager(qiniu.mac, &cfg)

	return bucketManager.Delete(bucket, qiniu.getKey(key))
}

func (qiniu *QiniuStorage) WebURL(key string) (string, error) {
	log.Infof("config %v", qiniu.Config)
	u, err := url.Parse(qiniu.Config.HttpPrefix)
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, qiniu.getKey(key))
	log.Infof("new url %s", u.String())
	return u.String(), nil
}

func (qiniu *QiniuStorage) BucketName() string {
	return qiniu.Config.Bucket
}

func (qiniu *QiniuStorage) BucketURI(key string) BucketURI {
	return BucketURI(fmt.Sprintf("%s://%s/%s", "qiniu", qiniu.Config.Bucket, qiniu.getKey(key)))
}

var _ Storage = &QiniuStorage{}
