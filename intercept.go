package job

import (
	"go.uber.org/zap"
	"myCron/mylog"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// delayIfStillRunning serializes jobs, delaying subsequent runs until the
// previous one is complete. Jobs running after a delay of more than a minute
// have the delay logged at Info.
func delayIfStillRunning(logger *mylog.Logger) JobWrapper {
	return func(j CronJob) cron.Job {
		var mu sync.Mutex
		return cron.FuncJob(func() {
			start := time.Now()
			mu.Lock()
			defer mu.Unlock()
			if dur := time.Since(start); dur > time.Minute {
				logger.Info("cron delay", zap.String("duration", dur.String()))
			}
			j.Run()
		})
	}
}

// skipIfStillRunning skips an invocation of the Job if a previous invocation is
// still running. It logs skips to the given logger at Info level.
func skipIfStillRunning(logger *mylog.Logger) JobWrapper {
	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
	return func(j CronJob) CronJob {
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				j.Run()
				ch <- v
			default:
				logger.Info("cron skip")
			}
		})
	}
}
