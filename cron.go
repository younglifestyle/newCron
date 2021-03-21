package job

import (
	"fmt"
	"myCron/mylog"
	"myCron/util"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type (
	// JobWrapper ...
	JobWrapper = cron.JobWrapper
	// EntryID ...
	EntryID = cron.EntryID
	// Schedule ...
	Schedule = cron.Schedule
	// Job ...
	CronJob = cron.Job
	//NamedJob ..
	NamedJob interface {
		Run() error
	}
)

// FuncJob ...
type FuncJob func() error

// Run ...
func (f FuncJob) Run() error { return f() }

// Name ...
func (f FuncJob) Name() string { return util.FunctionName(f) }

type wrappedLogger struct {
	*mylog.Logger
}

// Info logs routine messages about cron's operation.
func (wl *wrappedLogger) Info(msg string, keysAndValues ...interface{}) {
	wl.Infow("cron "+msg, keysAndValues...)
}

// Error logs an error condition.
func (wl *wrappedLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	wl.Errorw("cron "+msg, append(keysAndValues, "err", err)...)
}

// Cron ...
type Cron struct {
	*Worker
	*cron.Cron
	entries map[string]EntryID
}

func newCron(config *Worker) *Cron {
	c := &Cron{
		Worker: config,
		Cron: cron.New(
			cron.WithLogger(&wrappedLogger{config.logger}),
			cron.WithChain(config.wrappers...),
			cron.WithParser(config.parser),
		),
	}
	return c
}

// Schedule ...
func (c *Cron) Schedule(schedule Schedule, job NamedJob) EntryID {
	if c.ImmediatelyRun {
		schedule = &immediatelyScheduler{
			Schedule: schedule,
		}
	}
	innnerJob := &wrappedJob{
		NamedJob: job,
		logger:   c.Worker.logger,
	}

	return c.Cron.Schedule(schedule, innnerJob)
}

// AddJob ...
func (c *Cron) AddJob(spec string, cmd NamedJob) (EntryID, error) {
	schedule, err := c.Worker.parser.Parse(spec)
	if err != nil {
		return 0, err
	}
	return c.Schedule(schedule, cmd), nil
}

// AddFunc ...
func (c *Cron) AddFunc(spec string, cmd func() error) (EntryID, error) {
	return c.AddJob(spec, FuncJob(cmd))
}

// Remove an entry from being run in the future.
func (c *Cron) Remove(id EntryID) {
	c.Cron.Remove(id)
}

// Run ...
func (c *Cron) Run() {
	c.Worker.logger.Info("run Worker", zap.Int("number of scheduled jobs", len(c.Cron.Entries())))
	c.Cron.Start()
}

// Stop ...
func (c *Cron) Stop() error {
	_ = c.Cron.Stop()
	return nil
}

type immediatelyScheduler struct {
	Schedule
	initOnce uint32
}

// Next ...
func (is *immediatelyScheduler) Next(curr time.Time) (next time.Time) {
	if atomic.CompareAndSwapUint32(&is.initOnce, 0, 1) {
		return curr
	}

	return is.Schedule.Next(curr)
}

type wrappedJob struct {
	NamedJob
	logger *mylog.Logger
}

// Run ...
func (wj wrappedJob) Run() {
	_ = wj.run()
}

func (wj wrappedJob) run() (err error) {
	var fields = []zap.Field{}
	var beg = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			switch rec := rec.(type) {
			case error:
				err = rec
			default:
				err = fmt.Errorf("%v", rec)
			}

			stack := make([]byte, 4096)
			length := runtime.Stack(stack, true)
			fields = append(fields, zap.ByteString("stack", stack[:length]))
		}
		if err != nil {
			fields = append(fields, zap.String("err", err.Error()), zap.Duration("cost", time.Since(beg)))
			wj.logger.Error("Worker", fields...)
		}
	}()

	return wj.NamedJob.Run()
}
