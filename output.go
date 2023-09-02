package rpirf

import (
	"runtime"

	"github.com/stianeikeland/go-rpio/v4"
	"github.com/warthog618/gpiod"
)

type HardwareOutputKind uint8

const (
	HardwareOutputRpi HardwareOutputKind = iota
	HardwareOutputCdev
)

type HardwareOutput interface {
	Kind() HardwareOutputKind
	Low() error
	High() error
}

type HardwareOutputRaspberryPi struct {
	Pin rpio.Pin
}

func (self HardwareOutputRaspberryPi) Kind() HardwareOutputKind {
	return HardwareOutputRpi
}

func (self HardwareOutputRaspberryPi) Low() error {
	self.Pin.Low()
	return nil
}

func (self HardwareOutputRaspberryPi) High() error {
	self.Pin.High()
	return nil
}

func NewRaspberryPi(pinNumber uint) (HardwareOutputRaspberryPi, error) {
	// If the architecture is not arm, the library can not run (not a Raspberry Pi)
	if runtime.GOARCH != "arm" {
		return HardwareOutputRaspberryPi{}, ErrNonArm
	}

	if err := rpio.Open(); err != nil {
		// If the memory range from /dev/mem could not be opened, return the failed to initialize error
		return HardwareOutputRaspberryPi{}, ErrInitialize
	}

	// Initialize the bcm pin
	pin := rpio.Pin(pinNumber)

	// Set pin as output
	pin.Output()

	return HardwareOutputRaspberryPi{
		Pin: pin,
	}, nil
}

type HardwareOutputCharacterdev struct {
	LineHandle *gpiod.Line
}

func (self HardwareOutputCharacterdev) Kind() HardwareOutputKind {
	return HardwareOutputCdev
}

func (self HardwareOutputCharacterdev) Low() error {
	return self.LineHandle.SetValue(0)
}

func (self HardwareOutputCharacterdev) High() error {
	return self.LineHandle.SetValue(1)
}

func NewCharacterDev(devicePath string, pinNumber int) (HardwareOutputCharacterdev, error) {
	out, err := gpiod.RequestLine(devicePath, pinNumber, gpiod.AsOutput(0))
	if err != nil {
		return HardwareOutputCharacterdev{}, err
	}

	return HardwareOutputCharacterdev{
		LineHandle: out,
	}, nil
}
