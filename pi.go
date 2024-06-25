package pi

import (
	"context"

	"go.viam.com/rdk/components/board"
	picommon "go.viam.com/rdk/components/board/pi/common"
	piImpl "go.viam.com/rdk/components/board/pi/impl"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	"go.viam.com/utils"
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

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("pi"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) error {
	pigpio, err := module.NewModuleFromArgs(ctx, logger)

	if err != nil {
		return err
	}
	pigpio.AddModelFromRegistry(ctx, board.API, picommon.Model)

	err = pigpio.Start(ctx)
	defer pigpio.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
