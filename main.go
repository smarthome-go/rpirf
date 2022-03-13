/*
Experimental library for interacting with 433mz devices, will only support sending, only tested with protocol 1
Protocol and logic ported from https://github.com/milaq/rpi-rf

Example
func main() {
	device := rpirf.NewRF(17, 1, 10, 180, 24)
	device.Send(5121438)
	device.Cleanup()
}
*/

package rpirf

import (
	"fmt"
	"runtime"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// Initializes the GPIO device
// Provide a pin number, a protocol, how often the signal should be sent, the pulse length and the data length
// The pin number will be the `BCM / bcm2835` pin, not the physical one
func NewRF(pinNumber uint8, protocolIndex uint8, repeat uint8, pulseLength uint16, length uint8) (RFDevice, error) {
	// If the architecture is not arm, the library can not run (not a Raspberry Pi)
	if runtime.GOARCH != "arm" {
		return RFDevice{}, ErrNonArm
	}
	if err := rpio.Open(); err != nil {
		// If the memory range from /dev/mem could not be opened, return the failed to initialize error
		return RFDevice{}, ErrInitialize
	}
	// Initialize the bcm2835 pin
	pin := rpio.Pin(pinNumber)
	// Set pin as output
	pin.Output()
	device := RFDevice{
		Pin:           pin,
		TxEnabled:     true,
		TxProto:       protocolIndex - 1,
		TxRepeat:      repeat,
		TxLength:      length,
		TxPulseLength: pulseLength,
	}
	return device, nil
}

// Disables the transmitter and frees the allocated GPIO pin
func (device *RFDevice) Cleanup() error {
	if !device.TxEnabled {
		return ErrCleanWOInitialized
	}
	device.TxEnabled = false
	if err := rpio.Close(); err != nil {
		return ErrCleanup
	}
	return nil
}

// Sends the provided decimal number as a binary code
func (device *RFDevice) Send(code int) error {
	if !device.TxEnabled {
		return ErrNotInitialized
	}
	binary := fmt.Sprintf("%0*b", device.TxLength+2, code)[2:]
	if code > 16777216 {
		device.TxLength = 32
	}
	device.sendBinary(binary)
	return nil
}

// Sends the preprocessed binary code
func (device *RFDevice) sendBinary(code string) {
	for i := 0; i < int(device.TxRepeat); i++ {
		for b := 0; b < int(device.TxLength); b++ {
			if len(code) > b && code[b] == '0' {
				device.txL0()
			} else {
				device.txL1()
			}
		}
		device.txSync()
		time.Sleep(time.Microsecond)
	}
}

// Send a `0` bit
func (device *RFDevice) txL0() {
	device.txWaveform(protocols[device.TxProto].ZeroHigh, protocols[device.TxProto].ZeroLow)
}

// Send a `1` bit
func (device *RFDevice) txL1() {
	device.txWaveform(protocols[device.TxProto].OneHigh, protocols[device.TxProto].OneLow)
}

// Send a sync
func (device *RFDevice) txSync() {
	device.txWaveform(protocols[device.TxProto].SyncHigh, protocols[device.TxProto].SyncLow)
}

// Sends a generic waveform
func (device *RFDevice) txWaveform(highPulses uint8, lowPulses uint8) {
	device.Pin.High()
	device.sleep((float64(highPulses) * float64(device.TxPulseLength)) / 1000000)
	device.Pin.Low()
	device.sleep((float64(lowPulses) * float64(device.TxPulseLength)) / 1000000)
}

// customized sleep function
func (device *RFDevice) sleep(delay float64) {
	newDelay := delay / 100
	end := float64(time.Now().UnixMicro())/float64(1000000) + delay - newDelay
	for float64(time.Now().UnixMicro())/float64(1000000) < end {
		time.Sleep(time.Microsecond * time.Duration((newDelay)*100000))
	}
}
