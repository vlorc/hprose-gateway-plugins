package hash

import (
	"context"
	"github.com/vlorc/hprose-gateway-types"
	"sync"
	"hash"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/vlorc/hprose-gateway-core/plugin"
)

type hashFactory struct{}


func init() {
	plugin.Register(hashFactory{}, "hash")
}

var table = map[string]func()hash.Hash {
	"md5": md5.New,
	"sha1": sha1.New,
	"sha256": sha256.New,
	"sha512": sha512.New,
}
func (hashFactory) Instance(ctx context.Context, param map[string]string) types.Plugin {
	name := types.String(param["name"], "hash")
	factory := table[types.String(param["hash"], "md5")]
	return &hasher{
		name: name,
		pool: sync.Pool{
			New: func() interface{} {
				return factory()
			},
		},
	}
}
