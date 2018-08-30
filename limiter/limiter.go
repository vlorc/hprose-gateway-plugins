package limiter

import (
	"context"
	"github.com/vlorc/hprose-gateway-types"
	"golang.org/x/time/rate"
	"reflect"
	"sync"
)

type limiter struct {
	tab      sync.Map
	err      error
	interval int64
	total    int64
	allow    func(context.Context, *rate.Limiter) bool
}

func (l *limiter) Level() int {
	return 65500
}

func (l *limiter) Close() error {
	return nil
}

func (l *limiter) Name() string {
	return "limiter"
}

func (l *limiter) Handler(next types.InvokeHandler, ctx context.Context, method string, args []reflect.Value) (val []reflect.Value, err error) {
	lim, _ := l.tab.LoadOrStore(method, rate.NewLimiter(rate.Limit(l.interval), int(l.total)))
	if l.allow(ctx, lim.(*rate.Limiter)) {
		val, err = next(ctx, method, args)
	} else {
		err = l.err
	}
	return
}
