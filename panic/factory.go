package panic

import (
	"context"
	"github.com/vlorc/hprose-gateway-core/option"
	"github.com/vlorc/hprose-gateway-core/plugin"
	"github.com/vlorc/hprose-gateway-types"
)

type panicRecordFactory struct{}

func init() {
	plugin.Register(panicRecordFactory{}, "panic")
}

func (panicRecordFactory) Instance(ctx context.Context, param map[string]string) types.Plugin {
	return &panicRecord{
		log: ctx.Value("option").(*option.GatewayOption).Log,
	}
}
