package session

import (
	"context"
	"github.com/vlorc/hprose-gateway-types"
	"reflect"
)

type sessionParam struct {
	factory SessionFactory
	err     error
	ignore  func(string) bool
	id      string
	level   int
}

func (s *sessionParam) Level() int {
	return 65530
}

func (s *sessionParam) Close() error {
	return nil
}

func (s *sessionParam) Name() string {
	return "session"
}

func (s *sessionParam) Handler(next types.InvokeHandler, ctx context.Context, method string, args []reflect.Value) ([]reflect.Value, error) {
	if len(args) <= 0 {
		return nil, s.err
	}
	if s.ignore(method) {
		return next(context.WithValue(ctx, "appid", s.id), method, args)
	}
	str, ok := args[len(args)-1].Interface().(string)
	if !ok {
		return nil, s.err
	}
	session, err := s.factory.Instance(str)
	if nil != err {
		return nil, s.err
	}
	ctx = context.WithValue(ctx, "appid", session.Appid())
	args[len(args)-1] = reflect.ValueOf(session)
	return next(ctx, method, args)
}
