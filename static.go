package rpirf

import (
	"errors"

	"github.com/stianeikeland/go-rpio/v4"
)

type Protocol struct {
	Pulselength uint16
	SyncHigh    uint8
	SyncLow     uint8
	ZeroHigh    uint8
	ZeroLow     uint8
	OneHigh     uint8
	OneLow      uint8
}

type RFDevice struct {
	Pin           rpio.Pin
	TxEnabled     bool
	TxProto       uint8
	TxRepeat      uint8
	TxLength      uint8
	TxPulseLength uint16
}

var protocols = []Protocol{
	{Pulselength: 350, SyncHigh: 1, SyncLow: 31, ZeroHigh: 1, ZeroLow: 3, OneHigh: 3, OneLow: 1},
	{Pulselength: 650, SyncHigh: 1, SyncLow: 10, ZeroHigh: 1, ZeroLow: 2, OneHigh: 2, OneLow: 1},
	{Pulselength: 100, SyncHigh: 30, SyncLow: 71, ZeroHigh: 4, ZeroLow: 11, OneHigh: 9, OneLow: 6},
	{Pulselength: 380, SyncHigh: 1, SyncLow: 6, ZeroHigh: 1, ZeroLow: 3, OneHigh: 3, OneLow: 1},
	{Pulselength: 500, SyncHigh: 6, SyncLow: 14, ZeroHigh: 1, ZeroLow: 2, OneHigh: 2, OneLow: 1},
	{Pulselength: 200, SyncHigh: 1, SyncLow: 10, ZeroHigh: 1, ZeroLow: 5, OneHigh: 1, OneLow: 1},
}

var (
	ErrNotInitialized     = errors.New("cannot send code: device is not initialized. make sure to initialize the device first")
	ErrCleanWOInitialized = errors.New("cannot cleanup a non-initialized device")
	ErrCleanup            = errors.New("failed to cleanup: could not close rpio")
	ErrInitialize         = errors.New("failed to initialize device: could not open rpio")
	ErrNonArm             = errors.New("unsupported architecture: this library only works on the raspberry pi (arm)")
)
