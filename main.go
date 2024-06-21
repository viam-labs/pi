// Package pi implements raspberry pi model 4
package pi

import (
	"fmt"

	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/board/mcp3008helper"
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
