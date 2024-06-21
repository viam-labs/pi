package main

import (
	"context"

	goutils "go.viam.com/utils"
	"main.go/pi"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
)

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) (err error) {
	modalModule, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	modalModule.AddModelFromRegistry(ctx, board.API, pi.Model)

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
