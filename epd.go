package epd

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ecc1/gpio"
	"github.com/ecc1/spi"
)

const hoge = 1

func New() *EPD {
	var pinBUSY gpio.InputPin
	var pinCS gpio.OutputPin
	var pinDC gpio.OutputPin
	var pinRST gpio.OutputPin

	d, err := spi.Open("/dev/spidev0.0", 2000000, 0)
	if err != nil {
		panic("failed to open SPI device: " + err.Error())
	}
	d.SetMode(0)

	{
		p, err := gpio.Input(defPinBusy, false)
		if err != nil {
			panic("failed to setup BUSY pin: " + err.Error())
		}
		pinBUSY = p
	}

	{
		p, err := gpio.Output(defPinCS, false, false)
		if err != nil {
			panic("failed to setup CS pin: " + err.Error())
		}
		pinCS = p
	}

	{
		p, err := gpio.Output(defPinDC, false, false)
		if err != nil {
			panic("failed to setup DC pin: " + err.Error())
		}
		pinDC = p
	}

	{
		p, err := gpio.Output(defPinRST, false, false)
		if err != nil {
			panic("failed to setup RST pin: " + err.Error())
		}
		pinRST = p
	}

	e := &EPD{
		busy:   pinBUSY,
		cs:     pinCS,
		dc:     pinDC,
		lut:    lutFullUpdate,
		reset:  pinRST,
		spi:    d,
		width:  Width,
		height: Height,
	}

	e.Reinitialize()
	return e
}

func (e *EPD) Reinitialize() {
	e.SendCommand(DriverOutputControl)
	e.SendData((Height - 1) & 0xff)
	e.SendData(((Height - 1) >> 8) & 0xff)
	e.SendData(0x00)
	e.SendCommand(BoosterSoftStartControl)
	e.SendData(0xd7)
	e.SendData(0xd6)
	e.SendData(0x9D)
	e.SendCommand(WriteVCOMRegister)
	e.SendData(0xA8)
	e.SendCommand(SetDummyLinePeriod)
	e.SendData(0x1A)
	e.SendCommand(SetGateTime)
	e.SendData(0x08)
	e.SendCommand(DataEntryModeSetting)
	e.SendData(0x03)
	e.SetLUT(e.lut)
}

func (e *EPD) Reset() {
	e.reset.Write(false)
	time.Sleep(200 * time.Millisecond)
	e.reset.Write(true)
	time.Sleep(200 * time.Millisecond)
}

func (e *EPD) SendCommand(cmd byte) {
	e.dc.Write(false)
	e.spi.Transfer([]byte{cmd})
}

// SendData sends data through SPI. Arbitrary number of bytes can be
// passed to this method, and they will be each sent in succession
func (e *EPD) SendData(data ...byte) {
	for _, b := range data {
		e.dc.Write(true)
		e.spi.Transfer([]byte{b})
	}
}

func (e *EPD) SetLUT(lut []byte) {
	e.lut = lut
	e.SendCommand(WriteLUTRegister)
	for _, v := range lut {
		e.SendData(v)
	}
}

func (e *EPD) SetMemoryArea(startX, startY, endX, endY uint8) {
	log.Printf("SetMemoryArea")
	// x point must be multiple of 8 or the last 3 bits will be ignored
	e.SendCommand(SetRamXAddressStartEndPosition)
	e.SendData((startX >> 3) & 0xFF)
	e.SendData((endX >> 3) & 0xFF)
	e.SendCommand(SetRamYAddressStartEndPosition)
	e.SendData(startY & 0xFF)
	e.SendData((startY >> 8) & 0xFF)
	e.SendData(endY & 0xFF)
	e.SendData((endY >> 8) & 0xFF)
}

func (e *EPD) SetMemoryPointer(x, y uint8) {
	log.Printf("SetMemoryPointer")
	e.SendCommand(SetRamXAddressCounter)
	e.SendData((x >> 3) & 0xFF)
	e.SendCommand(SetRamYAddressCounter)
	e.SendData(y & 0xFF)
	e.SendData((y >> 8) & 0xFF)
	e.WaitUntilIdle(nil)
}

func (e *EPD) WaitUntilIdle(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	for {
		b, err := e.busy.Read()
		if err == nil && !b {
			return nil
		}
		log.Printf("b = %t", b)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}
	}
	return errors.New(`not reached`)
}

func (e *EPD) ClearFrameMemory(color byte) {
	log.Printf("ClearFrameMemory")
	defer log.Printf("done")
	e.SetMemoryArea(0, 0, e.width-1, e.height-1)
	e.SetMemoryPointer(0, 0)
	log.Printf("Start writing to RAM")
	e.SendCommand(WriteRAM)
	for i := uint8(0); i < (e.width/8)*e.height; i++ {
		e.SendData(color)
	}
}

func (e *EPD) DisplayFrame() {
	e.SendCommand(DisplayUpdateControl2)
	e.SendData(0xC4)
	e.SendCommand(MasterActivation)
	e.SendCommand(TerminateFrameReadWrite)
	e.WaitUntilIdle(nil)
}
