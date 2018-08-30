package limiter


import (
	"context"
	"errors"
	"github.com/vlorc/hprose-gateway-types"
	"golang.org/x/time/rate"
	"github.com/vlorc/hprose-gateway-core/plugin"
)

type limiterFactory struct{}

func init() {
	plugin.Register(limiterFactory{}, "limiter")
}

func (limiterFactory) Instance(ctx context.Context, param map[string]string) types.Plugin {
	l := &limiter{
		err: errors.New(types.String(param["error"],"Too frequent requests")),
		interval: types.Integer(param["interval"], 1000),
	}
	l.total = types.Integer(param["total"], l.interval)
	switch types.String(param["mode"], "discard") {
	case "wait":
		l.allow = func(ctx context.Context, lim *rate.Limiter) bool {
			return nil == lim.Wait(ctx)
		}
	case "discard":
		l.allow = func(ctx context.Context, lim *rate.Limiter) bool {
			return lim.Allow()
		}
	}
	return l
}
