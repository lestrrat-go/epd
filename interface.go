package epd

import (
	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

const (
	Width  = 128 // hardcoded for 2.13 inch ePaper for now
	Height = 250
)

const (
	OffsetSET = 7
	OffsetCLR = 10
)

const (
	DriverOutputControl            byte = 0x01
	BoosterSoftStartControl             = 0x0C
	GateScanStartPosition               = 0x0F
	DeepSleepMode                       = 0x10
	DataEntryModeSetting                = 0x11
	SwReset                             = 0x12
	TemperatuteSendorControl            = 0x1A
	MasterActivation                    = 0x20
	DisplayUpdateControl1               = 0x21
	DisplayUpdateControl2               = 0x22
	WriteRAM                            = 0x24
	WriteVCOMRegister                   = 0x2C
	WriteLUTRegister                    = 0x32
	SetDummyLinePeriod                  = 0x3A
	SetGateTime                         = 0x3B
	BorderWaveformControl               = 0x3C
	SetRamXAddressStartEndPosition      = 0x44
	SetRamYAddressStartEndPosition      = 0x45
	SetRamXAddressCounter               = 0x4E
	SetRamYAddressCounter               = 0x4F
	TerminateFrameReadWrite             = 0xF
)

const (
	defPinBusy = 24
	defPinCS   = 8
	defPinDC   = 25
	defPinRST  = 17
)

type EPD struct {
	busy   gpio.InputPin
	cs     gpio.OutputPin
	dc     gpio.OutputPin
	lut    []byte
	reset  gpio.OutputPin
	spi    *spi.Device
	width  uint8
	height uint8
}

var lutFullUpdate = []byte{
	0x22, 0x55, 0xAA, 0x55, 0xAA, 0x55, 0xAA, 0x11,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E, 0x1E,
	0x01, 0x00, 0x00, 0x00, 0x00, 0x0,
}
