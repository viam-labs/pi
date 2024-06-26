package main

import (
	"context"
	customPi "pi-module"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("rpi4"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) error {
	pigpio, err := module.NewModuleFromArgs(ctx, logger)

	if err != nil {
		return err
	}
	pigpio.AddModelFromRegistry(ctx, board.API, customPi.Model)

	err = pigpio.Start(ctx)
	defer pigpio.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
