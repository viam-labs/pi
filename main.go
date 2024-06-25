package main

import (
	"context"

	goutils "go.viam.com/utils"

	"go.viam.com/rdk/components/board"
	picommon "go.viam.com/rdk/components/board/pi/common"
	piImpl "go.viam.com/rdk/components/board/pi/impl"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
)

func init() {
	resource.RegisterComponent(board.API, picommon.Model,
		resource.Registration[board.Board, *piImpl.Config]{
			Constructor: piImpl.NewPigpio,
		},
	)
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	modalModule, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	modalModule.AddModelFromRegistry(ctx, board.API, picommon.Model)

	err = modalModule.Start(ctx)
	defer modalModule.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func main() {
	goutils.ContextualMain(mainWithArgs, logging.NewLogger("RPI"))
}
