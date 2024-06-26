package customPi

import (
	"context"

	"go.viam.com/rdk/components/board"
	piImpl "go.viam.com/rdk/components/board/pi/impl"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils"
)

var Model = resource.NewModel("viam-labs", "board", "rpi4")

func init() {
	resource.RegisterComponent(board.API, Model, resource.Registration[board.Board, *Config]{
		Constructor: newPigpio,
	},
	)
}

type Config struct {
	field1 int    `json: "one"`
	field2 string `json: "two"`
}

func (cfg *Config) Validate(path string) ([]string, error) {
	if cfg.field1 == 0 {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "one")
	}

	if cfg.field2 == "" {
		return nil, utils.NewConfigValidationFieldRequiredError(path, "two")
	}

	// TODO(7): return implicit dependencies if needed as the first value
	return []string{}, nil
}
func newPigpio(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (board.Board, error) {
	return piImpl.NewPigpio(ctx, conf.ResourceName(), conf, logger)
}
