package storage

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type BucketURI string

var (
	UnknownURI = "http://127.0.0.1:9000/keyue/unknown.png"
)

func SetUnknownURI(uri string) {
	UnknownURI = uri
}

func (uri BucketURI) UnknownURI() string {
	return UnknownURI
}

func (uri BucketURI) String() string {
	if Empty(string(uri)) {
		return uri.UnknownURI()
	}

	u, err := url.Parse(string(uri))
	if err != nil {
		return string(uri)
	}

	host, ok := GetBucketHost(u.Scheme, u.Host)
	if !ok {
		return string(uri)
	}

	// if isPrivatehost(host) {
	// 	if !config.GetBool("cloudmode.private") {
	// 		utils.ExternalIP()
	// 	}
	// }

	switch u.Scheme {
	case "minio", "s3", "qiniu":
		return fmt.Sprintf("http://%s/%s%s", host, u.Host, u.Path)
	default:
		return string(uri)
	}
}

func isPrivatehost(host string) bool {
	_host, _, _ := net.SplitHostPort(host)
	if _host == "localhost" {
		return true
	}

	if ip := net.ParseIP(_host); ip == nil {
		return false
	} else {
		return IsPrivateIP(ip)
	}
}

func (uri BucketURI) MarshalJSON() ([]byte, error) {
	return json.Marshal(uri.String())
}

func Empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func QuoteBytes(u string) []byte {
	return []byte(strconv.Quote(u))
}

var stores sync.Map

type storeHost struct {
	Bucket string
	Host   string
	store  Storage
}

func register(scheme string, store Storage, host string) {
	var sts = []storeHost{{
		Bucket: store.BucketName(),
		Host:   host,
		store:  store,
	}}

	if olds, load := stores.LoadOrStore(scheme, sts); load {
		if stss, ok := olds.([]storeHost); ok {
			stss = append(stss, sts...)
			stores.Store(scheme, stss)
		}
	}
}

func GetBucketHost(scheme string, bucket string) (string, bool) {
	sts, ok := stores.Load(scheme)
	if !ok {
		return "", false
	}
	if stss, ok := sts.([]storeHost); !ok {
		return "", false
	} else {
		for _, st := range stss {
			if st.Bucket == bucket {
				return st.Host, true
			}
		}
	}

	return "", false
}
