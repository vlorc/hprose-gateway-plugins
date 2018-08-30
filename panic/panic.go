package panic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vlorc/hprose-gateway-types"
	"go.uber.org/zap"
	"reflect"
	"runtime"
)

type panicRecord struct {
	log func() *zap.Logger
}

func (p *panicRecord) Level() int {
	return 65530
}

func (p *panicRecord) Close() error {
	return nil
}

func (p *panicRecord) Name() string {
	return "panic"
}

func (p *panicRecord) Handler(next types.InvokeHandler, ctx context.Context, method string, args []reflect.Value) (val []reflect.Value, err error) {
	defer func() {
		if it := recover(); nil != it {
			err = errors.New(fmt.Sprint(it))

			p.log().Debug("Panic",
				zap.String("method", method),
				zap.Error(err),
				zap.String("params", params(args)),
				zap.String("stack", string(stack())),
			)
		}
	}()
	val, err = next(ctx, method, args)
	return
}

func params(args []reflect.Value) string {
	arr := make([]interface{}, len(args))
	for i, v := range args {
		if v.IsValid() {
			arr[i] = v.Interface()
		} else {
			arr[i] = nil
		}
	}
	buf, _ := json.Marshal(arr)
	return string(buf)
}

func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}
