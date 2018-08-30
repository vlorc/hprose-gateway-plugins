package session

import (
	"context"
	"errors"
	"github.com/vlorc/hprose-gateway-core/plugin"
	"github.com/vlorc/hprose-gateway-types"
	"regexp"
	"strings"
)

type sessionParamFactory struct{}

func init() {
	plugin.Register(sessionParamFactory{}, "session")
}

func ignore(mode, match string) (result func(string) bool) {
	switch mode {
	case "prefix":
		result = func(s string) bool {
			return strings.HasPrefix(s, match)
		}
	case "suffix":
		result = func(s string) bool {
			return strings.HasSuffix(s, match)
		}
	case "find":
		result = func(s string) bool {
			return strings.Index(s, match) >= 0
		}
	case "regexp":
		result = regexp.MustCompile(match).MatchString
	default:
		result = func(string) bool {
			return false
		}
	}
	return
}

func (sessionParamFactory) Instance(ctx context.Context, param map[string]string) types.Plugin {
	factory := ctx.Value("SessionFactory").(func(string) SessionFactory)
	return &sessionParam{
		factory: factory(param["secret"]),
		err:     errors.New(types.String(param["error"],"Illegal token")),
		ignore:  ignore(param["ignore.mode"], param["ignore.match"]),
		id:      param["appid"],
		level:   int(types.Integer(param["level"],60000)),
	}
}
