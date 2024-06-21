// Package pi implements raspberry pi model 4
package pi

// #include <stdlib.h>
// #include <pigpio.h>
// #include "pi.h"
// #cgo LDFLAGS: -lpigpio
import "C"
import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/board/mcp3008helper"
	picommon "go.viam.com/rdk/components/board/pi/common"
	"go.viam.com/rdk/components/board/pinwrappers"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var (
	Model             = resource.NewModel("viam-labs", "board", "rpi")
	pigpioInitialized bool
	instanceMu        sync.RWMutex
	instances         = map[*piPigpio]struct{}{}
)

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

func initializePigpio() error {
	instanceMu.Lock()
	defer instanceMu.Unlock()

	if pigpioInitialized {
		return nil
	}

	resCode := C.gpioInitialise()
	if resCode < 0 {
		// failed to init, check for common causes
		_, err := os.Stat("/sys/bus/platform/drivers/raspberrypi-firmware")
		if err != nil {
			return errors.New("not running on a pi")
		}
		if os.Getuid() != 0 {
			return errors.New("not running as root, try sudo")
		}
		return picommon.ConvertErrorCodeToMessage(int(resCode), "error")
	}

	pigpioInitialized = true
	return nil
}

func newBoard(ctx context.Context, _ resource.Dependencies, conf resource.Config, logger logging.Logger,
) (board.Board, error) {
	internals := C.gpioCfgGetInternals()
	internals |= C.PI_CFG_NOSIGHANDLER
	resCode := C.gpioCfgSetInternals(internals)
	if resCode < 0 {
		return nil, picommon.ConvertErrorCodeToMessage(int(resCode), "gpioCfgSetInternals failed with code")
	}

	if err := initializePigpio(); err != nil {
		return nil, err
	}

	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	piInstance := &piPigpio{
		Named:      conf.ResourceName().AsNamed(),
		logger:     logger,
		isClosed:   false,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
	}

	if err := piInstance.Reconfigure(ctx, nil, cfg); err != nil {
		// This has to happen outside of the lock to avoid a deadlock with interrupts.
		C.gpioTerminate()
		instanceMu.Lock()
		pigpioInitialized = false
		instanceMu.Unlock()
		logger.CError(ctx, "Pi GPIO terminated due to failed init.")
		return nil, err
	}
	return piInstance, nil
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
