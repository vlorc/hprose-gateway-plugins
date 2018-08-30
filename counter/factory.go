package counter

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/vlorc/hprose-gateway-core/plugin"
	"github.com/vlorc/hprose-gateway-types"
)

type counterFactory struct{}

func init() {
	plugin.Register(counterFactory{}, "counter")
}

func (counterFactory) Instance(ctx context.Context, param map[string]string) types.Plugin {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",types.String(param["redis.addr"], "localhost"),types.Integer(param["redis.port"], 6379)),
		Password: types.String(param["redis.pass"], ""),
		DB:       int(types.Integer(param["redis.db"], 0)),
	})
	if _, err := client.Ping().Result(); nil != err {
		panic(err)
	}
	c := &counter{
		named: func(k string) string {
			return k
		},
		incr: func(k string) {
			client.Incr(k)
		},
	}
	if prefix := types.String(param["prefix"], ""); "" != prefix {
		c.named = func(k string) string {
			return prefix + k
		}
	}
	if key := types.String(param["key"], ""); "" != key {
		c.incr = func(k string) {
			client.HIncrBy(key,k,1)
		}
	}
	return c
}
