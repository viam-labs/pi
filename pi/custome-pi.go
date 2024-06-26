package pi

import (
	"context"

	"go.viam.com/rdk/components/board"
	piImpl "go.viam.com/rdk/components/board/pi/impl"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var Model = resource.NewModel("viam-labs", "board", "rpi4")

func init() {
	resource.RegisterComponent(board.API, Model, resource.Registration[board.Board, *piImpl.Config]{
		Constructor: newPigpio,
	},
	)
}

func newPigpio(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (board.Board, error) {
	return piImpl.NewPigpio(ctx, conf.ResourceName(), conf, logger)
}
