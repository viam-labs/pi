// Package pi implements raspberry pi model 4
package pi

import (
	"context"
	"fmt"
	"sync"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/board/mcp3008helper"
	"go.viam.com/rdk/components/board/pinwrappers"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var Model = resource.NewModel("viam-labs", "board", "rpi")

func init() {
	resource.RegisterComponent(board.API, Model, resource.Registration[board.Board, *Config]{
		Constructor: newRPI,
	})
}

type Config struct {
	AnalogReaders     []mcp3008helper.MCP3008AnalogConfig `json:"analogs,omitempty"`
	DigitalInterrupts []board.DigitalInterruptConfig      `json:"digital_interrupts,omitempty"`
}

func (conf *Config) Validate(path string) ([]string, error) {
	for idx, c := range conf.AnalogReaders {
		if err := c.Validate(fmt.Sprintf("%s.%s.%d", path, "analogs", idx)); err != nil {
			return nil, err
		}
	}
	for idx, c := range conf.DigitalInterrupts {
		if err := c.Validate(fmt.Sprintf("%s.%s.%d", path, "digital_interrupts", idx)); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func newRPI(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (board.Board, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}

	// Create a cancelable context for custom sensor
	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	r := &piPigpio{
		name:       rawConf.ResourceName(),
		logger:     logger,
		cfg:        conf,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}

	// TODO(10): If your custom component has dependencies, perform any checks you need to on them.

	// The Reconfigure() method changes the values on the customResource based on the attributes in the component config
	if err := r.Reconfigure(ctx, deps, rawConf); err != nil {
		logger.Error("Error configuring module with ", rawConf)
		return nil, err
	}

	return r, nil
}

type piPigpio struct {
	resource.Named

	mu            sync.Mutex
	cancelCtx     context.Context
	cancelFunc    context.CancelFunc
	duty          int // added for mutex
	gpioConfigSet map[int]bool
	analogReaders map[string]*pinwrappers.AnalogSmoother
	// `interrupts` maps interrupt names to the interrupts. `interruptsHW` maps broadcom addresses
	// to these same values. The two should always have the same set of values.
	interrupts   map[string]ReconfigurableDigitalInterrupt
	interruptsHW map[uint]ReconfigurableDigitalInterrupt
	logger       logging.Logger
	isClosed     bool

	activeBackgroundWorkers sync.WaitGroup
}
