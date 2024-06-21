package main

import "go.viam.com/rdk/resource"

// DigitalInterruptConfig describes the configuration of digital interrupt for a board.
type DigitalInterruptConfig struct {
	Name string `json:"name"`
	Pin  string `json:"pin"`
}

// Validate ensures all parts of the config are valid.
func (config *DigitalInterruptConfig) Validate(path string) error {
	if config.Name == "" {
		return resource.NewConfigValidationFieldRequiredError(path, "name")
	}
	if config.Pin == "" {
		return resource.NewConfigValidationFieldRequiredError(path, "pin")
	}
	return nil
}
