package etcd

import (
	"github.com/coreos/etcd/clientv3"
	"go.uber.org/zap"
	"myCron/mylog"
	"time"
)

// Config ...
type (
	Config struct {
		Endpoints []string `json:"endpoints"`
		CertFile  string   `json:"certFile"`
		KeyFile   string   `json:"keyFile"`
		CaCert    string   `json:"caCert"`
		BasicAuth bool     `json:"basicAuth"`
		UserName  string   `json:"userName"`
		Password  string   `json:"-"`
		// 连接超时时间
		ConnectTimeout time.Duration `json:"connectTimeout"`
		Secure         bool          `json:"secure"`
		// 自动同步member list的间隔
		AutoSyncInterval time.Duration `json:"autoAsyncInterval"`
		TTL              int           // 单位：s
		logger           *mylog.Logger
	}
)

// Duration ...
// panic if parse duration failed
func Duration(str string) time.Duration {
	dur, err := time.ParseDuration(str)
	if err != nil {
		panic(err)
	}

	return dur
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BasicAuth:      false,
		ConnectTimeout: Duration("5s"),
		Secure:         false,
		logger:         mylog.DefaultLogger,
	}
}

// New ...
func NewClient(config *Config) *clientv3.Client {
	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
	}

	if config.Endpoints == nil {
		config.logger.Panic("client etcd endpoints empty", zap.String("err ", "client.etcd"))
	}

	client, err := clientv3.New(conf)
	if err != nil {
		config.logger.Panic("client etcd start panic", zap.String("err ", "client.etcd"), zap.Error(err))
	}

	config.logger.Info("dial etcd server")
	return client
}
