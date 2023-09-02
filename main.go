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
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// Initializes the GPIO device
// Provide a pin number, a protocol, how often the signal should be sent, the pulse length and the data length
// The pin number will be the `BCM / bcm2835` pin, not the physical one
func NewRF(hardware HardwareOutput, protocolIndex uint8, repeat uint8, pulseLength uint16, length uint8) RFDevice {
	device := RFDevice{
		Output:        hardware,
		TxEnabled:     true,
		TxProto:       protocolIndex - 1,
		TxRepeat:      repeat,
		TxLength:      length,
		TxPulseLength: pulseLength,
	}
	return device
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

	return device.sendBinary(binary)
}

// Sends the preprocessed binary code
func (device *RFDevice) sendBinary(code string) error {
	for i := 0; i < int(device.TxRepeat); i++ {
		for b := 0; b < int(device.TxLength); b++ {
			if len(code) > b && code[b] == '0' {
				if err := device.txL0(); err != nil {
					return err
				}
			} else {
				if err := device.txL1(); err != nil {
					return err
				}
			}
		}
		if err := device.txSync(); err != nil {
			return err
		}
		time.Sleep(time.Microsecond)
	}

	return nil
}

// Send a `0` bit
func (device *RFDevice) txL0() error {
	return device.txWaveform(protocols[device.TxProto].ZeroHigh, protocols[device.TxProto].ZeroLow)
}

// Send a `1` bit
func (device *RFDevice) txL1() error {
	return device.txWaveform(protocols[device.TxProto].OneHigh, protocols[device.TxProto].OneLow)
}

// Send a sync
func (device *RFDevice) txSync() error {
	return device.txWaveform(protocols[device.TxProto].SyncHigh, protocols[device.TxProto].SyncLow)
}

// Sends a generic waveform
func (device *RFDevice) txWaveform(highPulses uint8, lowPulses uint8) error {
	if err := device.Output.High(); err != nil {
		return err
	}
	device.sleep((float64(highPulses) * float64(device.TxPulseLength)) / 1_000_000)
	if err := device.Output.Low(); err != nil {
		return err
	}
	device.sleep((float64(lowPulses) * float64(device.TxPulseLength)) / 1000000)

	return nil
}

// customized sleep function
func (device *RFDevice) sleep(delay float64) {
	newDelay := delay / 100
	end := float64(time.Now().UnixMicro())/float64(1_000_000) + delay - newDelay
	for float64(time.Now().UnixMicro())/float64(1_000_000) < end {
		time.Sleep(time.Microsecond * time.Duration((newDelay)*100_000))
	}
}
