package etcd

import (
	"context"
	"go.uber.org/zap"
	"myCron/mylog"
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// Watch A watch only tells the latest revision
type Watch struct {
	revision  int64
	cancel    context.CancelFunc
	eventChan chan *clientv3.Event
	lock      *sync.RWMutex
	logger    *mylog.Logger

	incipientKVs []*mvccpb.KeyValue
}

// C ...
func (w *Watch) C() chan *clientv3.Event {
	return w.eventChan
}

// IncipientKeyValues incipient key and values
func (w *Watch) IncipientKeyValues() []*mvccpb.KeyValue {
	return w.incipientKVs
}

// NewWatch ...
func WatchPrefix(client *clientv3.Client, ctx context.Context, prefix string) (*Watch, error) {
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var w = &Watch{
		revision:     resp.Header.Revision,
		eventChan:    make(chan *clientv3.Event, 100),
		incipientKVs: resp.Kvs,
	}

	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		w.cancel = cancel
		rch := client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
		for {
			for n := range rch {
				if n.CompactRevision > w.revision {
					w.revision = n.CompactRevision
				}
				if n.Header.GetRevision() > w.revision {
					w.revision = n.Header.GetRevision()
				}
				if err := n.Err(); err != nil {
					mylog.DefaultLogger.Error("watch request err",
						zap.String("error string", prefix+" register err"), zap.Error(err))
					continue
				}
				for _, ev := range n.Events {
					select {
					case w.eventChan <- ev:
					default:
						mylog.DefaultLogger.Error("watch etcd with prefix", zap.Any("err", "block event chan, drop event message"))
					}
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			w.cancel = cancel
			if w.revision > 0 {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
			} else {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
			}
		}
	}()

	return w, nil
}

// Close close watch
func (w *Watch) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}
