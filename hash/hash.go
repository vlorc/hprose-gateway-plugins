package hash

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/vlorc/hprose-gateway-types"
	"gopkg.in/vmihailenco/msgpack.v3"
	"hash"
	"reflect"
	"sync"
)

type hasher struct {
	name string
	pool sync.Pool
}

func (p *hasher) Level() int {
	return 65501
}

func (p *hasher) Close() error {
	return nil
}

func (p *hasher) Name() string {
	return "hash"
}

func (p *hasher) Handler(next types.InvokeHandler, ctx context.Context, method string, args []reflect.Value) ([]reflect.Value, error) {
	md5.New()
	h := p.pool.Get().(hash.Hash)
	w := msgpack.NewEncoder(h).SortMapKeys(true)
	w.Encode(method)
	for _, v := range args {
		w.Encode(v.Interface())
	}
	ctx = context.WithValue(ctx, p.name, hex.EncodeToString(h.Sum(nil)))
	h.Reset()
	p.pool.Put(h)
	return next(ctx, method, args)
}
