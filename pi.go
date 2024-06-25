package pi

import (
	"context"

	"go.viam.com/rdk/components/board"
	picommon "go.viam.com/rdk/components/board/pi/common"
	piImpl "go.viam.com/rdk/components/board/pi/impl"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

func init() {
	resource.RegisterComponent(board.API, picommon.Model, resource.Registration[board.Board, *piImpl.Config]{
		Constructor: newPigpio,
	},
	)
}

func newPigpio(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (board.Board, error) {
	return piImpl.NewPigpio(ctx, conf.ResourceName(), conf, logger)
}
