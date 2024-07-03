package main

import (
	"context"
	piimpl "pi-module/pi"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"
)

//var Model = resource.NewModel("viam-labs", "board", "rpi4")

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("pi"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) error {
	rtkSystem, err := module.NewModuleFromArgs(ctx, logger)

	if err != nil {
		return err
	}
	rtkSystem.AddModelFromRegistry(ctx, board.API, piimpl.Model)

	err = rtkSystem.Start(ctx)
	defer rtkSystem.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
