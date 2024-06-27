package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pi-module/pi"
	"syscall"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("custom-pi"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	pigpio, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}
	pigpio.AddModelFromRegistry(ctx, board.API, pi.Model)

	err = pigpio.Start(ctx)
	if err != nil {
		return err
	}
	fmt.Print("HEEEEEEEELP!!!")
	go func() {
		select {
		case sig := <-signalChan:
			logger.Infof("Received signal: %s", sig)

			file, err := os.Create("received_signal.txt")
			if err != nil {
				logger.Errorf("Failed to create file: %v", err)
			} else {
				defer file.Close()
				_, err = file.WriteString(sig.String())
				if err != nil {
					logger.Errorf("Failed to write to file: %v", err)
				}
			}
		}
	}()
	defer pigpio.Close(ctx)

	<-ctx.Done()
	return nil
}
