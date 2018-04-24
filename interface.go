package epd

import (
	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

const (
	Width  uint8 = 128 // hardcoded for 2.13 inch ePaper for now
	Height uint8 = 250
)

const (
	OffsetSET = 7
	OffsetCLR = 10
)

type Command byte

const (
	DriverOutputControl            Command = 0x01
	BoosterSoftStartControl        Command = 0x0C
	GateScanStartPosition          Command = 0x0F
	DeepSleepMode                  Command = 0x10
	DataEntryModeSetting           Command = 0x11
	SwReset                        Command = 0x12
	TemperatuteSendorControl       Command = 0x1A
	MasterActivation               Command = 0x20
	DisplayUpdateControl1          Command = 0x21
	DisplayUpdateControl2          Command = 0x22
	WriteRAM                       Command = 0x24
	WriteVCOMRegister              Command = 0x2C
	WriteLUTRegister               Command = 0x32
	SetDummyLinePeriod             Command = 0x3A
	SetGateTime                    Command = 0x3B
	BorderWaveformControl          Command = 0x3C
	SetRamXAddressStartEndPosition Command = 0x44
	SetRamYAddressStartEndPosition Command = 0x45
	SetRamXAddressCounter          Command = 0x4E
	SetRamYAddressCounter          Command = 0x4F
	TerminateFrameReadWrite        Command = 0xFF
)

const (
	defPinBusy = 24
	defPinCS   = 8
	defPinDC   = 25
	defPinRST  = 17
)

type EPD struct {
	buffer []byte
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
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
}
