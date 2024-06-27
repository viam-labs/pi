package pi

import (
	"context"
	"syscall"

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

func SigStr(sig syscall.Signal) string {

	switch sig {
	case syscall.SIGHUP:
		return "SIGHUP"
	case syscall.SIGINT:
		return "SIGINT"
	case syscall.SIGQUIT:
		return "SIGQUIT"
	case syscall.SIGABRT:
		return "SIGABRT"
	case syscall.SIGUSR1:
		return "SIGUSR1"
	case syscall.SIGUSR2:
		return "SIGUSR2"
	case syscall.SIGTERM:
		return "SIGTERM"
	default:
		return "<UNKNOWN>"
	}
}
