package job

import (
	"go.uber.org/zap"
	"myCron/etcd"
	"myCron/mylog"
	"myCron/parser"
	"myCron/util"

	"github.com/robfig/cron/v3"
)

var (
	myParser = parser.NewParser(parser.Second | parser.Minute | parser.Hour | parser.Dom | parser.Month | parser.Dow | parser.Descriptor)
)

const (
	JobsKeyPrefix   = "/cronjob/job/"    // job prefix
	OnceKeyPrefix   = "/cronjob/once/"   // job that run immediately
	LockKeyPrefix   = "/cronjob/lock/"   // job lock (only for single-node mode job)
	ProcKeyPrefix   = "/cronjob/proc/"   // running process
	ResultKeyPrefix = "/cronjob/result/" // task result (logs and status)
)

type Config struct {
	Enable bool

	etcdConf        *etcd.Config
	ReqTimeout      int   // 请求操作ETCD的超时时间，单位秒
	RequireLockTime int64 // 抢锁等待时间，单位秒

	HostName string
	AppIP    string

	logger   *mylog.Logger
	parser   parser.Parser
	wrappers []cron.JobWrapper
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		ReqTimeout: 3,
		Enable:     false,
	}
}

// StdConfig returns standard configuration information
func StdConfig(key string) *Config {

	// 解析配置

	return nil
}

// Build new a instance
func (c *Config) Build() *Worker {
	if !c.Enable {
		return nil
	}
	c.HostName = util.ReturnHostName()
	c.AppIP = util.ReturnAppIp()

	if c.logger == nil {
		c.logger = mylog.NewLogger()
	}
	c.logger = c.logger.With(zap.String("mod", "worker"))

	// default
	c.parser = myParser
	// 默认前面有任务执行，则直接跳过不执行
	c.wrappers = append(c.wrappers, skipIfStillRunning(c.logger))

	return NewWorker(c)
}
