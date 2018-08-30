package counter

import (
	"context"
	"github.com/vlorc/hprose-gateway-types"
	"reflect"
)

type counter struct {
	named func(string) string
	incr func (string)
}

func (c *counter) Level() int {
	return 65502
}

func (c *counter) Close() error {
	return nil
}

func (c *counter) Name() string {
	return "counter"
}

func (c *counter) Handler(next types.InvokeHandler, ctx context.Context, method string, args []reflect.Value) ([]reflect.Value, error) {
	c.incr(c.named(method))
	return next(ctx, method, args)
}
