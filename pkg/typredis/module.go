package typredis

import (
	"github.com/typical-go/typical-go/pkg/common"
	"github.com/typical-go/typical-go/pkg/typcore"
	"github.com/urfave/cli/v2"
)

// Module of Redis
type Module struct{}

// New Redis Module
func New() *Module {
	return &Module{}
}

// BuildCommands of module
func (r *Module) BuildCommands(c *typcore.BuildContext) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "redis",
			Usage: "Redis Tool",
			Before: func(ctx *cli.Context) error {
				return common.LoadEnvFile()
			},
			Subcommands: []*cli.Command{
				r.consoleCmd(c),
			},
		},
	}
}

// Configure Redis
func (r *Module) Configure(loader typcore.ConfigLoader) (prefix string, spec, loadFn interface{}) {
	prefix = "REDIS"
	spec = &Config{}
	loadFn = func() (cfg Config, err error) {
		err = loader.Load(prefix, &cfg)
		return
	}
	return
}

// Provide dependencies
func (r *Module) Provide() []interface{} {
	return []interface{}{
		r.connect,
	}
}

// Prepare the module
func (r *Module) Prepare() []interface{} {
	return []interface{}{
		r.ping,
	}
}

// Destroy dependencies
func (r *Module) Destroy() []interface{} {
	return []interface{}{
		r.disconnect,
	}
}
